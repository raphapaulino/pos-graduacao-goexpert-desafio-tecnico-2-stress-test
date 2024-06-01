package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	url         string
	requests    int
	concurrency int
)

func init() {
	flag.StringVar(&url, "url", "", "URL do serviço a ser testado")
	flag.IntVar(&requests, "requests", 100, "Número total de requests")
	flag.IntVar(&concurrency, "concurrency", 10, "Número de chamadas simultâneas")
}

func main() {
	flag.Parse()

	if url == "" {
		fmt.Println("O parâmetro --url é obrigatório")
		return
	}

	if requests <= 0 {
		fmt.Println("O parâmetro --requests é obrigatório e deve ser maior que 0")
		return
	}

	if concurrency <= 0 {
		fmt.Println("O parâmetro --concurrency é obrigatório e deve ser maior que 0")
		return
	}

	fmt.Printf("Testando %s com %d requests e concorrência de %d...\n", url, requests, concurrency)

	start := time.Now()
	statusCodes := make(map[int]int)
	var wg sync.WaitGroup
	var mutex sync.Mutex
	semaphore := make(chan struct{}, concurrency)

	for i := 0; i < requests; i++ {
		wg.Add(1)
		semaphore <- struct{}{}

		go func() {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				mutex.Lock()
				statusCodes[400]++
				mutex.Unlock()
				<-semaphore
				return
			}
			defer resp.Body.Close()

			mutex.Lock()
			statusCodes[resp.StatusCode]++
			mutex.Unlock()

			fmt.Print(".")
			<-semaphore
		}()
	}

	wg.Wait()
	totalTime := time.Since(start)

	fmt.Println("")
	fmt.Println("Relatório de Stress Test")
	fmt.Printf("Tempo Total: %v\n", totalTime)
	fmt.Printf("Requests Totais Realizadas: %d\n", requests)
	fmt.Printf("Concorrência: %d\n", concurrency)

	mutex.Lock()
	for statusCode, count := range statusCodes {
		fmt.Printf("Status %d: %d vezes\n", statusCode, count)
	}
	mutex.Unlock()
}
