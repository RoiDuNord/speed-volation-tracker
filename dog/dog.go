// Package dog simulates the operation of a database.
// Minimal functionality implemented.
package dog

import (
	"fmt"
	"math/rand"
	"sync/atomic"
)

type Dog struct {
	connected bool
	id        atomic.Int32
}

// New returns new dog-db entity
func New() *Dog {
	return new(Dog)
}

// Connect connects to server using conn
func (d *Dog) Connect(conn string) error {
	d.connected = true
	return nil
}

// Insert inserts new entry using key and value, returns id that entry and error if any
func (d *Dog) Insert(key string, value []byte) (int, error) {
	if !d.connected {
		return -1, ErrHasNoConn
	}

	if rand.Intn(10) == 0 {
		return -1, ErrInternal
	}

	id := int(d.id.Load())
	fmt.Printf("new db entry; id: %d; key: <%s>; data len: %d bytes\n", id, key, len(value))

	d.upd()
	sleep()

	return id, nil
}

// Close close cat connection
func (d *Dog) Close() error {
	d.connected = false
	return nil
}
