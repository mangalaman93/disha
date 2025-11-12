package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"time"
)

type videoMeta struct {
	VideoID       string
	Name          string
	VideoDuration time.Duration
	Language      string
	ClickURL      string
	PublishYear   int
	PublishMonth  time.Month
	ThumbnailURL  string
}

type filterParam struct {
	lang        string
	durationMin time.Duration
	durationMax time.Duration
	publishYear int
	source      string
}

func main() {
	lang := flag.String("lang", "", "filter by language [en-US, hi-IN]")
	durationMin := flag.Duration("minDuration", 0, "filter by minimum duration [such as 30s, 20m, 1h]")
	durationMax := flag.Duration("maxDuration", 0, "filter by maximum duration [such as 30s, 20m, 1h]")
	publishYear := flag.Int("publishYear", 0, "filter by publish year [such as 2022, 2023, 2024]")
	source := flag.String("source", "", "filter by source [youtube, tt]")
	flag.Parse()

	if err := cache.setup(); err != nil {
		panic(err)
	}

	if *source != "" && *source == "tt" {
		*source = "timelesstoday"
	}
	params := filterParam{
		lang:        *lang,
		durationMin: *durationMin,
		durationMax: *durationMax,
		publishYear: *publishYear,
		source:      *source,
	}

	filteredVideos, err := filterContent(cache.Videos, params)
	if err != nil {
		panic(err)
	}
	fmt.Println("total filtered videos latest to oldest:", len(filteredVideos))

	for _, video := range filteredVideos {
		fmt.Printf("[%v] in [%v-%v] of [%v]: %v\n", video.Name, video.PublishMonth,
			video.PublishYear, video.VideoDuration, video.ClickURL)
	}
}

func filterContent(videos map[string]videoMeta, param filterParam) ([]videoMeta, error) {
	var filteredVideos []videoMeta
	for _, video := range videos {
		if param.lang != "" && video.Language != param.lang {
			continue
		}
		if param.durationMin != 0 && video.VideoDuration < param.durationMin {
			continue
		}
		if param.durationMax != 0 && video.VideoDuration > param.durationMax {
			continue
		}
		if param.publishYear != 0 && video.PublishYear != param.publishYear {
			continue
		}
		if param.source != "" && !strings.Contains(video.ClickURL, param.source) {
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
