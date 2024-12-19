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
	// "github.com/alpacahq/alpaca-trade-api-go/alpaca"
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

func main () {

	dryRun, _ := os.LookupEnv("DRY_RUN")
	isDryRun, _ := strconv.ParseBool(dryRun)

	var alpacaKey, alpacaSecret string
	if !isDryRun {
		alpacaKey, _ = os.LookupEnv("ALPACA_KEY")
		alpacaSecret, _ = os.LookupEnv("ALPACA_SECRET")
	} else {
		alpacaKey, _ = os.LookupEnv("PAPER_ALPACA_KEY")
		alpacaSecret, _ = os.LookupEnv("PAPER_ALPACA_SECRET")
	}

	assetsUrl := "https://api.alpaca.markets/v2/assets"
	params := "?status=active&exchange=NASDAQ"

	res := alpacaRequest("GET", alpacaKey, alpacaSecret, assetsUrl, params, nil)

	var items []Item 
	err := json.Unmarshal(res, &items)

	if err != nil {
		panic(err)
	}

	first, second := findBiggestLosers(items, alpacaKey, alpacaSecret)

	fmt.Printf("%s:%F, %s:%F", first.Symbol, first.Change, second.Symbol ,second.Change)

	buyOrder(first.Symbol, alpacaKey, alpacaSecret, isDryRun)
	buyOrder(second.Symbol, alpacaKey, alpacaSecret, isDryRun)
	 
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


func findBiggestLosers(assets []Item, alpacaKey string, alpacaSecret string) (SymbolChange, SymbolChange){
	var first, second SymbolChange
	first.Change = math.MaxFloat64

	for _, asset := range assets {
		if asset.Symbol == "BAND" {
			continue
		}

		percentChange := getPercentChange(asset.Symbol, alpacaKey, alpacaSecret)

		if percentChange == 0 {
			continue
		}

		data := SymbolChange{
			Symbol: asset.Symbol,
			Change: percentChange,
		}

		if data.Change < first.Change {
			// Update first and second smallest values
			second = first
			first = data
		} else if data.Change < second.Change {
			// Update only the second smallest value
			second = data
		}
	}
	return first, second
}

func getPercentChange(symbol string, alpacaKey string, alpacaSecret string) float64 {
	url := fmt.Sprintf("https://data.alpaca.markets/v2/stocks/%s/bars", symbol)
	params := "?timeframe=1D&feed=iex"

	res := alpacaRequest("GET", alpacaKey, alpacaSecret, url, params, nil)
	
	var data map[string]interface{}
	err := json.Unmarshal(res, &data)

	if err != nil {
	  panic(err)
	}

	if data["bars"] == nil {
		return 0
	}

	bar := data["bars"].([]interface{})[0].(map[string]interface{})
	open := bar["o"].(float64)
	close := bar["c"].(float64)

	if close < 5 {
		return 0
	}
	
	return ((close - open) / open ) * 100
}

func buyOrder(symbol string, alpacaKey string, alpacaSecret string, dryRun bool) {
	var apiDomain string
	if !dryRun {
		apiDomain = "api"
	} else {
		apiDomain = "paper-api"
	}

	fmt.Printf("Routing requests to %s\n", apiDomain)

	accountURL := fmt.Sprintf("https://%s.alpaca.markets/v2/account", apiDomain)
	params := ""

	res := alpacaRequest("GET", alpacaKey, alpacaSecret, accountURL, params, nil)

	var account map[string]interface{}
	err := json.Unmarshal(res, &account)

	if err != nil {
		panic(err)
	}

	fmt.Print(account)

	buying_power := account["buying_power"].(float32)

	url := fmt.Sprintf("https://%s.alpaca.markets/v2/orders", apiDomain)
	payload := strings.NewReader(fmt.Sprintf("{\"symbol\":\"%s\",\"notional\":\"%f\",\"side\":\"buy\",\"type\":\"market\",\"time_in_force\":\"day\"}", symbol, buying_power/2))

	alpacaRequest("POST", alpacaKey, alpacaSecret, url, params, payload)
}
