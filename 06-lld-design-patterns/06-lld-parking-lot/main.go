// LLD parking lot: entities, Park/Unpark, pluggable pricing, thread-safe.
// Run: go run ./06-lld-design-patterns/06-lld-parking-lot
package main

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
)

type Size int

const (
	Small Size = iota
	Medium
	Large
)

type Vehicle struct {
	Plate string
	Size  Size
}

type Slot struct {
	ID   int
	Size Size
	used bool
}

type Ticket struct {
	ID      int
	Vehicle Vehicle
	Slot    *Slot
	Entry   time.Time
}

// Pricing is a strategy: change the fee rule without touching the lot.
type Pricing interface {
	Fee(size Size, d time.Duration) float64
}

type HourlyPricing struct{ Rate map[Size]float64 }

func (p HourlyPricing) Fee(s Size, d time.Duration) float64 {
	hours := max(math.Ceil(d.Hours()), 1) // round up, minimum one hour
	return hours * p.Rate[s]
}

var ErrFull = errors.New("no free slot")

type Lot struct {
	mu      sync.Mutex
	slots   []*Slot
	tickets map[int]*Ticket
	nextID  int
	pricing Pricing
	now     func() time.Time // injected so tests and demos control time
}

func NewLot(slots []*Slot, p Pricing) *Lot {
	return &Lot{slots: slots, tickets: map[int]*Ticket{}, pricing: p, now: time.Now}
}

// Park picks the first free slot that is big enough (smallest fit first,
// because slots are listed small to large).
func (l *Lot) Park(v Vehicle) (*Ticket, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, s := range l.slots {
		if !s.used && s.Size >= v.Size {
			s.used = true
			l.nextID++
			t := &Ticket{ID: l.nextID, Vehicle: v, Slot: s, Entry: l.now()}
			l.tickets[t.ID] = t
			return t, nil
		}
	}
	return nil, ErrFull
}

func (l *Lot) Unpark(ticketID int) (float64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	t, ok := l.tickets[ticketID]
	if !ok {
		return 0, fmt.Errorf("ticket %d not found", ticketID)
	}
	t.Slot.used = false
	delete(l.tickets, ticketID)
	return l.pricing.Fee(t.Vehicle.Size, l.now().Sub(t.Entry)), nil
}

func main() {
	fmt.Println("=== parking lot ===")
	lot := NewLot(
		[]*Slot{{ID: 1, Size: Small}, {ID: 2, Size: Medium}, {ID: 3, Size: Large}},
		HourlyPricing{Rate: map[Size]float64{Small: 10, Medium: 20, Large: 40}},
	)
	clock := time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)
	lot.now = func() time.Time { return clock }

	t1, _ := lot.Park(Vehicle{"BIKE-1", Small})
	t2, _ := lot.Park(Vehicle{"CAR-1", Medium})
	fmt.Printf("  %s -> slot %d, %s -> slot %d\n", t1.Vehicle.Plate, t1.Slot.ID, t2.Vehicle.Plate, t2.Slot.ID)

	t3, _ := lot.Park(Vehicle{"BIKE-2", Small}) // small slot taken, falls to the large one
	fmt.Printf("  BIKE-2 -> slot %d\n", t3.Slot.ID)
	_, err := lot.Park(Vehicle{"CAR-2", Medium})
	fmt.Println("  CAR-2:", err)

	clock = clock.Add(150 * time.Minute)
	fee, _ := lot.Unpark(t2.ID)
	fmt.Printf("  CAR-1 parked 2h30m, fee = %.0f\n", fee)
	_, err = lot.Unpark(t2.ID)
	fmt.Println("  unpark twice:", err)
}
