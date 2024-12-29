package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Location represents the important data of a store
type Location struct {
	Lat string
	Lon string
}

func getLocations(config Configuration) []Location {
	if _, err := os.Stat("location_cache.txt"); err == nil {
		// cache file exists
		return readCache()
	}

	return executeGet(config)
}

func readCache() []Location {
	var locations []Location

	file, err := os.Open("location_cache.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var parts = strings.Split(scanner.Text(), ",")

		locations = append(locations, Location{parts[0], parts[1]})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return locations
}

func executeGet(config Configuration) []Location {
	client := &http.Client{}

	var locations []Location

	for _, city := range config.Locations.City {
		req, err := http.NewRequest("GET", strings.Replace(config.Urls.Location, "{city}", city, 1), nil)
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Set("User-Agent", "UK-Webuy-Finder/1.0")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var locs []Location
		json.Unmarshal(body, &locs)

		locations = append(locations, locs...)
	}

	writeCache(locations)

	return locations
}

func writeCache(locations []Location) {
	file, err := os.Create("location_cache.txt")

	if err != nil {
		return
	}
	defer file.Close()

	for _, l := range locations {
		fmt.Fprintln(file, l.Lat+","+l.Lon)
	}
}
