package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	hindiLang   = "hi-IN"
	englishLang = "en-US"

	ttListContentURL = "https://api3.timelesstoday.io/v2/cms/products/en-US/language/%v/10000/0"
	ttVideoURL       = "https://www.timelesstoday.tv/%v/home/product/%v"
)

var (
	httpClient = &http.Client{Timeout: 10 * time.Second}
)

func getTTContent() ([]videoMeta, error) {
	var videoList []videoMeta

	for _, lang := range []string{hindiLang, englishLang} {
		videos, err := getContentForLang(lang)
		if err != nil {
			return nil, err
		}
		videoList = append(videoList, videos...)
	}

	return videoList, nil
}

func getContentForLang(lang string) ([]videoMeta, error) {
	log.Println("getting video list for lang:", lang)

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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("error while closing response body for lang %s: %v", lang, err)
		}
	}()

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

	var videoList []videoMeta
	for _, video := range respstruct.Data {
		publishTs, err := time.Parse("2006-01-02T15:04:05.999", video.PublishDate)
		if err != nil {
			return nil, fmt.Errorf("error parsing publish date for video [%+v]: %w", video, err)
		}

		videoList = append(videoList, videoMeta{
			VideoID:       video.MediaUUID,
			Name:          video.Name,
			VideoDuration: time.Duration(video.DurationSec) * time.Second,
			Language:      video.SourceLanguage,
			ClickURL:      getClickURL(video.MediaUUID, video.SourceLanguage),
			PublishYear:   publishTs.Year(),
			PublishMonth:  publishTs.Month(),
			ThumbnailURL:  video.ThumbnailURL,
		})
	}

	return videoList, nil
}

func getClickURL(mediaUUID, lang string) string {
	return fmt.Sprintf(ttVideoURL, getLangTwoLetterCode(lang), mediaUUID)
}

func getLangTwoLetterCode(lang string) string {
	switch lang {
	case "en-US":
		return "en"
	case "hi-IN":
		return "hi"
	default:
		return "en"
	}
}
