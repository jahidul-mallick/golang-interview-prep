// Shows JSON struct tags, strict decoding and reading tags with reflect.
// Run: go run ./01-go-fundamentals/15-json-reflection-struct-tags
package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	Age   *int   `json:"age,omitempty"` // pointer: nil means "not sent"
	Pass  string `json:"-"`
	note  string // unexported: never encoded
}

func main() {
	fmt.Println("=== Marshal: tags, omitempty, skipped fields ===")
	b, _ := json.Marshal(User{ID: 1, Name: "Alice", Pass: "secret", note: "x"})
	fmt.Println(string(b))

	fmt.Println("=== Decode: unknown fields ignored by default ===")
	var u User
	err := json.NewDecoder(strings.NewReader(`{"id":2,"name":"Bob","age":0,"extra":1}`)).Decode(&u)
	fmt.Println(err, "| age is set:", u.Age != nil)

	fmt.Println("=== DisallowUnknownFields makes it strict ===")
	dec := json.NewDecoder(strings.NewReader(`{"id":3,"extra":1}`))
	dec.DisallowUnknownFields()
	fmt.Println(dec.Decode(&User{}))

	fmt.Println("=== numbers decoded into any are float64 ===")
	var m map[string]any
	_ = json.Unmarshal([]byte(`{"n":42}`), &m)
	fmt.Printf("%T\n", m["n"])

	fmt.Println("=== reading tags with reflect ===")
	t := reflect.TypeOf(User{})
	for i := range t.NumField() {
		f := t.Field(i)
		fmt.Printf("%-6s kind=%-7s json=%q\n", f.Name, f.Type.Kind(), f.Tag.Get("json"))
	}
}
