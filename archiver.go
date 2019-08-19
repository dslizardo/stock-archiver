package archiver

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const basedURL = "http://www.pse.com.ph/stockMarket/home.html"

// Stock is a struct representation of individual stock information
type Stock struct {
	SecurityAlias   string `json:"securityAlias"`
	SecuritySymbol  string `json:"securitySymbol"`
	LastTradedPrice string `json:"lastTradedPrice"`
}

func getStocks() []Stock {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", basedURL, nil)
	req.Header.Set("Referer", basedURL)

	params := req.URL.Query()
	params.Set("method", "getSecuritiesAndIndicesForPublic")
	params.Set("ajax", "true")
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	stocks := []Stock{}
	if err := json.Unmarshal(body, &stocks); err != nil {
		log.Fatal(err)
	}

	return stocks
}
