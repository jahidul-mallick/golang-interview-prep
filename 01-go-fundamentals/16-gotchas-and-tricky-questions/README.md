# 16 · Gotchas & Tricky "What Does This Print?"

[← Section index](../README.md) · Prev: [JSON, Reflection & Struct Tags](../15-json-reflection-struct-tags/README.md)

**Interview weight: 1.2%**  ·  Revise in: 12 min  ·  Run: `go run ./01-go-fundamentals/16-gotchas-and-tricky-questions`

## In one line
Most Go gotchas come from copies, two-word interfaces, shared slice arrays and evaluation timing; name the rule and the answer follows.

## Must know
- An interface is (type, value). A nil `*T` inside an interface is not a nil interface.
- Slices share a backing array. `append` with spare capacity overwrites the original; `s[:n:n]` forces `append` to copy.
- `defer` evaluates its arguments immediately and runs LIFO. Closures read variables when they run.
- `range` gives copies of elements. Modify with `items[i]`, not the loop variable.
- Go 1.22+: each loop iteration has its own loop variable, so closures and goroutines capture the right value.
- `range` over a string yields byte offsets and runes; `len(s)` is bytes.
- Nil map: reads work, writes panic. Nil slice: `append` works. Receive from a closed channel returns zero value and `false`.
- Map iteration order is random. `select` with several ready cases picks one at random.

## Code to remember
```go
a := []int{1, 2, 3}
b := a[:2]
b = append(b, 99)    // a is now [1 2 99]
c := a[:2:2]
c = append(c, 7)     // c is a new array, a unchanged

var p *T
var i any = p        // i != nil
```

## Common traps
- Shadowing `err` with `:=` inside a block, then checking the outer `err`.
- `wg.Add(1)` inside the goroutine instead of before it. (Go 1.25 adds `wg.Go(f)`.)
- Closing a channel twice, or sending on a closed channel: both panic.
- `m[k].f = v` on a map of structs does not compile; copy out, change, store back.

## Interview Q&A

**Q1. What does `var p *T; var i any = p; fmt.Println(i == nil)` print?**

**A:** `false`. The interface holds type `*T` with a nil value, and it is nil only when both are nil.

**Q2. After `b := a[:2]; b = append(b, 99)`, what is `a`?**

**A:** `a[2]` becomes 99 because `b` has spare capacity and shares the array.

**Q3. `for i := 0; i < 3; i++ { defer fmt.Print(i) }` prints?**

**A:** `210`. Arguments are evaluated at the defer statement and defers run last in, first out.

**Q4. What does `for i, r := range "héllo"` give for `i`?**

**A:** `0 1 3 4 5`. `i` is the byte offset and `é` takes two bytes.

**Q5. `func f(s []int) { s[0] = 9; s = append(s, 4) }` is called with `[1 2 3]`. What does the caller see?**

**A:** `[9 2 3]`. The slice header is copied: element writes are shared, the append is not visible.

**Q6. `for _, u := range users { u.Age++ }` changes what?**

**A:** Nothing in `users`. `u` is a copy of each element. Write `users[i].Age++`, or range over pointers.

**Q7. `if ok { v, err := f(); _ = v }; return err` returns what?**

**A:** The outer `err`, usually `nil`. `:=` inside the block declares a new shadowing `err`. Declare `v` first and use `=`.
