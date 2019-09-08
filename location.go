package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Location represents the important data of a store
type Location struct {
	Lat string
	Lon string
}

func getLocations(config Configuration) []Location {
	client := &http.Client{}

	req, err := http.NewRequest("GET", strings.Replace(config.Urls.Location, "{city}", config.Locations.City, 1), nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "UK-Webuy-Finder/1.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var locations []Location
	json.Unmarshal(body, &locations)

	return locations
}
