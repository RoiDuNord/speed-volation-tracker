// Package cat simulates the operation of a message broker.
// Minimal functionality implemented.
package cat

type CatChan chan message

// New returns new cat-brocker entity
func New() *cat {
	return new(cat)
}

// Connect connect to server using conn
func (c *cat) Connect(conn string) error {
	c.connected = true
	return nil
}

// Subscript returns CatChan and start broadcast
func (c *cat) Subscipt() (CatChan, error) {
	if !c.connected {
		return nil, ErrHasNoConn
	}

	c.dataCh = make(CatChan)
	c.stopCh = make(chan struct{})

	c.wg.Add(1)
	go c.broadcast()

	return c.dataCh, nil
}

// Close close cat connection
func (c *cat) Close() error {
	c.connected = false

	close(c.stopCh)
	c.wg.Wait()

	c.dataCh = nil
	return nil
}

// Bytes returns a bytes body message
func (m *message) Bytes() []byte {
	return m.b
}
