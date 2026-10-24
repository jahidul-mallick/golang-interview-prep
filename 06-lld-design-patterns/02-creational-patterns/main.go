// Creational patterns: singleton (sync.Once), factory, functional options, builder.
// Run: go run ./06-lld-design-patterns/02-creational-patterns
package main

import (
	"fmt"
	"strings"
	"sync"
)

// Singleton: sync.Once makes lazy init safe across goroutines.
type Config struct{ Env string }

var (
	cfg  *Config
	once sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		fmt.Println("  (init runs once)")
		cfg = &Config{Env: "prod"}
	})
	return cfg
}

// Factory: callers get an interface and never name the concrete type.
type Notifier interface{ Send(msg string) string }

type email struct{}
type sms struct{}

func (email) Send(m string) string { return "email: " + m }
func (sms) Send(m string) string   { return "sms: " + m }

func NewNotifier(kind string) (Notifier, error) {
	switch kind {
	case "email":
		return email{}, nil
	case "sms":
		return sms{}, nil
	}
	return nil, fmt.Errorf("unknown notifier %q", kind)
}

// Functional options: defaults plus optional overrides, no giant constructor.
type Server struct {
	host    string
	port    int
	tlsOn   bool
	timeout int
}

type Option func(*Server)

func WithPort(p int) Option { return func(s *Server) { s.port = p } }
func WithTLS() Option       { return func(s *Server) { s.tlsOn = true } }

func NewServer(host string, opts ...Option) *Server {
	s := &Server{host: host, port: 80, timeout: 30}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Builder: step-by-step construction, validated in Build().
type QueryBuilder struct {
	table string
	where []string
}

func From(t string) *QueryBuilder { return &QueryBuilder{table: t} }

func (b *QueryBuilder) Where(c string) *QueryBuilder {
	b.where = append(b.where, c)
	return b
}

func (b *QueryBuilder) Build() (string, error) {
	if b.table == "" {
		return "", fmt.Errorf("table required")
	}
	q := "SELECT * FROM " + b.table
	if len(b.where) > 0 {
		q += " WHERE " + strings.Join(b.where, " AND ")
	}
	return q, nil
}

func main() {
	fmt.Println("=== singleton ===")
	a, b := GetConfig(), GetConfig()
	fmt.Println("  same pointer:", a == b)

	fmt.Println("=== factory ===")
	for _, k := range []string{"email", "sms", "pager"} {
		n, err := NewNotifier(k)
		if err != nil {
			fmt.Println(" ", err)
			continue
		}
		fmt.Println(" ", n.Send("hi"))
	}

	fmt.Println("=== functional options ===")
	fmt.Printf("  %+v\n", *NewServer("localhost", WithPort(8443), WithTLS()))

	fmt.Println("=== builder ===")
	q, _ := From("users").Where("age > 18").Where("active = true").Build()
	fmt.Println(" ", q)
}
