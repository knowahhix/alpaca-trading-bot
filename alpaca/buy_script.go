package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"encoding/json"
	// "github.com/alpacahq/alpaca-trade-api-go/alpaca"
)
type SymbolChange struct {
	Symbol string
	Change float32
}

func main () {
	alpacaKey, _ := os.LookupEnv("ALPACA_KEY")
	alpacaSecret, _ := os.LookupEnv("ALPACA_SECRET")
	assetsUrl := "https://api.alpaca.markets/v2/assets"

	params := "?status=active&asset_class=us_equity&exchange=NASDAQ"

	res := alpacaRequest("GET", alpacaKey, alpacaSecret, assetsUrl, params)

	fmt.Println(res)
}

func alpacaRequest(method string, alpacaKey string, alpacaSecret string, url string, params string) []map[string] {
	req, _ := http.NewRequest(method, url + params, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("APCA-API-KEY-ID", alpacaKey)
	req.Header.Add("APCA-API-SECRET-KEY", alpacaSecret)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var data []map[string]
	err := json.Unmarshal(body, &data)

	if err != nil {
	  panic(err)
	}

	fmt.Print(data)

	return data
}

// func main() {
// 	alpacaKey, _ := os.LookupEnv("ALPACA_KEY")
// 	alpaceSecret, _ := os.LookupEnv("ALPACA_SECRET")
// 	fmt.Printf("%s", alpacaKey)
// 	baseUrl := "https://paper-api.alpaca.markets"

// 	alpaca.SetBaseUrl(baseUrl)
// 	creds := common.APIKey{
// 		ID: alpacaKey,
// 		Secret: alpaceSecret,
// 	}

// 	client := alpaca.NewClient(&creds)

// 	status := "active"
// 	assets, err := client.ListAssets(&status)
// 	if err != nil {
// 		panic(err)
// 	}
	
// 	nasdaq_assets := []alpaca.Asset{}
// 	for _, asset := range assets {
// 		if asset.Exchange == "NASDAQ" {
// 			nasdaq_assets = append(nasdaq_assets, asset)
// 		}
// 	}

// 	first, second := findBiggestLosers(client, nasdaq_assets)

// 	buyOrder(first.Symbol, client)
// 	buyOrder(second.Symbol, client)
// }

// func findBiggestLosers(client *alpaca.Client, assets []alpaca.Asset) (SymbolChange, SymbolChange){
// 	losers := []SymbolChange{}
// 	for _, asset := range assets {
// 		open := getStartPrice(asset.Symbol, client)
// 		current := getCurrentPrice(asset.Symbol, client)
// 		percentChange := (current / open) * 100
// 		data := SymbolChange{
// 			Symbol: asset.Symbol,
// 			Change: percentChange,
// 		}
// 		losers = append(losers, data)
// 	}
	
// 	var first, second SymbolChange
// 	first.Change = math.MaxFloat32

// 	for _, sc := range losers {
// 		if sc.Change < first.Change {
// 			// Update first and second smallest values
// 			second = first
// 			first = sc
// 		} else if sc.Change < second.Change {
// 			// Update only the second smallest value
// 			second = sc
// 		}
// 	}

// 	return first, second

// }

// func getCurrentPrice(symbol string, client *alpaca.Client) float32 {
// 	resp, err := client.GetLastQuote(symbol)

// 	if err != nil {
// 		panic(err)
// 	}

// 	return resp.Last.AskPrice
// }

// func getStartPrice(symbol string, client *alpaca.Client) float32 {
// 	limit := 1
// 	params := alpaca.ListBarParams{
// 		Timeframe: "1D",
// 		Limit: &limit,
// 	}

// 	bars, err := client.GetSymbolBars(symbol, params)

// 	if err != nil {
// 		panic(err)
// 	}

// 	openPrice := bars[0].Open

// 	return openPrice
// }

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
