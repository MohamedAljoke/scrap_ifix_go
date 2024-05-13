package b3

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var url = "https://sistemaswebb3-listados.b3.com.br/indexProxy/indexCall/GetPortfolioDay/eyJsYW5ndWFnZSI6InB0LWJyIiwicGFnZU51bWJlciI6MSwicGFnZVNpemUiOjEyMCwiaW5kZXgiOiJJRklYIiwic2VnbWVudCI6IjEifQ=="

type Response struct {
	Results []Asset
}
type Asset struct {
	Code  string `json:"cod"`
	Asset string `json:"asset"`
	Type  string `json:"type"`
	Part  string `json:"part"`
}

func GetB3IfixData() ([]Asset, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, errors.New("error starting client")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.New("error getting body")
	}
	var parsedResponse Response

	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return nil, errors.New("body must have only a single JSON value")
	}

	return parsedResponse.Results, nil
}
