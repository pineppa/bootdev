package pokeAPI

import (
	"io"
	"net/http"
	"pokeCLI/internal/pokecache"
	"time"
)

var cache = pokecache.NewCache(5 * time.Minute)

func fetchFromURL(url string) ([]byte, error) {
	content, ok := cache.Get(url)
	if ok {
		return content, nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	cache.Add(url, body)
	return body, err
}
