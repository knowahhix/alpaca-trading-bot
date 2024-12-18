package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
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
	alpacaKey, _ := os.LookupEnv("ALPACA_KEY")
	alpacaSecret, _ := os.LookupEnv("ALPACA_SECRET")
	assetsUrl := "https://api.alpaca.markets/v2/assets"
	params := "?status=active&exchange=NASDAQ"

	res := alpacaRequest("GET", alpacaKey, alpacaSecret, assetsUrl, params)

	fmt.Print(res)
	fmt.Print("\n")

	var data []Item
	err := json.Unmarshal(res, &data)

	if err != nil {
	  panic(err)
	}

	first, second := findBiggestLosers(data, alpacaKey, alpacaSecret)

	fmt.Printf("%f , %f", first.Change, second.Change)
	
}

func alpacaRequest(method string, alpacaKey string, alpacaSecret string, url string, params string) []byte {
	req, _ := http.NewRequest(method, url + params, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("APCA-API-KEY-ID", alpacaKey)
	req.Header.Add("APCA-API-SECRET-KEY", alpacaSecret)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Print(body)
	fmt.Print("\n")

	return body
}

func findBiggestLosers(assets []Item, alpacaKey string, alpacaSecret string) (SymbolChange, SymbolChange){
	losers := []SymbolChange{}
	for _, asset := range assets {
		current := getCurrentPrice(asset.Symbol, alpacaKey, alpacaSecret)
		open := getStartPrice(asset.Symbol, alpacaKey, alpacaSecret)
		percentChange := (current / open) * 100
		data := SymbolChange{
			Symbol: asset.Symbol,
			Change: percentChange,
		}
		losers = append(losers, data)
	}
	
	var first, second SymbolChange
	first.Change = math.MaxFloat32

	for _, sc := range losers {
		if sc.Change < first.Change {
			// Update first and second smallest values
			second = first
			first = sc
		} else if sc.Change < second.Change {
			// Update only the second smallest value
			second = sc
		}
	}

	return first, second

}

func getCurrentPrice(symbol string, alpacaKey string, alpacaSecret string) float64 {
	url := fmt.Sprintf("https://data.alpaca.markets/v2/stocks/%s/quotes/latest", symbol)
	params := "?feed=iex"
	res := alpacaRequest("GET", alpacaKey, alpacaSecret, url, params)

	var data map[string]interface{}
	err := json.Unmarshal(res, &data)

	if err != nil {
	  panic(err)
	}

	fmt.Print(data)

	quote := data["quote"].(map[string]interface{})

	return quote["ap"].(float64)
}

func getStartPrice(symbol string, alpacaKey string, alpacaSecret string) float64 {
	url := fmt.Sprintf("https://data.alpaca.markets/v2/stocks/%s/bars", symbol)
	params := "?timeframe=1D&feed=iex"

	res := alpacaRequest("GET", alpacaKey, alpacaSecret, url, params)
	
	var data map[string]interface{}
	err := json.Unmarshal(res, &data)

	if err != nil {
	  panic(err)
	}

	bars := data["bars"].([]interface{})

	fmt.Print(bars)

	bar := bars[0].(map[string]interface{})

	return bar["o"].(float64)
}

// func buyOrder(symbol string, client *alpaca.Client) {
// 	accountID, err := client.GetAccount()
// 	if err != nil {
// 		panic(err)
// 	}

// 	amt := accountID.BuyingPower.Div(decimal.NewFromInt(2))

// 	order := alpaca.PlaceOrderRequest{
// 		AccountID: accountID.ID,
// 		AssetKey: &symbol,
// 		Notional: amt,
// 		Side: "buy",
// 		Type: "market",
// 		TimeInForce: "day",
// 	}
// 	client.PlaceOrder(order)
// }
