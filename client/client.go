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

	io.Copy(os.Stdout, res.Body)

	if err != nil {
		panic(err)
	}

	res1, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res1.Body.Close()

	body, err := io.ReadAll(res1.Body)

	if err != nil {
		panic(err)
	}
	
	var c CotacaoResponse

	_ = json.Unmarshal(body, &c)

	file, err := os.Create("cotacao.txt")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dolar: %v", c.Dolar))

	if err != nil {
		panic(err)
	}
}