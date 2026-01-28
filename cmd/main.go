package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gabrielpgava/stress-test-fullcycle/internal/storage"
)

func main() {

	url := flag.String("url", "", "A URL para testar a conexão HTTP.")
	requests := flag.Int("requests", 1, "Número de requisições HTTP a serem feitas.")
	concurrency := flag.Int("concurrency", 1, "Número de requisições concorrentes.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: %s --url=http://exemplo.com --requests=100 --concurrency=10\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if err := validateFlags(*url, *requests, *concurrency); err != nil {
		fmt.Fprintln(os.Stderr, "Erro:", err)
		flag.Usage()
		os.Exit(1)
	}

	storage.Reset()

	startTime := time.Now()
	client := &http.Client{Timeout: 10 * time.Second}

	runLoad(*url, *requests, *concurrency, client)
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
	fmt.Printf("Total de respostas 200:           %d\n", storage.GetStatusCounts()[200])
	fmt.Printf("Total de erros de conexão:        %d\n", storage.GetErrorCount())
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("Distribuição de códigos de status:")
	totalProcessed := storage.PrintReport()
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Total de requisições processadas: %d\n", totalProcessed)
	fmt.Println(strings.Repeat("=", 50) + "\n")

}

func validateFlags(url string, requests int, concurrency int) error {
	if strings.TrimSpace(url) == "" {
		return fmt.Errorf("a flag --url é obrigatória")
	}
	if requests <= 0 {
		return fmt.Errorf("a flag --requests deve ser maior que zero")
	}
	if concurrency <= 0 {
		return fmt.Errorf("a flag --concurrency deve ser maior que zero")
	}
	if concurrency > requests {
		return fmt.Errorf("a flag --concurrency não pode ser maior que --requests")
	}
	return nil
}

func runLoad(url string, requests int, concurrency int, client *http.Client) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	wg.Add(requests)
	for i := 0; i < requests; i++ {
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			makeRequest(client, url)
			<-semaphore
		}()
	}
	wg.Wait()
}

func makeRequest(client *http.Client, url string) {

	resp, err := client.Get(url)
	if err != nil {
		storage.IncrementError()
		return
	}
	defer resp.Body.Close()

	storage.IncrementStatus(resp.StatusCode)
}
