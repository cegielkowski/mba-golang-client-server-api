package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:8080/cotacao", nil)
	if err != nil {
		log.Println(err)
		return
	}

	resChan := make(chan *http.Response)
	errChan := make(chan error)
	go func() {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			errChan <- err
		}
		resChan <- res
	}()

	var res *http.Response
	select {
	case <-time.After(300 * time.Millisecond):
		log.Println("Request timeout")
		return
	case err = <-errChan:
		log.Println(err)
		return
	case res = <-resChan:
	}
	defer res.Body.Close()

	respValue, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	content := "Dolar: " + string(respValue)
	err = os.WriteFile("cotacao.txt", []byte(content), 0644)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(content)
}
