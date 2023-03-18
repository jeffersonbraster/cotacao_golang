package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Cotacao struct {
	Usdbrl struct {
	Code string `json:"code"`
	Codein string `json:"codein"`
	Name string `json:"name"`
	High string `json:"high"`
	Low string `json:"low"`
	VarBid string `json:"varBid"`
	PctChange string `json:"pctChange"`
	Bid string `json:"bid"`
	Ask string `json:"ask"`
	Timestamp string `json:"timestamp"`
	Create_date string `json:"create_date"`
	} `json:"USDBRL"`
}

type CotacaoDB struct {
	Code string `json:"code"`
	Codein string `json:"codein"`
	Name string `json:"name"`
	Bid string `json:"bid"`
	Create_date string `json:"create_date"`
	gorm.Model
}

type CotacaoResponse struct {
	Moeda string `json:"Dolar"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", handleCotacao)
	http.ListenAndServe(":8080", mux)
}

func handleCotacao(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cotacao, err := getCotacao()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	insertCotacaoDB(ctx, cotacao)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := &CotacaoResponse{Moeda: cotacao.Usdbrl.Bid}
	json.NewEncoder(w).Encode(resp)
}

func insertCotacaoDB(ctx context.Context, c *Cotacao) {
	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&CotacaoDB{})

	select {
	case <-ctx.Done():
		fmt.Println("Limite de tempo atingido")
		return
	case <-time.After(10 * time.Millisecond):
		db.Create(&CotacaoDB{
			Code: c.Usdbrl.Code,
			Codein: c.Usdbrl.Codein,
			Name: c.Usdbrl.Name,
			Bid: c.Usdbrl.Bid,
			Create_date: c.Usdbrl.Create_date,
		})
	}	
}

func getCotacao() (*Cotacao, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var c Cotacao

	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

