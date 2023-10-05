package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Quote struct {
	Data QuoteRequest `json:"USDBRL"`
}

type QuoteRequest struct {
	Code      string `json:"code"`
	Codein    string `json:"codein"`
	Name      string `json:"name"`
	High      string `json:"high"`
	Low       string `json:"low"`
	VarBid    string `json:"varBid"`
	PctChange string `json:"pctChange"`
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	Timestamp string `json:"timestamp"`
	gorm.Model
}

type QuoteResponse struct {
	Bid string `json:"bid"`
}

func main() {
	log.Println("Starting server...")
	initDB()
	http.HandleFunc("/cotacao", getQuoteHandler)
	http.ListenAndServe(":8080", nil)
}

func getQuoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()
	quoteRequest, error := getQuoteOnCurrencyAPI(ctx)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctxDB := context.Background()
	ctxDB, cancelDB := context.WithTimeout(ctxDB, time.Millisecond*10)
	defer cancelDB()
	insertQuoteRequest(ctxDB, quoteRequest)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := QuoteResponse{Bid: quoteRequest.Bid}
	json.NewEncoder(w).Encode(resp)
}

func initDB() {
	db, err := gorm.Open(sqlite.Open("quotes.db"), &gorm.Config{})
	if err != nil {
		log.Println("Error while initializing database")
		panic(err)
	}
	db.AutoMigrate(&QuoteRequest{})
	log.Println("DB initialized")
	DB = db
}

func insertQuoteRequest(ctx context.Context, quoteRequest *QuoteRequest) {
	currencyAPIContext, _ := context.WithTimeout(ctx, time.Millisecond*20)
	DB.WithContext(currencyAPIContext)
	DB.Create(quoteRequest)
}

func getQuoteOnCurrencyAPI(ctx context.Context) (*QuoteRequest, error) {
	log.Println("Calling prices API...")
	currencyAPIContext, cancel := context.WithTimeout(ctx, time.Millisecond*1000)
	defer cancel()
	request, error := http.NewRequestWithContext(currencyAPIContext, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if error != nil {
		return nil, error
	}
	response, error := http.DefaultClient.Do(request)
	if error != nil {
		return nil, error
	}
	defer response.Body.Close()
	body, error := io.ReadAll(response.Body)
	if error != nil {
		return nil, error
	}
	var Quote Quote
	error = json.Unmarshal(body, &Quote)
	if error != nil {
		log.Println("Error unmarshal...")
		return nil, error
	}
	return &Quote.Data, nil
}
