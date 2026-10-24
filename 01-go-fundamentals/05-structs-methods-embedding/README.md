# 05 · Structs, Methods & Embedding

[← Section index](../README.md) · Prev: [04 · Strings, Runes & Bytes](../04-strings-runes-bytes/README.md) · Next: [06 · Interfaces, Type Assertion & Type Switch](../06-interfaces-type-assertion-type-switch/README.md)

**Interview weight: 1.1%**  ·  Revise in: 6 min  ·  Run: `go run ./01-go-fundamentals/05-structs-methods-embedding`

## In one line
Go has no classes: structs hold data, methods attach behaviour to any named type, and embedding gives composition with promoted fields and methods instead of inheritance.

## Must know
- A struct is a value type. Assigning or passing it copies all fields.
- Unset fields get zero values. Use `NewX(...)` functions as constructors; make the zero value useful where possible.
- A method is a function with a receiver: `func (d Dog) Speak()`. Receivers can be value or pointer (see topic 07).
- Embedding `Animal` inside `Dog` promotes its fields and methods: `d.Name`, `d.Speak()`.
- The outer type can override a promoted method. The embedded one is still reachable as `d.Animal.Speak()`.
- Embedding is NOT inheritance: the embedded type never sees the outer type, so there is no virtual dispatch.
- Embedding an interface in a struct satisfies that interface by delegation (handy for partial mocks).
- A struct is comparable (`==`, map key) only if all its fields are comparable.

## Code to remember
```go
type Animal struct{ Name string }

func (a Animal) Speak() string { return a.Name + " makes a sound" }

type Dog struct {
	Animal // embedded, fields and methods promoted
	Breed string
}

func (d Dog) Speak() string { return d.Name + " barks" } // overrides
```

## Common traps
- Expecting polymorphism through embedding: `Animal`'s methods call `Animal`'s `Speak`, never `Dog`'s.
- Two embedded types with the same method name make the call ambiguous (compile error when called).
- Comparing structs that contain slices or maps (compile error).
- Copying a struct that holds a `sync.Mutex` (use a pointer; `go vet` flags this).

## Interview Q&A

**Q1. How does Go do inheritance?**

**A:** It does not. It uses composition: embed a type to reuse its fields and methods, and use interfaces for polymorphism.

**Q2. What happens if the outer struct and the embedded struct both define a method?**

**A:** The outer one wins. You can still call the embedded one explicitly with `outer.Embedded.Method()`.

**Q3. How do you create an object with required setup in Go?**

**A:** A `NewX` function that returns `*X` (and an error if needed). Optionally use functional options for many settings.

**Q4. When are two structs equal?**

**A:** When every field is equal, and only if all fields are comparable types.

**Q5. Can I pass a `Dog` where an `Animal` is expected?**

**A:** No, `Dog` is a different type, not a subtype. Use an interface (for example `Speaker`): `Dog` satisfies it through its own or promoted methods.
