// Shows implicit interfaces, type assertion, type switch and the typed-nil trap.
// Run: go run ./01-go-fundamentals/06-interfaces-type-assertion-type-switch
package main

import (
	"errors"
	"fmt"
	"math"
)

type Shape interface{ Area() float64 }

type Rect struct{ W, H float64 }
type Circle struct{ R float64 }

func (r Rect) Area() float64    { return r.W * r.H }
func (c *Circle) Area() float64 { return math.Pi * c.R * c.R } // pointer receiver

var _ Shape = (*Circle)(nil) // compile-time check that *Circle implements Shape

type MyErr struct{}

func (*MyErr) Error() string { return "my error" }

// Bug: returns a typed nil pointer wrapped in a non-nil error interface.
func mayFailBad() error {
	var e *MyErr
	return e
}

func mayFailGood() error { return nil }

func describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case int:
		return fmt.Sprintf("int %d", x)
	case string:
		return "string " + x
	case Shape:
		return fmt.Sprintf("shape with area %.1f", x.Area())
	case error:
		return "error " + x.Error()
	default:
		return fmt.Sprintf("other %T", x)
	}
}

func main() {
	fmt.Println("=== implicit satisfaction ===")
	shapes := []Shape{Rect{2, 3}, &Circle{1}} // Circle needs &, its method has a pointer receiver
	for _, s := range shapes {
		fmt.Printf("%T area=%.1f\n", s, s.Area())
	}

	fmt.Println("=== type assertion ===")
	var s Shape = Rect{2, 3}
	r, ok := s.(Rect)
	fmt.Println("Rect:", r, ok)
	c, ok := s.(*Circle) // comma-ok never panics
	fmt.Println("Circle:", c, ok)

	fmt.Println("=== type switch ===")
	for _, v := range []any{nil, 42, "go", Rect{1, 1}, errors.New("boom"), 3.5} {
		fmt.Println(describe(v))
	}

	fmt.Println("=== typed nil trap ===")
	fmt.Println("bad == nil:", mayFailBad() == nil)
	fmt.Println("good == nil:", mayFailGood() == nil)
}
