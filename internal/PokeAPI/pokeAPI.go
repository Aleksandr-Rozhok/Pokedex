package pokeAPI

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func PokeAPI(url string) []byte {
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

	return body
}
