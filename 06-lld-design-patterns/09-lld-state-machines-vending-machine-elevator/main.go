// Vending machine as a state machine: a transition table plus per-state actions.
// Run: go run ./06-lld-design-patterns/09-lld-state-machines-vending-machine-elevator
package main

import (
	"fmt"
	"sync"
)

type State string
type Event string

const (
	Idle     State = "Idle"
	HasMoney State = "HasMoney"
	Dispense State = "Dispensing"

	InsertCoin Event = "insert_coin"
	Select     Event = "select"
	Cancel     Event = "cancel"
	Done       Event = "done"
)

// The table is the single source of truth; a missing entry is illegal.
var transitions = map[State]map[Event]State{
	Idle:     {InsertCoin: HasMoney},
	HasMoney: {InsertCoin: HasMoney, Select: Dispense, Cancel: Idle},
	Dispense: {Done: Idle},
}

type Machine struct {
	mu      sync.Mutex // serialises events from concurrent callers
	state   State
	balance int
	price   int
}

func (m *Machine) Handle(e Event, coin int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	to, ok := transitions[m.state][e]
	if !ok {
		return fmt.Errorf("%s not allowed in %s", e, m.state)
	}
	switch e {
	case InsertCoin:
		m.balance += coin
	case Select:
		if m.balance < m.price {
			return fmt.Errorf("need %d more", m.price-m.balance) // guard: stay in state
		}
	case Cancel:
		fmt.Println("    refund", m.balance)
		m.balance = 0
	case Done:
		fmt.Println("    item out, change", m.balance-m.price)
		m.balance = 0
	}
	m.state = to
	return nil
}

func main() {
	fmt.Println("=== vending machine ===")
	m := &Machine{state: Idle, price: 20}
	steps := []struct {
		e    Event
		coin int
	}{
		{Select, 0}, {InsertCoin, 10}, {Select, 0}, {InsertCoin, 15}, {Select, 0}, {Done, 0},
	}
	for _, s := range steps {
		err := m.Handle(s.e, s.coin)
		fmt.Printf("  %-11s -> state=%-10s balance=%d err=%v\n", s.e, m.state, m.balance, err)
	}
}
