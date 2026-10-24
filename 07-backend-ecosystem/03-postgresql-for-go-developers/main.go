// Shows two PostgreSQL habits without a database: the withTx helper
// (commit on success, rollback on error or panic) and keyset vs offset paging.
// Run: go run ./07-backend-ecosystem/03-postgresql-for-go-developers
package main

import (
	"errors"
	"fmt"
	"sort"
)

// Tx mirrors the few methods of *sql.Tx we care about.
type Tx interface {
	Commit() error
	Rollback() error
}

type fakeTx struct{ name string }

func (t fakeTx) Commit() error   { fmt.Println("  COMMIT", t.name); return nil }
func (t fakeTx) Rollback() error { fmt.Println("  ROLLBACK", t.name); return nil }

// withTx is the pattern to remember: one place owns commit and rollback.
func withTx(name string, fn func() error) (err error) {
	tx := Tx(fakeTx{name}) // real code: tx, err := db.BeginTx(ctx, nil)
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	if err = fn(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

type Row struct {
	ID   int
	Name string
}

// offsetPage = LIMIT n OFFSET k: the database still walks k rows.
func offsetPage(rows []Row, offset, limit int) []Row {
	return rows[offset:min(offset+limit, len(rows))]
}

// keysetPage = WHERE id > $1 ORDER BY id LIMIT n: an index seek, stable under inserts.
func keysetPage(rows []Row, afterID, limit int) []Row {
	i := sort.Search(len(rows), func(i int) bool { return rows[i].ID > afterID })
	return rows[i:min(i+limit, len(rows))]
}

func main() {
	fmt.Println("=== transaction: success ===")
	fmt.Println("  err:", withTx("transfer-1", func() error { return nil }))

	fmt.Println("=== transaction: error rolls back ===")
	err := withTx("transfer-2", func() error { return errors.New("insufficient funds") })
	fmt.Println("  err:", err)

	fmt.Println("=== transaction: panic rolls back, then re-panics ===")
	func() {
		defer func() { fmt.Println("  recovered:", recover()) }()
		withTx("transfer-3", func() error { panic("bug") })
	}()

	rows := []Row{{1, "a"}, {2, "b"}, {3, "c"}, {4, "d"}, {5, "e"}}
	fmt.Println("=== offset paging, page 2 of size 2 ===")
	fmt.Println(" ", offsetPage(rows, 2, 2))
	fmt.Println("=== keyset paging, after id=2, size 2 ===")
	fmt.Println(" ", keysetPage(rows, 2, 2))
}
