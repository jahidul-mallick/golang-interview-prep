// Classic string problems.
// Run: go run ./04-dsa/05-strings-problems
package main

import (
	"fmt"
	"strings"
	"unicode"
)

// isAnagram: count letters up for s and down for t. O(n), O(1) space.
func isAnagram(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	var freq [26]int
	for i := range s {
		freq[s[i]-'a']++
		freq[t[i]-'a']--
	}
	return freq == [26]int{}
}

// isPalindrome ignores case and non-alphanumerics. It works on runes, not
// bytes, so multi-byte characters stay whole. Two pointers, O(n) time and space.
func isPalindrome(s string) bool {
	var r []rune
	for _, c := range s {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			r = append(r, unicode.ToLower(c))
		}
	}
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		if r[i] != r[j] {
			return false
		}
	}
	return true
}

// longestPalindrome expands around every center (odd and even). O(n^2).
func longestPalindrome(s string) string {
	start, end := 0, 0
	expand := func(l, r int) {
		for l >= 0 && r < len(s) && s[l] == s[r] {
			l--
			r++
		}
		if r-l-1 > end-start {
			start, end = l+1, r
		}
	}
	for i := range s {
		expand(i, i)   // odd length
		expand(i, i+1) // even length
	}
	return s[start:end]
}

// reverseWords: strings.Fields drops extra spaces, then swap from both ends.
func reverseWords(s string) string {
	w := strings.Fields(s)
	for i, j := 0, len(w)-1; i < j; i, j = i+1, j-1 {
		w[i], w[j] = w[j], w[i]
	}
	return strings.Join(w, " ")
}

func main() {
	fmt.Println("=== valid anagram ===")
	fmt.Println(isAnagram("anagram", "nagaram"), isAnagram("rat", "car"))

	fmt.Println("=== valid palindrome ===")
	fmt.Println(isPalindrome("A man, a plan, a canal: Panama"), isPalindrome("hello"))

	fmt.Println("=== longest palindromic substring ===")
	fmt.Println(longestPalindrome("babad"), longestPalindrome("cbbd")) // bab bb

	fmt.Println("=== reverse words ===")
	fmt.Printf("%q\n", reverseWords("  the sky  is blue "))
}
