package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type CotacaoResponse struct {
	Dolar string `json:"Dolar"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}
	
	var c CotacaoResponse

	err = json.Unmarshal(body, &c)

	if err != nil {
		panic(err)
	}

	file, err := os.Create("cotacao.txt")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	file.WriteString(fmt.Sprintf("Dolar: %v", c.Dolar))
}