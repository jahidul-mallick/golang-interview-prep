# 11 · Generics

[← Section index](../README.md) · Prev: [10 · Panic & Recover](../10-panic-recover/README.md) · Next: [12 · Packages, Modules & Initialization Order](../12-packages-modules-init-order/README.md)

**Interview weight: 0.8%**  ·  Revise in: 8 min  ·  Run: `go run ./01-go-fundamentals/11-generics`

## In one line
Generics (Go 1.18+) let one function or type work over many types, with a constraint saying which operations are allowed.

## Must know
- Type parameters go in square brackets: `func Map[T, U any](...)`, `type Stack[T any] struct{...}`.
- A constraint is an interface; it can list a type set (the allowed types, `~int | ~float64`) and methods.
- `any` allows no operators. `comparable` allows `==`. `cmp.Ordered` allows `<` and `>`.
- `~int` means "any type whose underlying type is int", so `type Celsius float64` fits `~float64`.
- Methods can use the type's parameters but cannot add new ones.
- Type arguments are usually inferred from function arguments, not from the return type.
- The zero value of `T` is `var zero T`.
- Do not use generics when one concrete type or a plain interface is enough.

## Code to remember
```go
type Number interface{ ~int | ~int64 | ~float64 }

func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x
	}
	return total
}

type Stack[T any] struct{ items []T }

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }
```

## Common traps
- Constraint without `~` rejects named types like `type ID int`.
- `Zero[int]()` needs an explicit type argument because `T` appears only in the result.
- `comparable` accepts interface types, but comparing uncomparable dynamic values panics at run time.
- Over-engineering: a function taking an interface is often clearer than a generic one.

## Interview Q&A

**Q1. What is a constraint?**

**A:** An interface that limits which types may be used as a type argument, by a type set, by methods, or both.

**Q2. `any` vs `comparable`?**

**A:** `any` is every type but allows no operators. `comparable` allows `==` and `!=`, so the type can be a map key.

**Q3. When would you not use generics?**

**A:** When behaviour varies, not data shape (use an interface), or when only one type is used today.

**Q4. Can a method have its own type parameters?**

**A:** No. A method can only use the parameters of its receiver type. Write a generic function instead.

**Q5. How are generics implemented?**

**A:** The compiler shares code between types with the same memory shape and passes a dictionary, so it is neither full per-type copying nor boxing everything. Treat the exact cost as an implementation detail.

**Q6. What does `~int` mean in a constraint?**

**A:** It matches every type whose underlying type is `int`, so `type ID int` is accepted. Plain `int` matches only `int` itself.
