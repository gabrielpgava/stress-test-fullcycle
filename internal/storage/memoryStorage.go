package storage

import (
	"fmt"
	"sort"
	"sync"
)

var (
	StatusCounts = make(map[int]int)
	ErrorCount   int
	mu           sync.Mutex
)

func Reset() {
	mu.Lock()
	defer mu.Unlock()
	StatusCounts = make(map[int]int)
	ErrorCount = 0
}

func IncrementStatus(code int) {
	mu.Lock()
	defer mu.Unlock()
	StatusCounts[code]++
}

func IncrementError() {
	mu.Lock()
	defer mu.Unlock()
	ErrorCount++
}

func GetStatusCounts() map[int]int {
	mu.Lock()
	defer mu.Unlock()
	copied := make(map[int]int, len(StatusCounts))
	for code, count := range StatusCounts {
		copied[code] = count
	}
	return copied
}

func GetErrorCount() int {
	mu.Lock()
	defer mu.Unlock()
	return ErrorCount
}

func PrintReport() int {
	mu.Lock()
	defer mu.Unlock()
	total := 0
	codes := make([]int, 0, len(StatusCounts))
	for code := range StatusCounts {
		codes = append(codes, code)
	}
	sort.Ints(codes)
	for _, code := range codes {
		count := StatusCounts[code]
		fmt.Printf("Código: %d, Contagem: %d\n", code, count)
		total += count
	}
	if ErrorCount > 0 {
		fmt.Printf("Erros de conexão: %d\n", ErrorCount)
		total += ErrorCount
	}
	return total
}
