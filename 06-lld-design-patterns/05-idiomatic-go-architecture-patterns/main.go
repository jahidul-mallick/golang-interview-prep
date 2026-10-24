// Layered architecture: handler -> service -> repository, wired in main().
// Run: go run ./06-lld-design-patterns/05-idiomatic-go-architecture-patterns
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
)

// Domain: no imports of http or database packages.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var ErrNotFound = errors.New("user not found")

// The service owns this interface (consumer-defined), so it can be faked in tests.
type UserRepo interface {
	FindByID(id int) (*User, error)
}

// Repository layer: one implementation per storage engine.
type memRepo struct {
	mu sync.RWMutex
	m  map[int]*User
}

func (r *memRepo) FindByID(id int) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if u, ok := r.m[id]; ok {
		return u, nil
	}
	return nil, ErrNotFound
}

// Service layer: business rules only.
type UserService struct{ repo UserRepo }

func NewUserService(r UserRepo) *UserService { return &UserService{repo: r} }

func (s *UserService) Get(id int) (*User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}
	return s.repo.FindByID(id)
}

// Handler layer: translates HTTP to service calls and errors to status codes.
type UserHandler struct{ svc *UserService }

func (h UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var id int
	fmt.Sscan(r.PathValue("id"), &id)
	u, err := h.svc.Get(id)
	switch {
	case errors.Is(err, ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case err != nil:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		json.NewEncoder(w).Encode(u)
	}
}

func main() {
	// Composition root: the only place that knows concrete types.
	repo := &memRepo{m: map[int]*User{1: {ID: 1, Name: "asha"}}}
	mux := http.NewServeMux()
	mux.Handle("GET /users/{id}", UserHandler{NewUserService(repo)})

	fmt.Println("=== layered request flow ===")
	for _, path := range []string{"/users/1", "/users/9", "/users/0"} {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		fmt.Printf("  GET %s -> %d %s", path, rec.Code, rec.Body.String())
	}
}
