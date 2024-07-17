package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func getFromAPI(url string) (string, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	return string(body), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <cep>")
		return
	}

	cep := os.Args[1]
	api1Chan := make(chan string)
	api2Chan := make(chan string)

	go func() {
		result, err := getFromAPI(fmt.Sprintf("https://viacep.com.br/ws/%s/json", cep))
		if err != nil {
			panic(err)
		}
		api1Chan <- result
	}()

	go func() {
		result, err := getFromAPI("https://brasilapi.com.br/api/cep/v1/" + cep)
		if err != nil {
			panic(err)
		}
		api2Chan <- result
	}()

	select {
	case result := <-api1Chan:
		fmt.Println(result)

	case result := <-api2Chan:
		fmt.Println(result)

	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
		os.Exit(1)
	}
}
