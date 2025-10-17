package main

import (
	"log"
	"speed_violation_tracker/app"
)

func main() {
	if err := app.MustRun(); err != nil {
		log.Fatal(err)
	}
}
