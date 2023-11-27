package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Exchange sends request and receives response.
// It returns with the response body, with the response status code and with the error.
func Exchange(request *http.Request) (string, int, error) {
	client := &http.Client{}

	start := time.Now()
	response, err := client.Do(request)
	//	measuring the exchange
	duration := time.Since(start)
	fmt.Println("Duration:\n   ", duration)

	if err != nil {
		fmt.Println(err)
		return "", 0, err
	}
	defer func(body io.ReadCloser) {
		err2 := body.Close()
		if err2 != nil {
			fmt.Println("Failed closing response body.")
		}
	}(response.Body)

	// Extract response body
	var buf bytes.Buffer
	n, err := buf.ReadFrom(response.Body)
	if n < 1 {
		fmt.Println("Body was empty. Something went wrong.")
	}
	if err != nil {
		fmt.Println(err)
		return "", 0, err
	}

	fmt.Println("\nResponse body: ", buf.String())
	fmt.Println("Response status: ", response.Status)
	logHeaders(response.Header)

	return buf.String(), response.StatusCode, nil
}
