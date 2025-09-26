package dog

import (
	"errors"
	"math/rand"
	"time"
)

var (
	ErrHasNoConn = errors.New("dog has no connection")
	ErrInternal  = errors.New("failed new insertion: unexpected internal error")
)

func (d *Dog) upd() {
	d.id.Add(1)
}

func sleep() {
	dura := rand.Intn(10)*100 + 1000
	time.Sleep(time.Duration(dura) * time.Millisecond)
}
