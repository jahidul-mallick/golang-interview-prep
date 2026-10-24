// Shows struct basics, constructors, embedding with promotion, overriding and comparison.
// Run: go run ./01-go-fundamentals/05-structs-methods-embedding
package main

import "fmt"

type Animal struct {
	Name string
}

func (a Animal) Speak() string { return a.Name + " makes a sound" }
func (a Animal) Kind() string  { return "animal" }
func (a Animal) Intro() string { return a.Speak() } // always calls Animal.Speak, never Dog's

// Dog embeds Animal: fields and methods are promoted (this is composition, not inheritance).
type Dog struct {
	Animal
	Breed string
}

// Overrides the promoted Speak. The embedded one is still reachable as d.Animal.Speak().
func (d Dog) Speak() string { return d.Name + " barks" }

type Logger interface{ Log(string) }

type ConsoleLogger struct{}

func (ConsoleLogger) Log(msg string) { fmt.Println("log:", msg) }

// Embedding an interface lets a struct satisfy it by delegation.
type Service struct {
	Logger
	Name string
}

type Point struct{ X, Y int }

// Constructor: Go has no constructors, a NewX function is the convention.
func NewDog(name, breed string) *Dog {
	return &Dog{Animal: Animal{Name: name}, Breed: breed}
}

func main() {
	fmt.Println("=== struct literals and zero values ===")
	p := Point{X: 1}
	fmt.Printf("%v %+v\n", p, p)

	fmt.Println("=== embedding: promotion and override ===")
	d := NewDog("Rex", "Lab")
	fmt.Println(d.Name, "|", d.Kind(), "|", d.Speak(), "|", d.Animal.Speak())
	fmt.Println("no virtual dispatch:", d.Intro())

	fmt.Println("=== embedded interface ===")
	svc := Service{Logger: ConsoleLogger{}, Name: "billing"}
	svc.Log("started " + svc.Name)

	fmt.Println("=== structs are comparable if all fields are ===")
	fmt.Println(Point{1, 2} == Point{1, 2})
	visited := map[Point]bool{{1, 2}: true}
	fmt.Println("struct as map key:", visited[Point{1, 2}])

	fmt.Println("=== anonymous struct (handy in tests) ===")
	cfg := struct {
		Host string
		Port int
	}{"localhost", 8080}
	fmt.Println(cfg)
}
