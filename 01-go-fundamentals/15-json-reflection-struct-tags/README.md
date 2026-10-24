# JSON, Reflection & Struct Tags

[← Section index](../README.md) · Prev: [Iterators & range-over-func (Go 1.23+)](../14-iterators-range-over-func/README.md) · Next: [16 · Gotchas & Tricky "What Does This Print?"](../16-gotchas-and-tricky-questions/README.md)

**Interview weight: 0.6%**  ·  Revise in: 8 min  ·  Run: `go run ./01-go-fundamentals/15-json-reflection-struct-tags`

## In one line
`encoding/json` uses reflection to read struct tags, so you control field names and omission with tags, and only exported fields are encoded.

## Must know
- Tags are raw strings: `` `json:"name,omitempty"` ``; `json:"-"` skips the field.
- Unexported fields are always ignored, tag or not.
- `omitempty` omits false, 0, nil pointer/interface, and empty slice, map or string. It does not omit a zero struct; Go 1.24 `omitzero` does (handy for `time.Time`).
- Use a pointer field (`*int`) to tell "missing" from "zero" when decoding.
- Unknown JSON fields are ignored by default; `Decoder.DisallowUnknownFields()` makes them an error.
- `json.Marshal` returns `[]byte`; `json.NewEncoder(w)` and `json.NewDecoder(r)` stream to a writer or from a reader.
- Decoding into `any` gives `map[string]any`, and numbers become `float64`.
- Reflection: `reflect.TypeOf` for the type, `reflect.ValueOf` for the value, `Kind()` is the underlying category. It is slow and loses compile-time safety.

## Code to remember
```go
type User struct {
	ID    int    `json:"id"`
	Email string `json:"email,omitempty"`
	Pass  string `json:"-"`
}

dec := json.NewDecoder(r)
dec.DisallowUnknownFields()
err := dec.Decode(&u)
```

## Common traps
- Lowercase field names: silently skipped by JSON.
- Numbers decoded into `any` are `float64`, which loses precision for big ints.
- `omitempty` on an int hides a legitimate `0`. Use `*int`.
- `Decode` on a non-pointer target fails; always pass `&v`.

## Interview Q&A

**Q1. `Marshal` vs `Encoder`?**

**A:** `Marshal` builds the whole `[]byte` in memory. `Encoder` writes to an `io.Writer`, better for HTTP responses and large data.

**Q2. How do you distinguish a missing field from a zero value?**

**A:** Use a pointer type, or `json.RawMessage` when you need full control.

**Q3. How do you customize encoding of a type?**

**A:** Implement `MarshalJSON` and `UnmarshalJSON`, or `encoding.TextMarshaler`.

**Q4. What is reflection and why avoid it in hot paths?**

**A:** Package `reflect` inspects types and values at run time; `encoding/json` and ORMs use it. It is slower than direct access and errors appear at run time. Measure first.

**Q5. A field I set is missing from the JSON output. Why?**

**A:** It is unexported (lowercase), tagged `json:"-"`, or hidden by `omitempty` because it holds a zero value.
