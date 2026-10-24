// SOLID in Go: OCP, ISP, DIP (+ SRP) with small interfaces.
// Run: go run ./06-lld-design-patterns/01-solid-principles-in-go
package main

import "fmt"

// OCP: a new discount is a new type; Checkout never changes.
type Discount interface{ Apply(total float64) float64 }

type Percent struct{ Pct float64 }
type Flat struct{ Off float64 }

func (p Percent) Apply(t float64) float64 { return t * (1 - p.Pct/100) }
func (f Flat) Apply(t float64) float64    { return max(t-f.Off, 0) }

func Checkout(total float64, d Discount) float64 { return d.Apply(total) }

// ISP: interfaces are small and defined by the consumer.
type Reader interface{ Get(id int) (string, bool) }
type Writer interface{ Put(id int, v string) }

type memStore struct{ m map[int]string }

func (s *memStore) Get(id int) (string, bool) { v, ok := s.m[id]; return v, ok }
func (s *memStore) Put(id int, v string)      { s.m[id] = v }

// Report only needs a Reader, so it cannot write and is easy to fake.
func Report(r Reader, id int) string {
	if v, ok := r.Get(id); ok {
		return "found " + v
	}
	return "missing"
}

// DIP + SRP: OrderService has one job (orders) and depends on an
// interface, so the real notifier can be swapped for a fake.
type Notifier interface{ Notify(msg string) }

type OrderService struct{ notifier Notifier }

func NewOrderService(n Notifier) *OrderService { return &OrderService{notifier: n} }

func (s *OrderService) Place(item string) { s.notifier.Notify("order placed: " + item) }

type emailNotifier struct{}

func (emailNotifier) Notify(msg string) { fmt.Println("  email sent:", msg) }

type fakeNotifier struct{ got []string }

func (f *fakeNotifier) Notify(msg string) { f.got = append(f.got, msg) }

func main() {
	fmt.Println("=== OCP: discounts ===")
	for _, d := range []Discount{Percent{10}, Flat{30}, Flat{500}} {
		fmt.Printf("  %T%v on 200 -> %.0f\n", d, d, Checkout(200, d))
	}

	fmt.Println("=== ISP: Report takes only a Reader ===")
	st := &memStore{m: map[int]string{}}
	var w Writer = st
	w.Put(1, "alice")
	fmt.Println(" ", Report(st, 1), "/", Report(st, 2))

	fmt.Println("=== DIP: inject the dependency ===")
	NewOrderService(emailNotifier{}).Place("book")
	fake := &fakeNotifier{}
	NewOrderService(fake).Place("pen")
	fmt.Println("  fake captured:", fake.got)
}
