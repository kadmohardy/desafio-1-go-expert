package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type QuoteResponse struct {
	Bid string `json:"bid"`
}

func main() {
	log.Println("Starting client...")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*300)
	defer cancel()
	quote, error := getQuote(ctx)
	if error != nil {
		log.Println("Error when tried to get quote.")
		return
	}
	createFile(quote.Bid)
}

func getQuote(ctx context.Context) (*QuoteResponse, error) {
	log.Println("Sending request to server...")
	request, error := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if error != nil {
		log.Println("Request error...")

		return nil, error
	}
	response, error := http.DefaultClient.Do(request)
	if error != nil {
		return nil, error
	}

	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	if error != nil {
		return nil, error
	}

	var QuoteResponse QuoteResponse
	error = json.Unmarshal(body, &QuoteResponse)
	if error != nil {
		log.Println("Error unmarshal...")

		return nil, error
	}
	select {
	case <-ctx.Done():
		fmt.Println("Request failed. Timeout reached")

	case <-time.After(time.Millisecond * 300):
		fmt.Println("Request successful")
	}
	return &QuoteResponse, nil
}

func createFile(bid string) {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	tamanho, err := f.Write([]byte("Dólar: " + bid))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Arquivo criado com sucesso! Tamanho: %d bytes\n", tamanho)
	f.Close()
}
