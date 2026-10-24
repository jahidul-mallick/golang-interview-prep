// Shows bytes vs runes, range over strings, safe slicing and strings.Builder.
// Run: go run ./01-go-fundamentals/04-strings-runes-bytes
package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func main() {
	s := "héllo, 世界"

	fmt.Println("=== len counts bytes, not characters ===")
	fmt.Println("len:", len(s), "| runes:", utf8.RuneCountInString(s))
	fmt.Println("s[1] is a byte:", s[1], "| not the character é")

	fmt.Println("=== range yields runes with BYTE indexes ===")
	for i, r := range "aé世" {
		fmt.Printf("index=%d rune=%c (%d bytes)\n", i, r, utf8.RuneLen(r))
	}

	fmt.Println("=== reverse needs []rune ===")
	fmt.Println(reverse(s))

	fmt.Println("=== slicing in the middle of a rune breaks it ===")
	fmt.Printf("%q valid=%v\n", s[:2], utf8.ValidString(s[:2]))

	fmt.Println("=== strings are immutable ===")
	b := []byte("go")
	b[0] = 'G' // s[0] = 'G' would not compile
	fmt.Println(string(b))

	fmt.Println("=== concatenation in a loop vs strings.Builder ===")
	var sb strings.Builder
	for i := 0; i < 3; i++ {
		fmt.Fprintf(&sb, "item%d,", i)
	}
	fmt.Println(sb.String())

	fmt.Println("=== common helpers ===")
	fmt.Println(strings.Fields(" a  b c "), strings.Split("a,b,c", ","), strings.TrimSpace("  hi  "))
	fmt.Println(strings.Contains("golang", "lan"), strings.EqualFold("Go", "GO"))
}
