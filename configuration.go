package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Configuration is the top most level of the config
type Configuration struct {
	Urls      URLData
	Locations LocationData
	Stores    StoreData
}

// URLData contains the needed URLs
type URLData struct {
	Detail string
	Store  string
	Location string
}

// LocationData contains the name of the place to get coordinates for
type LocationData struct {
	City    string
	Country string
}

// StoreData contains the names of stores to match
type StoreData struct {
	MatchName string
}

func getConfig() Configuration {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	return configuration
}
