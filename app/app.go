package app

import (
	"errors" // Для wrapping ошибок (fmt.Errorf с %w)
	// Для логирования ошибок закрытия
	"speed_violation_tracker/cat"
	"speed_violation_tracker/config"
	"speed_violation_tracker/dog"
)

func MustRun() error {
	var errorsSlice []error

	_, err := config.MustLoad()
	if err != nil {
		errorsSlice = append(errorsSlice, err)
	}

	db := dog.New()
	dbConn := "postgres"
	if err := db.Connect(dbConn); err != nil {
		errorsSlice = append(errorsSlice, err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			errorsSlice = append(errorsSlice, closeErr)
		}
	}()

	messageBroker := cat.New()
	mbConn := "kafka"
	if err := messageBroker.Connect(mbConn); err != nil {
		errorsSlice = append(errorsSlice, err)
	}
	defer func() {
		if closeErr := messageBroker.Close(); closeErr != nil {
			errorsSlice = append(errorsSlice, closeErr)
		}
	}()

	if len(errorsSlice) > 0 {
		return errors.Join(errorsSlice...)
	}
	return nil
}
