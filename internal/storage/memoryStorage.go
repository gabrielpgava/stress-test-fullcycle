package storage

import (
	"fmt"
	"sync"

	models "github.com/gabrielpgava/stress-test-fullcycle/internal/model"
)

var (
	StatusCode = make(map[string]*models.StatusCode)
	mu         sync.Mutex
)

func GetCode(statuscode string) (*models.StatusCode, bool) {
	mu.Lock()
	defer mu.Unlock()
	state, exists := StatusCode[statuscode]
	return state, exists
}

func SetCode(statuscode string, state *models.StatusCode) {
	mu.Lock()
	defer mu.Unlock()
	StatusCode[statuscode] = state
}

func IncrementCode(statuscode string, code int) {
	mu.Lock()
	defer mu.Unlock()
	if existing, exists := StatusCode[statuscode]; exists {
		existing.Count++
	} else {
		StatusCode[statuscode] = &models.StatusCode{
			Code:  code,
			Count: 1,
		}
	}
}

func DeleteCode(statuscode string) {
	mu.Lock()
	defer mu.Unlock()
	delete(StatusCode, statuscode)
}

func PrintReport() int {
	mu.Lock()
	defer mu.Unlock()
	total := 0
	for code, status := range StatusCode {
		if code == "Error" {
			fmt.Printf("Erros de conexão: %d\n", status.Count)
		} else {
			fmt.Printf("Código: %s, Contagem: %d\n", code, status.Count)
		}
		total += status.Count
	}
	return total
}
