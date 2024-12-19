package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	dryRun, _ := os.LookupEnv("DRY_RUN")
	isDryRun, _ := strconv.ParseBool(dryRun)

	var apiDomain string
	if !isDryRun {
		apiDomain = "api"
	} else {
		apiDomain = "paper-api"
	}

	fmt.Printf("Routing requests to %s", apiDomain)

	alpacaKey, _ := os.LookupEnv("ALPACA_KEY")
	alpacaSecret, _ := os.LookupEnv("ALPACA_SECRET")
	url := fmt.Sprintf("https://%s.alpaca.markets/v2/positions", apiDomain)

	req, _ := http.NewRequest("DELETE", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("APCA-API-KEY-ID",alpacaKey)
	req.Header.Add("APCA-API-SECRET-KEY", alpacaSecret)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}
