// Package dog simulates the operation of a database.
// Minimal functionality implemented.
package dog

import "fmt"

// New returns new dog-db entity
func New() *dog {
	return new(dog)
}

// Connect connect to server using conn
func (d *dog) Connect(conn string) error {
	d.connected = true
	return nil
}

// Insert inserts new entry using key and value, returns id that entry and error if any
func (d *dog) Insert(key string, value []byte) (int, error) {
	if !d.connected {
		return -1, ErrHasNoConn
	}

	id := d.id
	fmt.Printf("new db entry; id: %d; key: <%s>; data len: %d bytes\n", id, key, len(value))

	d.upd()
	sleep()

	return id, nil
}

// Close close cat connection
func (d *dog) Close() error {
	d.connected = false
	return nil
}
