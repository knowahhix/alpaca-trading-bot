package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)
type SymbolChange struct {
	Symbol string
	Change float64
}

type Item struct {
	ID                   string   `json:"id"`
	Class                string   `json:"class"`
	Exchange             string   `json:"exchange"`
	Symbol               string   `json:"symbol"`
	Name                 string   `json:"name"`
	Status               string   `json:"status"`
	Tradable             bool     `json:"tradable"`
	Marginable           bool     `json:"marginable"`
	Shortable            bool     `json:"shortable"`
	EasyToBorrow         bool     `json:"easy_to_borrow"`
	Fractionable         bool     `json:"fractionable"`
	MarginRequirementLong string   `json:"margin_requirement_long"`
	MarginRequirementShort string   `json:"margin_requirement_short"`
	Attributes           []string `json:"attributes"`
}

type Bar struct {
	T   string  `json:"t"`
	O   float64 `json:"o"`
	H   float64 `json:"h"`
	L   float64 `json:"l"`
	C   float64 `json:"c"`
	V   int     `json:"v"`
	N   int     `json:"n"`
	VW  float64 `json:"vw"`
}

type Response struct {
	Bars          map[string][]Bar `json:"bars"`
	NextPageToken string           `json:"next_page_token"`
}

func main () {

	dryRun, _ := os.LookupEnv("DRY_RUN")
	isDryRun, _ := strconv.ParseBool(dryRun)

	alpacaKey, _ := os.LookupEnv("ALPACA_KEY")
	alpacaSecret, _ := os.LookupEnv("ALPACA_SECRET")
	assetsUrl := "https://api.alpaca.markets/v2/assets"
	params := "?status=active&exchange=NYSE"
	pageToken := "" // Start with an empty token for the first page

	res := alpacaRequest("GET", alpacaKey, alpacaSecret, assetsUrl, params, nil)

	var items []Item 
	err := json.Unmarshal(res, &items)

	if err != nil {
		panic(err)
	}
	var itemList []string 
	for _, item := range items {
		if strings.Contains(item.Symbol, "/") || !item.Fractionable  || strings.Contains(strings.ToLower(item.Name), "short") || item.Symbol == "BAND" {
			// I don't want invalid symbols, non fractionable stocks, shorts, or Bandwidth
			continue
		}
		itemList = append(itemList, item.Symbol)
	}

	var first, second SymbolChange
	first.Change = math.MaxFloat64

	for {
		// Fetch data with current page token
		resp, err := getData(pageToken, alpacaKey, alpacaSecret, strings.Join(itemList, ","))
	
		if err != nil {
			panic(err)
		}

		// Iterate over all the stock symbols in the bars map
		for symbol, bars := range resp.Bars {
			for _, bar := range bars {
				symbolChange := SymbolChange{Symbol: symbol, Change: ((bar.C - bar.O) / bar.O ) * 100 }

				if symbolChange.Change < first.Change {
					// Update first and second smallest values
					second = first
					first = symbolChange
				} else if symbolChange.Change < second.Change {
					// Update only the second smallest value
					second = symbolChange
				}
			}
		}

		// Check if there is a next page
		if resp.NextPageToken == "" {
			break // No more pages
		}

		// Set the next page token for the next iteration
		pageToken = resp.NextPageToken
	}

	fmt.Printf("Biggest Losers: \n%s:%F, %s:%F\n\n", first.Symbol, first.Change, second.Symbol ,second.Change)

	buyingPower := findBuyingPower(alpacaKey, alpacaSecret, isDryRun)

	fmt.Printf("\nBuying Power: %f", buyingPower)

	buyOrder(first.Symbol, alpacaKey, alpacaSecret, isDryRun, buyingPower)
}

func alpacaRequest(method string, alpacaKey string, alpacaSecret string, url string, params string, body any) []byte {
	var req *http.Request
	if body == nil {
		req, _ = http.NewRequest(method, url + params, nil)
	} else {
		payload := body.(io.Reader)
		req, _ = http.NewRequest(method, url + params, payload)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("APCA-API-KEY-ID", alpacaKey)
	req.Header.Add("APCA-API-SECRET-KEY", alpacaSecret)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	return data
}

func buyOrder(symbol string, alpacaKey string, alpacaSecret string, dryRun bool, buyingPower float64) {
	var apiDomain string
	if !dryRun {
		apiDomain = "api"
	} else {
		apiDomain = "paper-api"
		alpacaKey, _ = os.LookupEnv("PAPER_ALPACA_KEY")
		alpacaSecret, _ = os.LookupEnv("PAPER_ALPACA_SECRET")
	}

	fmt.Printf("\n\nRouting requests to %s\n", apiDomain)

	url := fmt.Sprintf("https://%s.alpaca.markets/v2/orders", apiDomain)
	payload := strings.NewReader(fmt.Sprintf("{\"symbol\":\"%s\",\"notional\":\"%f\",\"side\":\"buy\",\"type\":\"market\",\"time_in_force\":\"day\"}", symbol, buyingPower))

	orderStatus := alpacaRequest("POST", alpacaKey, alpacaSecret, url, "", payload)

	fmt.Printf("\n\n\nOrder Info: \n%s", string(orderStatus))
}

func findBuyingPower(alpacaKey string, alpacaSecret string, isDryRun bool) float64 {
	var apiDomain string
	if !isDryRun {
		apiDomain = "api"
	} else {
		apiDomain = "paper-api"
		alpacaKey, _ = os.LookupEnv("PAPER_ALPACA_KEY")
		alpacaSecret, _ = os.LookupEnv("PAPER_ALPACA_SECRET")
	}

	accountURL := fmt.Sprintf("https://%s.alpaca.markets/v2/account", apiDomain)
	accountParams := ""

	accountRes := alpacaRequest("GET", alpacaKey, alpacaSecret, accountURL, accountParams, nil)

	var account map[string]interface{}
	e := json.Unmarshal(accountRes, &account)

	if e != nil {
		panic(e)
	}

	totalBuyingPower, _ := strconv.ParseFloat(account["cash"].(string), 32)

	return math.Floor( (totalBuyingPower) * 100 ) / 100 
}

func getData(pageToken string, alpacaKey string, alpacaSecret string, symbols string) (*Response, error) {
	url := "https://data.alpaca.markets/v2/stocks/bars"
	params := fmt.Sprintf("?symbols=%s&timeframe=1D&feed=iex", symbols)

	// If there's a page token, add it as a query parameter
	if pageToken != "" {
		params = fmt.Sprintf("%s&page_token=%s", params, pageToken)
	}

	// Prepare the request with the correct headers
	req, err := http.NewRequest("GET", url + params, nil)
	if err != nil {
		return nil, err
	}

	// Add the necessary headers
	req.Header.Add("accept", "application/json")
	req.Header.Add("APCA-API-KEY-ID", alpacaKey)
	req.Header.Add("APCA-API-SECRET-KEY", alpacaSecret)


	// Make the HTTP request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	// Parse the JSON response into the Response struct
	var resp Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
