package cat

import (
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"time"

	"speed_violation_tracker/models"
)

var ErrHasNoConn = errors.New("cat has no connection")

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func (c *Cat) broadcast() {
	defer c.wg.Done()

	for {
		sleep()

		passage := genPass()
		jPass, _ := json.Marshal(passage)
		if rnd.Intn(20) == 0 && len(jPass) > 10 {
			jPass = jPass[10:]
		}

		select {
		case c.dataCh <- Message{b: jPass}:
			continue
		case _, ok := <-c.stopCh:
			if !ok {
				return
			}
		}
	}
}

func genPass() models.Passage {
	count := rnd.Intn(20) + 10
	return models.Passage{
		Track:      genTrack(count),
		LicenseNum: genGRN(),
		Speeds:     genSpeeds(count),
		Classes:    genClasses(count),
		Sides:      genSides(count),
	}
}

func genSpeeds(cnt int) []float64 {
	baseSpeed := 70 + rand.Float64()*30 - 15
	res := make([]float64, cnt)

	for i := range res {
		speed := baseSpeed + rand.Float64()*3 - 1.5
		res[i] = float64(int(speed*4)) / 4
	}

	return res
}

func genClasses(cnt int) []models.VehicleClass {
	res := make([]models.VehicleClass, cnt)
	for i := range res {
		res[i] = models.VehicleClass(rand.Intn(5) - 1)
	}
	return res
}

func genSides(cnt int) []models.VehicleSide {
	res := make([]models.VehicleSide, cnt)
	for i := range res {
		res[i] = models.VehicleSide(rand.Intn(3) - 1)
	}
	return res
}

func genGRN() string {
	runes := []rune("abcdefgh12345678")
	res := make([]rune, 5)

	for i := range res {
		res[i] = runes[rnd.Intn(len(runes))]
	}

	return string(res)
}

func genTrack(cnt int) []models.TPoint {
	res := make([]models.TPoint, cnt)

	k, b := rnd.Float64()*0.2+0.2, rnd.Float64()*20+20

	for i := range res {
		x := float64(i) / float64(cnt-1) * 100
		res[i] = models.TPoint{
			X: math.Round(x*100) / 100,
			Y: math.Round((k*x+b)*100) / 100,
			T: int(time.Now().Unix()) - (cnt - i),
		}
	}

	for i := range cnt {
		j := rnd.Intn(cnt)
		res[i], res[j] = res[j], res[i]
	}

	return res
}

func sleep() {
	dura := rnd.Intn(5)*100 + 500
	time.Sleep(time.Duration(dura) * time.Millisecond)
}
