package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Missing file argument.")
		return
	}

	var config = getConfig()

	var locations = getLocations(config)

	if len(locations) == 0 {
		fmt.Println("No locations found for the current config.")
		return
	}

	input, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	dt := time.Now()

	f, err := os.Create("trade-games-" + dt.Format("20060102-150405") + ".txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var totalBuy float64 = 0
	var totalSell float64 = 0

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		var gameData = strings.Split(scanner.Text(), ",")

		printSeparatorLine(f, gameData[1])

		var detailResponse = getDetailResponse(gameData[0], config)
		printDetailData(gameData[1], detailResponse.Response.Data.BoxDetails, f)

		if gameData[1] == "buy" {
			var storesResponse = getStoresResponse(gameData[0], locations[0], config)
			if handleStoreData(storesResponse.Response.Data.NearestStores, &totalBuy, config, f) {
				totalBuy += detailResponse.Response.Data.BoxDetails[0].SellPrice
			}
		}

		if gameData[1] == "sell" {
			totalSell += detailResponse.Response.Data.BoxDetails[0].ExchangePrice
		}
	}

	printSeparatorLine(f, "sell")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	f.Sync()

	printToScreenAndFile(f, fmt.Sprintf("Total Buy Value: £%s", formatFloat(totalBuy, 2)))
	printToScreenAndFile(f, fmt.Sprintf("Total Sell Value: £%s", formatFloat(totalSell, 2)))
	printToScreenAndFile(f, fmt.Sprintf("Buy-Sell difference: £%s", formatFloat(totalSell-totalBuy, 2)))
}

func printDetailData(action string, details []ItemDetailResponse, f *os.File) {
	for _, detail := range details {
		printToScreenAndFile(f, fmt.Sprintf("%s (%s)", detail.BoxName, detail.CategoryFriendlyName))

		switch action {
		case "buy":
			printToScreenAndFile(f, fmt.Sprintf("    Buy for: £%s", formatFloat(detail.SellPrice, 2)))
			break
		case "sell":
			printToScreenAndFile(f, fmt.Sprintf("    Sell at: £%s", formatFloat(detail.ExchangePrice, 2)))
			break
		}

	}
}

func handleStoreData(nearestStores []NearestStoresResponse, totalBuy *float64, config Configuration, f *os.File) bool {
	var found = false

	printToScreenAndFile(f, "")

	for _, store := range nearestStores {
		if strings.Contains(store.StoreName, config.Stores.MatchName) {
			printToScreenAndFile(f, fmt.Sprintf("    %s: %s", store.StoreName, getString(store.QuantityOnHand)))
			found = true
		}
	}

	if !found {
		printToScreenAndFile(f, "    Not found in any store.")
	}

	return found
}

func getString(v interface{}) string {
	switch v := v.(type) {
	case float64:
		return formatFloat(v, 0)
	case string:
		return v
	}

	return ""
}

func formatFloat(f float64, d int) string {
	return strconv.FormatFloat(f, 'f', d, 64)
}

func printSeparatorLine(f *os.File, action string) {
	if action == "buy" {
		printToScreenAndFile(f, "--------------------------------------------------------------------------------")
	} else if action == "sell" {
		printToScreenAndFile(f, "================================================================================")
	}
}

func printToScreenAndFile(f *os.File, message string) {
	fmt.Println(message)

	f.WriteString(message + "\n")
}
