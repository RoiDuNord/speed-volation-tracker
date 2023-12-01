package cat

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/kvolis/tesgode/models"
)

var ErrHasNoConn = errors.New("cat has no connection")

func init() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("seeded")
}

type cat struct {
	connected bool
	dataCh    CatChan
	stopCh    chan struct{}
	wg        sync.WaitGroup
}

type message struct {
	b []byte
}

func (c *cat) broadcast() {
	defer c.wg.Done()

	for {
		sleep()

		passage := genPass()
		jPass, _ := json.Marshal(passage)

		select {
		case c.dataCh <- message{b: jPass}:
			continue
		case _, ok := <-c.stopCh:
			if !ok {
				return
			}
		}
	}
}

func genPass() models.Passage {
	return models.Passage{
		Track:      genTrack(),
		LicenseNum: genGRN(),
	}
}

func genGRN() string {
	runes := []rune("abcdefgh12345678")
	res := make([]rune, 5)

	for i := range res {
		res[i] = runes[rand.Intn(len(runes))]
	}

	return string(res)
}

func genTrack() []models.TPoint {
	cnt := rand.Intn(20) + 10
	res := make([]models.TPoint, cnt)

	k, b := rand.Float64()*0.2+0.2, rand.Float64()*20+20

	for i := range res {
		x := float64(i) / float64(cnt-1) * 100
		res[i] = models.TPoint{
			X: x,
			Y: k*x + b,
			T: int(time.Now().Unix()),
		}
	}

	for i := 0; i < cnt; i++ {
		j := rand.Intn(cnt)
		res[i], res[j] = res[j], res[i]
	}

	return res
}

func sleep() {
	dura := rand.Intn(5)*100 + 500
	time.Sleep(time.Duration(dura) * time.Millisecond)
}
