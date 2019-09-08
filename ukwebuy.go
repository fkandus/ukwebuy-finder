package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// StoresResponse represents the topmost level
type StoresResponse struct {
	Response CexStoresResponse
}

// CexStoresResponse is the Response to nearest stores request
type CexStoresResponse struct {
	Data StoresDataResponse
}

// StoresDataResponse is the list of nearest stores
type StoresDataResponse struct {
	NearestStores []NearestStoresResponse
}

// NearestStoresResponse represents the important data of a store
type NearestStoresResponse struct {
	StoreName      string
	QuantityOnHand interface{}
}

// DetailResponse represents the topmost level
type DetailResponse struct {
	Response CexDetailResponse
}

// CexDetailResponse is the Response to nearest stores request
type CexDetailResponse struct {
	Data DetailDataResponse
}

// DetailDataResponse is the list of nearest stores
type DetailDataResponse struct {
	BoxDetails []ItemDetailResponse
}

// ItemDetailResponse represents the important data of a store
type ItemDetailResponse struct {
	BoxName       string
	SellPrice     float64
	ExchangePrice float64
}

func getDetailResponse(gameID string, config Configuration) DetailResponse {
	resp, err := http.Get(strings.Replace(config.Urls.Detail, "{gameID}", gameID, 1))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var response DetailResponse
	json.Unmarshal(body, &response)

	return response
}

func getStoresResponse(gameID string, location Location, config Configuration) StoresResponse {
	r := strings.NewReplacer("{gameID}", gameID, "{latitude}", location.Lat, "{longitude}", location.Lon)

	resp, err := http.Get(r.Replace(config.Urls.Store))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var response StoresResponse
	json.Unmarshal(body, &response)

	return response
}
