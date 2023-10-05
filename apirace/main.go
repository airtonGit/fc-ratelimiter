package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	urlApiCEP = "https://cdn.apicep.com/file/apicep/%s.json" // + cep + ".json
	urlViaCEP = "http://viacep.com.br/ws/%s/json"            // + cep + \"/json/"
)

func main() {
	raceStart()
}

type signal struct{}

func raceStart() {
	client := &http.Client{
		Timeout: time.Second,
	}

	ctxApi, cancelApi := context.WithCancel(context.Background())
	ctxVia, cancelVia := context.WithCancel(context.Background())

	doneC := make(chan signal)

	go func() {
		getCep(ctxApi, doneC, "ApiCEP", client, cancelVia, fmt.Sprintf(urlApiCEP, "88955000"))

	}()

	go func() {
		getCep(ctxVia, doneC, "ViaCEP", client, cancelApi, fmt.Sprintf(urlViaCEP, "88955-000"))
	}()

	<-doneC
}

func getCep(ctx context.Context, doneC chan signal, cepSource string, client *http.Client, cancelFunc context.CancelFunc, url string) {

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("new request: %s\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("request doing: %s\n", err)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("body read: %s\n", err)
		return
	}

	// cancel other request
	cancelFunc()
	println(fmt.Sprintf("%s WINS!", cepSource))
	println(fmt.Sprintf("%s", body))

	doneC <- signal{}
}
