// Shows MongoDB schema thinking with plain structs: embed vs reference,
// and a tiny filter match like find({status: "paid"}).
// Run: go run ./07-backend-ecosystem/07-mongodb-basics-and-schema-design
package main

import (
	"encoding/json"
	"fmt"
)

// Embedded: line items always read with the order, so store them inside it.
type Order struct {
	ID         string `json:"_id"`
	CustomerID string `json:"customer_id"` // reference: customers change on their own
	Status     string `json:"status"`
	Items      []Item `json:"items"`
}

type Item struct {
	SKU string `json:"sku"`
	Qty int    `json:"qty"`
}

// find returns documents whose fields equal every key in filter.
func find(docs []map[string]any, filter map[string]any) []map[string]any {
	var out []map[string]any
	for _, d := range docs {
		match := true
		for k, v := range filter {
			if d[k] != v {
				match = false
				break
			}
		}
		if match {
			out = append(out, d)
		}
	}
	return out
}

func main() {
	fmt.Println("=== one document holds the order and its items ===")
	o := Order{"o1", "c42", "paid", []Item{{"pen", 2}, {"ink", 1}}}
	b, _ := json.Marshal(o)
	fmt.Println(" ", string(b))

	docs := []map[string]any{
		{"_id": "o1", "status": "paid"},
		{"_id": "o2", "status": "new"},
		{"_id": "o3", "status": "paid"},
	}
	fmt.Println("=== find({status: \"paid\"}) ===")
	for _, d := range find(docs, map[string]any{"status": "paid"}) {
		fmt.Println(" ", d["_id"])
	}
	fmt.Println("=== schema is flexible: same collection, different shape ===")
	docs = append(docs, map[string]any{"_id": "o4", "status": "paid", "coupon": "NEW10"})
	fmt.Println("  paid orders now:", len(find(docs, map[string]any{"status": "paid"})))
}
