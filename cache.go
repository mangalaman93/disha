package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	cacheFile = "cache.json"
)

var (
	cache videoCache
)

type videoCache struct {
	Videos      map[string]videoMeta `json:"videos"`
	LastUpdated time.Time            `json:"lastUpdated"`
}

func (c *videoCache) set(video videoMeta) {
	c.Videos[video.VideoID] = video
}

func (c *videoCache) get(id string) (videoMeta, bool) {
	video, ok := c.Videos[id]
	return video, ok
}

func (c *videoCache) load() error {
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return fmt.Errorf("error reading cache file [%v]: %w", cacheFile, err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("error unmarshalling cache file [%v]: %w", cacheFile, err)
	}

	if c.LastUpdated.Add(time.Hour * 24).Before(time.Now()) {
		log.Println("cache is old, downloading")
		return c.download()
	}

	return nil
}

func (c *videoCache) save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling cache: %w", err)
	}

	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		return fmt.Errorf("error writing cache file: %w", err)
	}

	return nil
}

func (c *videoCache) download() error {
	videosFromYouTube, err := getYouTubeContent()
	if err != nil {
		return fmt.Errorf("error getting video list from YouTube: %w", err)
	}
	log.Println("total videos retrieved from YouTube:", len(videosFromYouTube))

	videosFromTT, err := getTTContent()
	if err != nil {
		return fmt.Errorf("error getting video list from TT: %w", err)
	}
	log.Println("total videos retrieved from TT:", len(videosFromTT))

	if c.Videos == nil {
		c.Videos = make(map[string]videoMeta)
	} else {
		clear(c.Videos)
	}

	for _, video := range append(videosFromTT, videosFromYouTube...) {
		c.set(video)
	}
	c.LastUpdated = time.Now()

	customizeCache(c)
	return c.save()
}

func (c *videoCache) setup(updateCache bool) error {
	if updateCache {
		log.Println("update for cache requested!")
		return c.download()
	}

	if _, err := os.Stat(cacheFile); err == nil {
		log.Println("cache file already exists, no need to download")
		return c.load()
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error checking for cache file: %w", err)
	} else {
		log.Println("cache file does not exist, downloading")
		return c.download()
	}
}
