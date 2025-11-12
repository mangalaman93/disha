package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	hindiLang   = "hi-IN"
	englishLang = "en-US"

	localContentFile = "content.json"
	ttListContentURL = "https://api3.timelesstoday.io/v2/cms/products/en-US/language/%v/10000/0"
	ttVideoURL       = "https://www.timelesstoday.tv/en/home/product/%v"
)

var (
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	allLangs = []string{hindiLang, englishLang}
)

type videoContent struct {
	Name         string
	DurationSec  int
	Language     string
	ClickURL     string
	PublishYear  int
	PublishMonth time.Month
	ThumbnailURL string
}

func main() {
	if err := ensureMetadata(); err != nil {
		panic(err)
	}
}

func ensureMetadata() error {
	// TODO: if file is 24 hours old, delete it and hit APIs again
	if _, err := os.Stat(localContentFile); err == nil {
		fmt.Println("local content file already exists, no need to hit APIs")
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error checking for local content file: %w", err)
	}

	videos, err := getTTContent()
	if err != nil {
		return fmt.Errorf("error getting video list from TT: %w", err)
	}
	fmt.Println("total videos retrieved from TT: ", len(videos))

	data, err := json.Marshal(videos)
	if err != nil {
		return fmt.Errorf("error marshalling video list from TT: %w", err)
	}

	if err := os.WriteFile(localContentFile, data, 0644); err != nil {
		return fmt.Errorf("error writing local content file: %w", err)
	}

	return nil
}

func getTTContent() ([]videoContent, error) {
	var videoList []videoContent
	for _, lang := range allLangs {
		videos, err := getContentForLang(lang)
		if err != nil {
			return nil, err
		}
		videoList = append(videoList, videos...)
	}

	return videoList, nil
}

func getContentForLang(lang string) ([]videoContent, error) {
	fmt.Println("getting video list for lang: ", lang)

	req, err := http.NewRequest("GET", fmt.Sprintf(ttListContentURL, lang), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for lang [%v]: %w", lang, err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:144.0) Gecko/20100101 Firefox/144.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting video list for lang [%v]: %w", lang, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad http status [%v] while geting video list for lang [%v]", resp.Status, lang)
	}
	defer resp.Body.Close()

	var respstruct struct {
		Data []struct {
			Name           string `json:"tt_name"`
			DurationSec    int    `json:"tt_duration"`
			SourceLanguage string `json:"tt_source_language"`
			MediaUUID      string `json:"tt_media_uuid"`
			PublishDate    string `json:"tt_publishing_date"`
			ThumbnailURL   string `json:"tt_image_url"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&respstruct); err != nil {
		return nil, fmt.Errorf("error decoding video list for lang [%v]: %w", lang, err)
	}

	var videoList []videoContent
	for _, video := range respstruct.Data {
		publishTs, err := time.Parse("2006-01-02T15:04:05.999", video.PublishDate)
		if err != nil {
			return nil, fmt.Errorf("error parsing publish date for video [%+v]: %w", video, err)
		}

		videoList = append(videoList, videoContent{
			Name:         video.Name,
			DurationSec:  video.DurationSec,
			Language:     video.SourceLanguage,
			ClickURL:     fmt.Sprintf(ttVideoURL, video.MediaUUID),
			PublishYear:  publishTs.Year(),
			PublishMonth: publishTs.Month(),
			ThumbnailURL: video.ThumbnailURL,
		})
	}

	return videoList, nil
}
