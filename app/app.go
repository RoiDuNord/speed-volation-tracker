package app

import (
	"speed_violation_tracker/config"
)

func MustRun() error {
	_, err := config.MustLoad()
	if err != nil {
		return err
	}

	return err
}
