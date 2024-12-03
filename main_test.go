package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/MaximKlimenko/JavaCode/models"
)

const (
	requestsPerSec = 1200
	testDuration   = 1
	baseURL        = "http://localhost:8080/api/v1/wallet"
)

func TestMain(m *testing.M) {
	var wg sync.WaitGroup
	successCount := 0
	errorCount := 0
	totalReq := requestsPerSec * testDuration
	mu := sync.Mutex{}

	fmt.Printf("Starting load test: %d RPS for %d seconds\n", requestsPerSec, testDuration)

	// Канал для ограничения скорости
	ticker := time.NewTicker(time.Second / time.Duration(requestsPerSec))
	defer ticker.Stop()

	// Время окончания теста
	//endTime := time.Now().Add(time.Duration(testDuration) * time.Second)
	wg.Add(totalReq)
	for i := 0; i < totalReq; i++ {
		<-ticker.C

		go func() {
			defer wg.Done()

			// Создаём тело запроса
			reqBody := models.WalletReq{
				WalletID:      "1111-ac21",
				OperationType: "DEPOSIT", // Или "WITHDRAW"
				Amount:        1,
			}
			body, err := json.Marshal(reqBody)
			if err != nil {
				fmt.Printf("Error marshaling request body: %v\n", err)
				return
			}

			// Отправляем POST-запрос
			resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(body))
			if err != nil {
				mu.Lock()
				errorCount++
				mu.Unlock()
				fmt.Printf("Error: %v\n", err)
				return
			}
			defer resp.Body.Close()

			// Проверяем статус ответа
			if resp.StatusCode == http.StatusOK {
				mu.Lock()
				successCount++
				mu.Unlock()
			} else {
				mu.Lock()
				errorCount++
				mu.Unlock()
				fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
			}
		}()

	}

	// Ждём завершения всех горутин
	wg.Wait()

	// Результаты теста
	fmt.Println("Load test completed:")
	fmt.Printf("  Successful requests: %d\n", successCount)
	fmt.Printf("  Failed requests: %d\n", errorCount)
	fmt.Printf("  Total requests: %d\n", successCount+errorCount)
}
