package main

import (
	"flag"
	"fmt"
	"sort"
	"time"
)

type videoMeta struct {
	Name          string
	VideoDuration time.Duration
	Language      string
	ClickURL      string
	PublishYear   int
	PublishMonth  time.Month
	ThumbnailURL  string
}

func main() {
	lang := flag.String("lang", "", "filter by language [en-US, hi-IN]")
	durationMin := flag.Duration("minDuration", 0, "filter by minimum duration [such as 30s, 20m, 1h]")
	durationMax := flag.Duration("maxDuration", 0, "filter by maximum duration [such as 30s, 20m, 1h]")
	publishYear := flag.Int("publishYear", 0, "filter by publish year [such as 2022, 2023, 2024]")
	flag.Parse()

	if err := ensureCache(); err != nil {
		panic(err)
	}

	videos, err := readFromCache()
	if err != nil {
		panic(err)
	}

	filteredVideos, err := filterContent(videos, *lang, *durationMin, *durationMax, *publishYear)
	if err != nil {
		panic(err)
	}

	fmt.Println("total filtered videos latest to oldest:", len(filteredVideos))
	for _, video := range filteredVideos {
		fmt.Printf("[%v] in [%v-%v] of [%v]: %v\n", video.Name, video.PublishMonth,
			video.PublishYear, video.VideoDuration, video.ClickURL)
	}
}

func filterContent(videos []videoMeta, lang string, durationMin,
	durationMax time.Duration, publishYear int) ([]videoMeta, error) {

	var filteredVideos []videoMeta
	for _, video := range videos {
		if lang != "" && video.Language != lang {
			continue
		}
		if durationMin != 0 && video.VideoDuration < durationMin {
			continue
		}
		if durationMax != 0 && video.VideoDuration > durationMax {
			continue
		}
		if publishYear != 0 && video.PublishYear != publishYear {
			continue
		}

		filteredVideos = append(filteredVideos, video)
	}

	return sortVideosByPublishYear(filteredVideos), nil
}

func sortVideosByPublishYear(videos []videoMeta) []videoMeta {
	sort.Slice(videos, func(i, j int) bool {
		if videos[i].PublishYear != videos[j].PublishYear {
			return videos[i].PublishYear > videos[j].PublishYear
		}
		return videos[i].PublishMonth > videos[j].PublishMonth
	})
	return videos
}
