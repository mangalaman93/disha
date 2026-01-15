package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

//go:embed data/spotify.html
var spotifyHTML []byte
var sentenceBoundary = regexp.MustCompile(`([.!?])([A-Z])`)
var spotifyPublishDate = time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC)

// customizeSpotifyCache adds Spotify podcast episodes to the cache
func customizeSpotifyCache(cache *videoCache) error {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(spotifyHTML))
	if err != nil {
		return fmt.Errorf("error parsing data/spotify.html: %w", err)
	}

	doc.Find(`[data-testid^="episode-"]`).Each(func(_ int, s *goquery.Selection) {
		vm, err := parseEpisode(s)
		if err != nil {
			panic(err)
		}
		if vm == nil {
			return
		}

		cache.set(*vm)
	})

	return nil
}

func parseEpisode(s *goquery.Selection) (*videoMeta, error) {
	link := s.Find(`a[href^="/episode/"]`).First()
	href, ok := link.Attr("href")
	if !ok {
		return nil, nil
	}

	videoID := strings.TrimPrefix(href, "/episode/")
	if videoID == "" {
		return nil, fmt.Errorf("no video ID found")
	}

	title := normalizeText(link.Text())
	if title == "" {
		return nil, fmt.Errorf("no title found")
	}

	description := normalizeText(
		s.Find(`[data-encore-id="listRowDetails"]`).Text(),
	)

	thumbnailURL := extractBestThumbnail(s)

	dateText := strings.TrimSpace(
		s.Find(".mVTn6QQxDKAnD1EMA4Uk [data-encore-id=text]").First().Text(),
	)
	publishedAt, err := parseDate(dateText)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %w", err)
	}

	durationText := strings.TrimSpace(
		s.Find(`[data-testid="episode-progress-not-played"]`).Text(),
	)
	duration, err := parseDurationForSpotify(durationText)
	if err != nil {
		return nil, fmt.Errorf("error parsing duration [%v]: %w", durationText, err)
	}

	return &videoMeta{
		VideoID:       videoID,
		Name:          title,
		Description:   description,
		VideoDuration: duration,
		Language:      englishLang,
		ClickURL:      "https://open.spotify.com/episode/" + videoID,
		PublishYear:   publishedAt.Year(),
		PublishMonth:  publishedAt.Month(),
		PublishDay:    publishedAt.Day(),
		ThumbnailURL:  thumbnailURL,
		AudioOnly:     true,
	}, nil
}

func normalizeText(s string) string {
	out := strings.Join(strings.Fields(s), " ")
	out = sentenceBoundary.ReplaceAllString(out, "$1 $2")
	return strings.TrimSpace(out)
}

func extractBestThumbnail(s *goquery.Selection) string {
	img := s.Find("img").First()

	if srcset, ok := img.Attr("srcset"); ok {
		candidates := strings.Split(srcset, ",")
		for _, c := range candidates {
			c = strings.TrimSpace(c)
			if strings.HasSuffix(c, "640w") {
				parts := strings.Split(c, " ")
				if len(parts) > 0 {
					return parts[0]
				}
			}
		}
	}

	if src, ok := img.Attr("src"); ok {
		return src
	}

	return ""
}
func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)

	// 1. Handle weekday names (e.g. "Tuesday")
	if wd, ok := parseWeekday(s); ok {
		return mostRecentWeekday(spotifyPublishDate, wd), nil
	}

	// 2. Absolute date formats
	layouts := []string{
		"Jan 2, 2006",
		"Jan 2006",
	}

	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unrecognized date: %q", s)
}

func parseWeekday(s string) (time.Weekday, bool) {
	switch strings.ToLower(s) {
	case "sunday":
		return time.Sunday, true
	case "monday":
		return time.Monday, true
	case "tuesday":
		return time.Tuesday, true
	case "wednesday":
		return time.Wednesday, true
	case "thursday":
		return time.Thursday, true
	case "friday":
		return time.Friday, true
	case "saturday":
		return time.Saturday, true
	default:
		return 0, false
	}
}

func mostRecentWeekday(from time.Time, wd time.Weekday) time.Time {
	// Normalize to local midnight
	base := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())

	diff := (int(base.Weekday()) - int(wd) + 7) % 7
	return base.AddDate(0, 0, -diff)
}

func parseDurationForSpotify(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)

	var h, m, sec int
	var n int
	var err error

	// "1 hr 12 min"
	if strings.Contains(s, "hr") && strings.Contains(s, "min") {
		n, err = fmt.Sscanf(s, "%d hr %d min", &h, &m)
		if err != nil || n != 2 {
			return 0, fmt.Errorf("invalid duration format: %q", s)
		}
		return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute, nil
	}

	// "1 hr"
	if strings.Contains(s, "hr") {
		n, err = fmt.Sscanf(s, "%d hr", &h)
		if err != nil || n != 1 {
			return 0, fmt.Errorf("invalid duration format: %q", s)
		}
		return time.Duration(h) * time.Hour, nil
	}

	// "22 min 39 sec"
	if strings.Contains(s, "sec") {
		n, err = fmt.Sscanf(s, "%d min %d sec", &m, &sec)
		if err != nil || n != 2 {
			return 0, fmt.Errorf("invalid duration format: %q", s)
		}
		return time.Duration(m)*time.Minute + time.Duration(sec)*time.Second, nil
	}

	// "37 min"
	n, err = fmt.Sscanf(s, "%d min", &m)
	if err != nil || n != 1 {
		return 0, fmt.Errorf("invalid duration format: %q", s)
	}
	return time.Duration(m) * time.Minute, nil
}
