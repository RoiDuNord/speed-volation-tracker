package dog

import (
	"errors"
	"math/rand"
	"time"
)

var ErrHasNoConn = errors.New("dog has no connection")

type dog struct {
	connected bool
	id        int
}

func (d *dog) upd() {
	d.id++
}

func sleep() {
	dura := rand.Intn(10)*100 + 1000
	time.Sleep(time.Duration(dura) * time.Millisecond)
}
