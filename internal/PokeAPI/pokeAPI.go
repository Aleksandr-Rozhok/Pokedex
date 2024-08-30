package PokeAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type LocationResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationData struct {
	Count    int              `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"` // interface{} because it can be null
	Results  []LocationResult `json:"results"`
}

func PokeAPI(url string) LocationData {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	err = res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	if err != nil {
		log.Fatal(err)
	}

	var result LocationData

	err = json.Unmarshal(body, &result)
	if err != nil {
		return LocationData{}
	}

	return result
}
