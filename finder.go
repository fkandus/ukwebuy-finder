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

	var availableTotalBuy float64 = 0
	var totalBuy float64 = 0
	var totalSell float64 = 0
	var storeCount = make(map[string]int)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		var gameData = strings.Split(scanner.Text(), ",")

		printSeparatorLine(f, gameData[1])

		var detailResponse = getDetailResponse(gameData[0], config)
		printDetailData(gameData, detailResponse.Response.Data.BoxDetails, f)

		if gameData[1] == "buy" {
			var storesResponse = getStoresResponse(gameData[0], locations[0], config)
			if processStores(storesResponse, config, storeCount, f) {
				availableTotalBuy += detailResponse.Response.Data.BoxDetails[0].SellPrice
			}
			totalBuy += detailResponse.Response.Data.BoxDetails[0].SellPrice
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

	printToScreenAndFile(f, fmt.Sprintf("Available items total Buy Value: £%s", formatFloat(availableTotalBuy, 2)))
	printToScreenAndFile(f, fmt.Sprintf("All items total Buy Value: £%s", formatFloat(totalBuy, 2)))
	printToScreenAndFile(f, fmt.Sprintf("Total Sell Value: £%s", formatFloat(totalSell, 2)))
	printToScreenAndFile(f, fmt.Sprintf("Buy-Sell difference (available): £%s", formatFloat(totalSell-availableTotalBuy, 2)))
	printToScreenAndFile(f, fmt.Sprintf("Buy-Sell difference (total): £%s", formatFloat(totalSell-totalBuy, 2)))

	printSeparatorLine(f, "sell")

	printStoreCount(storeCount, f)
}

func printDetailData(gameData []string, details []ItemDetailResponse, f *os.File) {
	var id = gameData[0]
	var action = gameData[1]

	for _, detail := range details {
		printToScreenAndFile(f, fmt.Sprintf("%s (%s) - (%s)", detail.BoxName, detail.CategoryFriendlyName, id))

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

func processStores(storesResponse StoresResponse, config Configuration, storeCount map[string]int, f *os.File) bool {
	var stores = filterStores(storesResponse.Response.Data.NearestStores, config)

	if len(stores) > 0 {
		handleStores(stores, storeCount, f)

		return true
	} else {
		printToScreenAndFile(f, "    Not found in any store.")

		return false
	}
}

func filterStores(nearestStores []NearestStoresResponse, config Configuration) []NearestStoresResponse {
	var filteredStores []NearestStoresResponse

	for _, store := range nearestStores {
		if strings.Contains(store.StoreName, config.Stores.MatchName) {
			filteredStores = append(filteredStores, store)
		}
	}

	return filteredStores
}

func handleStores(stores []NearestStoresResponse, storeCount map[string]int, f *os.File) {
	printToScreenAndFile(f, "")

	for _, store := range stores {
		printToScreenAndFile(f, fmt.Sprintf("    %s: %s", store.StoreName, getString(store.QuantityOnHand)))
		_, ok := storeCount[store.StoreName]

		if ok {
			storeCount[store.StoreName]++
		} else {
			storeCount[store.StoreName] = 1
		}
	}
}

func printStoreCount(storeCount map[string]int, f *os.File) {
	printToScreenAndFile(f, "Store counter:")

	for key, value := range storeCount {
		printToScreenAndFile(f, fmt.Sprint(key, ": ", value))
	}
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
		printToScreenAndFile(f, "------------------------------------------------------------------------------------------")
	} else if action == "sell" {
		printToScreenAndFile(f, "==========================================================================================")
	}
}

func printToScreenAndFile(f *os.File, message string) {
	fmt.Println(message)

	f.WriteString(message + "\n")
}
