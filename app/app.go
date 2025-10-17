package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log" // Добавлено для логирования
	"sync"
	"time"

	"speed_violation_tracker/cat"
	"speed_violation_tracker/config"
	"speed_violation_tracker/dog"
	"speed_violation_tracker/models"
)

const (
	MAX_MESSAGES       = 7
	SECONDS_IN_MINUTES = 60
)

func MustRun() error {
	var errorsSlice []error

	cfg, err := config.MustLoad()
	if err != nil {
		errorsSlice = append(errorsSlice, fmt.Errorf("config load failed: %w", err))
	}
	log.Printf("INFO config loaded host=%s port=%d", cfg.Server.Host, cfg.Server.Port)

	db := dog.New()
	dbConn := "postgres"
	if err := db.Connect(dbConn); err != nil {
		errorsSlice = append(errorsSlice, fmt.Errorf("database connect failed: %w", err))
	}

	messageBroker := cat.New()
	mbConn := "kafka"
	if err := messageBroker.Connect(mbConn); err != nil {
		errorsSlice = append(errorsSlice, fmt.Errorf("message broker connect failed: %w", err))
	}

	ch, err := messageBroker.Subscript()
	if err != nil {
		errorsSlice = append(errorsSlice, fmt.Errorf("message broker subscription failed: %w", err))
	}

	start := time.Now()
	var wg sync.WaitGroup
	var counter int
	var mu sync.Mutex

	var passages []models.Passage

	// Первый цикл: Получаем и размаршаллим сообщения параллельно
	for message := range ch {
		wg.Add(1)
		go func(msg cat.Message) {
			defer wg.Done()
			var passage models.Passage
			if err := json.Unmarshal(msg.Bytes(), &passage); err != nil {
				mu.Lock()
				// Добавьте детальный лог для отладки invalid JSON (уберите в продакшене)
				log.Printf("ERROR: Failed to unmarshal message %s: %v", string(msg.Bytes()), err)
				errorsSlice = append(errorsSlice, fmt.Errorf("invalid JSON in message: %w", err))
				mu.Unlock()
				return
			}
			mu.Lock()
			passages = append(passages, passage)
			mu.Unlock()
		}(message)
		counter++
		if counter >= MAX_MESSAGES {
			break
		}
	}

	// Ждём завершения всех горутин-потребителей, чтобы избежать race на срезах
	wg.Wait()

	// Второй цикл: Обрабатываем passages параллельно с безопасными захватами
	for i, passage := range passages {
		wg.Add(1)
		go func(p models.Passage, idx int) {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()

			// Проверка границ: Пропускаем, если Track пустой
			if len(p.Track) == 0 {
				errorsSlice = append(errorsSlice, fmt.Errorf("passage at index %d has empty Track; skipping: %w", idx, errors.New("processing skipped")))
				return
			}

			maxTimeStamp := p.Track[0].T
			for _, point := range p.Track[1:] {
				if point.T > maxTimeStamp {
					maxTimeStamp = point.T
				}
			}

			if maxTimeStamp%SECONDS_IN_MINUTES >= 45 {
				pasBytes, err := json.Marshal(p)
				if err != nil {
					errorsSlice = append(errorsSlice, fmt.Errorf("failed to marshal passage at index %d: %w", idx, err))
					return
				}
				db.Insert(fmt.Sprintf("offender_%d", idx), pasBytes)
			}
		}(passage, i)
	}

	wg.Wait() // Ждём завершения всех обработчиков
	fmt.Println("Execution time:", time.Since(start))

	if len(errorsSlice) > 0 {
		// Логируйте итоговые ошибки перед return
		log.Printf("ERROR: Encountered %d errors during execution", len(errorsSlice))
		for _, err := range errorsSlice {
			log.Printf(" - %v", err)
		}
		return errors.Join(errorsSlice...)
	}

	if err := gracefulShutdown(messageBroker, db); err != nil {
		log.Printf("ERROR: Graceful shutdown failed: %v", err)
		return err
	}

	return nil
}

func gracefulShutdown(broker *cat.Cat, db *dog.Dog) error {
	var errorsSlice []error

	if err := broker.Close(); err != nil {
		errorsSlice = append(errorsSlice, fmt.Errorf("ошибка закрытия брокера: %w", err))
	}

	if err := db.Close(); err != nil {
		errorsSlice = append(errorsSlice, fmt.Errorf("ошибка закрытия БД: %w", err))
	}

	if len(errorsSlice) > 0 {
		return errors.Join(errorsSlice...)
	}
	return nil
}
