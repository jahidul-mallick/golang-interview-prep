// Shows sentinel errors, wrapping, errors.Is/As/Join, the typed-nil trap and
// boundary mapping.
// Run: go run ./01-go-fundamentals/09-error-handling-wrapping
package main

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrNotFound = errors.New("not found")

// ValidationError is a custom type: callers extract it with errors.As.
type ValidationError struct{ Field string }

func (e *ValidationError) Error() string { return "invalid field " + e.Field }

func findUser(id int) error {
	switch {
	case id < 0:
		return &ValidationError{Field: "id"}
	case id > 100:
		return ErrNotFound
	}
	return nil
}

func loadProfile(id int) error {
	if err := findUser(id); err != nil {
		return fmt.Errorf("load profile %d: %w", id, err)
	}
	return nil
}

func buggy() error {
	var ve *ValidationError
	return ve // non-nil error: the interface has a type
}

// httpStatus is the one place that maps domain errors to codes.
func httpStatus(err error) int {
	var ve *ValidationError
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.As(err, &ve):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func main() {
	fmt.Println("=== wrapping with w keeps the cause, with v loses it ===")
	wrapped := fmt.Errorf("a: %w", ErrNotFound)
	opaque := fmt.Errorf("a: %v", ErrNotFound)
	fmt.Println("wrapped:", errors.Is(wrapped, ErrNotFound), "opaque:", errors.Is(opaque, ErrNotFound))

	fmt.Println("=== Is and As walk the chain ===")
	err := loadProfile(999)
	fmt.Println(err, "| Is NotFound:", errors.Is(err, ErrNotFound))
	err = loadProfile(-1)
	var ve *ValidationError
	if errors.As(err, &ve) {
		fmt.Println(err, "| bad field:", ve.Field)
	}

	fmt.Println("=== errors.Join builds a tree ===")
	joined := errors.Join(ErrNotFound, &ValidationError{Field: "email"})
	fmt.Println("Is:", errors.Is(joined, ErrNotFound), "| Unwrap == nil:", errors.Unwrap(joined) == nil)

	fmt.Println("=== typed-nil trap ===")
	fmt.Println("buggy() == nil:", buggy() == nil)

	fmt.Println("=== boundary mapping ===")
	for _, e := range []error{nil, loadProfile(999), loadProfile(-1), errors.New("db down")} {
		fmt.Printf("%-40v -> %d\n", e, httpStatus(e))
	}
}
