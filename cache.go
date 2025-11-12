package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	cacheFile = "cache.json"
)

func ensureCache() error {
	// TODO: if file is 24 hours old, delete it and hit APIs again
	if _, err := os.Stat(cacheFile); err == nil {
		fmt.Println("cache file already exists, no need to hit APIs")
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error checking for cache file: %w", err)
	}

	videos, err := getTTContent()
	if err != nil {
		return fmt.Errorf("error getting video list from TT: %w", err)
	}
	fmt.Println("total videos retrieved from TT:", len(videos))

	data, err := json.MarshalIndent(videos, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling video list from TT: %w", err)
	}

	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		return fmt.Errorf("error writing cache file: %w", err)
	}

	return nil
}

func readFromCache() ([]videoMeta, error) {
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, fmt.Errorf("error reading cache file: %w", err)
	}

	var videoList []videoMeta
	if err := json.Unmarshal(data, &videoList); err != nil {
		return nil, fmt.Errorf("error unmarshalling cache file: %w", err)
	}

	return videoList, nil
}
