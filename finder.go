package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
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

	var inputFile = os.Args[1]

	input, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()

	inputBasename := strings.TrimSuffix(inputFile, path.Ext(inputFile))

	dt := time.Now()

	f, err := os.Create("result-" + inputBasename + "-" + dt.Format("20060102-150405") + ".txt")
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
			var storesResponses = getStoresResponse(gameData[0], locations, config)
			if processStores(storesResponses, config, storeCount, f) {
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

func processStores(storeResponses []StoresResponse, config Configuration, storeCount map[string]int, f *os.File) bool {
	var stores = filterStores(storeResponses, config)

	if len(stores) > 0 {
		handleStores(stores, storeCount, f)

		return true
	}

	printToScreenAndFile(f, "    Not found in any store.")

	return false
}

func filterStores(storeResponses []StoresResponse, config Configuration) []NearestStoresResponse {
	var filteredStores []NearestStoresResponse

	for _, storeResponse := range storeResponses {
		for _, store := range storeResponse.Response.Data.NearestStores {
			if matchStore(store.StoreName, config.Stores.MatchName) {
				filteredStores = append(filteredStores, store)
			}
		}
	}

	return filteredStores
}

func matchStore(storeName string, matchNames []string) bool {
	for _, matchName := range matchNames {
		if strings.Contains(storeName, matchName) {
			return true
		}
	}

	return false
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

	keys := sortedKeys(storeCount)

	for _, k := range keys {
		printToScreenAndFile(f, fmt.Sprint(k, ": ", storeCount[k]))
	}
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
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
