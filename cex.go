package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		var gameData = strings.Split(scanner.Text(), ",")

		if gameData[1] == "buy" {
			printToScreenAndFile(f, "-------------------------------------------------------")
		} else {
			printToScreenAndFile(f, "=======================================================")
		}

		var detailResponse = getDetailResponse(gameData[0])
		printDetailData(gameData[1], detailResponse.Response.Data.BoxDetails, f)

		if gameData[1] == "buy" {
			var storesResponse = getStoresResponse(gameData[0])
			printStoreData(storesResponse.Response.Data.NearestStores, f)
		}
	}

	printToScreenAndFile(f, "=======================================================")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	f.Sync()
}

func getDetailResponse(gameID string) DetailResponse {
	resp, err := http.Get("https://wss2.cex.uk.webuy.io/v3/boxes/" + gameID + "/detail")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//printToScreenAndFile(fmt.Sprintf("HTTP Status: %d", resp.StatusCode))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//printToScreenAndFile(string(body))

	var response DetailResponse
	json.Unmarshal(body, &response)

	return response
}

func getStoresResponse(gameID string) StoresResponse {
	resp, err := http.Get("https://wss2.cex.uk.webuy.io/v3/boxes/" + gameID + "/neareststores?latitude=55.864237&longitude=-4.2518059")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//printToScreenAndFile(fmt.Sprintf("HTTP Status: %d", resp.StatusCode))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//printToScreenAndFile(string(body))

	var response StoresResponse
	json.Unmarshal(body, &response)

	return response
}

func printDetailData(action string, details []ItemDetailResponse, f *os.File) {
	for _, detail := range details {
		printToScreenAndFile(f, detail.BoxName)

		switch action {
		case "buy":
			printToScreenAndFile(f, fmt.Sprintf("    Sell Price: %s", formatFloat(detail.SellPrice, 2)))
			break
		case "sell":
			printToScreenAndFile(f, fmt.Sprintf("    Exchange Price: %s", formatFloat(detail.ExchangePrice, 2)))
			break
		}

	}
}

func printStoreData(nearestStores []NearestStoresResponse, f *os.File) {
	var found = false

	printToScreenAndFile(f, "")

	for _, store := range nearestStores {
		if strings.Contains(store.StoreName, "Glasgow") {
			printToScreenAndFile(f, fmt.Sprintf("    %s: %s", store.StoreName, getString(store.QuantityOnHand)))
			found = true
		}
	}

	if !found {
		printToScreenAndFile(f, "    Not found in any Glasgow store.")
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

func printToScreenAndFile(f *os.File, message string) {
	fmt.Println(message)

	f.WriteString(message + "\n")
}
