package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gabrielpgava/stress-test-fullcycle/internal/storage"
)

type Request struct {
	URL        string
	StatusCode int
}

func main() {

	url := flag.String("url", "", "A URL para testar a conexão HTTP.")
	requests := flag.Int("requests", 1, "Número de requisições HTTP a serem feitas.")
	concurrency := flag.Int("concurrency", 1, "Número de requisições concorrentes.")

	flag.Parse()

	startTime := time.Now()

	done := make(chan bool)
	semaphore := make(chan struct{}, *concurrency)

	for i := 0; i < *requests; i++ {
		go func() {
			semaphore <- struct{}{}
			makeRequest(*url, done)
			<-semaphore
		}()
	}

	for i := 0; i < *requests; i++ {
		<-done
	}
	elapsedTime := time.Since(startTime)

	//Relatorio
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("                   RELATÓRIO FINAL")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("URL fornecida:                    %s\n", *url)
	fmt.Printf("Número de requisições:            %d\n", *requests)
	fmt.Printf("Número de requisições concorrentes: %d\n", *concurrency)
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Tempo total de execução:          %s\n", elapsedTime)
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("Distribuição de códigos de status:")
	totalProcessed := storage.PrintReport()
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Total de requisições processadas: %d\n", totalProcessed)
	fmt.Println(strings.Repeat("=", 50) + "\n")

}

func makeRequest(url string, done chan bool) {
	defer func() { done <- true }()

	resp, err := http.Get(url)
	if err != nil {
		storage.IncrementCode("Error", 0)
		return
	}
	defer resp.Body.Close()

	storage.IncrementCode(fmt.Sprintf("%d", resp.StatusCode), resp.StatusCode)
}
