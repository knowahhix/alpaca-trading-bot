package main

import (
	"os"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

func main() {
	alpacaKey, _ := os.LookupEnv("ALPACA_KEY")
	alpaceSecret, _ := os.LookupEnv("ALPACA_SECRET")

	client := alpaca.NewClient(alpaca.ClientOpts{
		APIKey: alpacaKey,
		APISecret: alpaceSecret,
		BaseURL: "https://paper-api.alpaca.markets",
	})
	// TODO: Find 2 biggest losesr of S & P  500 for the day


	// TODO: Buy the bigest losers
	
}
