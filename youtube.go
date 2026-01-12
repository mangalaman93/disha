package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	rvkHandle  = "@rajvidyakender"
	ttHandle   = "@TimelessToday"
	prHandle   = "@PremRawatOfficial"
	wopgHandle = "@wopgyt"

	youTubeVideoURL   = "https://www.youtube.com/watch?v=%v"
	baseYouTubeAPIURL = "https://www.googleapis.com/youtube/v3"
	playlistURL       = baseYouTubeAPIURL + "/channels?part=contentDetails&forHandle=%v&key=%v"
	videoListURL      = baseYouTubeAPIURL + "/playlistItems?part=snippet&maxResults=50&playlistId=%v&key=%v&pageToken=%v"
	videoMetaURL      = baseYouTubeAPIURL + "/videos?part=snippet,contentDetails&id=%v&key=%v"
)

var (
	allYtHandles = []string{wopgHandle, prHandle, rvkHandle, ttHandle}
)

func getYouTubeContent() ([]videoMeta, error) {
	var videos []videoMeta
	for _, handle := range allYtHandles {
		log.Printf("getting videos from handle: [%v]\n", handle)

		playlistID, err := getPlaylistID(handle)
		if err != nil {
			return nil, fmt.Errorf("error getting playlist ID for [%v]: %v", handle, err)
		}

		playlistVideos, err := getVideosFromPlaylist(playlistID)
		if err != nil {
			return nil, fmt.Errorf("error getting videos from playlist [%v]: %v", playlistID, err)
		}

		videos = append(videos, playlistVideos...)
	}

	return videos, nil
}

func getPlaylistID(handle string) (string, error) {
	youTubeAPIKey := os.Getenv("YOUTUBE_API_KEY")
	resp, err := httpClient.Get(fmt.Sprintf(playlistURL, handle, youTubeAPIKey))
	if err != nil {
		return "", fmt.Errorf("error getting playlist ID for [%v]: %v", handle, err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting playlist ID for [%v]: %v", handle, resp.Status)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("error while closing response body for handle %s: %v", handle, err)
		}
	}()

	var respstruct struct {
		Items []struct {
			ContentDetails struct {
				RelatedPlaylists struct {
					Uploads string `json:"uploads"`
				} `json:"relatedPlaylists"`
			} `json:"contentDetails"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respstruct); err != nil {
		return "", fmt.Errorf("error decoding response while getting playlist ID for [%v]: %v", handle, err)
	}
	if len(respstruct.Items) == 0 {
		return "", fmt.Errorf("no playlist found for [%v]", handle)
	}

	return respstruct.Items[0].ContentDetails.RelatedPlaylists.Uploads, nil
}

func getVideosFromPlaylist(playlistID string) ([]videoMeta, error) {
	youTubeAPIKey := os.Getenv("YOUTUBE_API_KEY")

	var videos []videoMeta
	nextPageToken := ""
	pageNo := 0
	for {
		pageNo++
		log.Printf("getting page [%v] of videos from playlist [%v], nextPageToken: [%v]\n",
			pageNo, playlistID, nextPageToken)

		resp, err := httpClient.Get(fmt.Sprintf(videoListURL, playlistID, youTubeAPIKey, nextPageToken))
		if err != nil {
			return nil, fmt.Errorf("error getting videos from playlist [%v]: %v", playlistID, err)
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error getting videos from playlist [%v]: %v", playlistID, resp.Status)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("error while closing response body for playlist %s: %v", playlistID, err)
			}
		}()

		var respstruct struct {
			NextPageToken string `json:"nextPageToken"`
			Items         []struct {
				Snippet struct {
					Title       string `json:"title"`
					Description string `json:"description"`
					PublishedAt string `json:"publishedAt"`
					Thumbnails  struct {
						Medium struct {
							URL string `json:"url"`
						} `json:"medium"`
					} `json:"thumbnails"`
					ResourceID struct {
						VideoID string `json:"videoId"`
					} `json:"resourceId"`
				} `json:"snippet"`
			} `json:"items"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&respstruct); err != nil {
			return nil, fmt.Errorf("error decoding response in getting videos from playlist [%v]: %v", playlistID, err)
		}
		if len(respstruct.Items) == 0 {
			break
		}

		for _, item := range respstruct.Items {
			video, ok := cache.get(item.Snippet.ResourceID.VideoID)
			if ok {
				videos = append(videos, video)
				continue
			}

			publishTs, err := time.Parse("2006-01-02T15:04:05Z", item.Snippet.PublishedAt)
			if err != nil {
				return nil, fmt.Errorf("error parsing publish date for video [%+v]: %w", item, err)
			}

			audioLang, duration, err := getMetaForYouTubeVideo(item.Snippet.ResourceID.VideoID)
			if err != nil {
				return nil, err
			}

			if duration == 0 || (audioLang != "hi-IN" && audioLang != "en-US") {
				continue
			}

			videos = append(videos, videoMeta{
				VideoID:       item.Snippet.ResourceID.VideoID,
				Name:          item.Snippet.Title,
				Description:   item.Snippet.Description,
				VideoDuration: duration,
				Language:      audioLang,
				ClickURL:      fmt.Sprintf(youTubeVideoURL, item.Snippet.ResourceID.VideoID),
				PublishYear:   publishTs.Year(),
				PublishMonth:  publishTs.Month(),
				PublishDay:    publishTs.Day(),
				ThumbnailURL:  item.Snippet.Thumbnails.Medium.URL,
				AudioOnly:     false,
			})
		}

		nextPageToken = respstruct.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

	return videos, nil
}

func getMetaForYouTubeVideo(videoID string) (string, time.Duration, error) {
	youTubeAPIKey := os.Getenv("YOUTUBE_API_KEY")
	resp, err := httpClient.Get(fmt.Sprintf(videoMetaURL, videoID, youTubeAPIKey))
	if err != nil {
		return "", 0, fmt.Errorf("error getting meta for video [%v]: %v", videoID, err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("error getting meta for video [%v]: %v", videoID, resp.Status)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("error while closing response body for video %s: %v", videoID, err)
		}
	}()

	var respstruct struct {
		Items []struct {
			Snippet struct {
				Title     string `json:"title"`
				AudioLang string `json:"defaultAudioLanguage"`
			} `json:"snippet"`
			ContentDetails struct {
				Duration string `json:"duration"`
			} `json:"contentDetails"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respstruct); err != nil {
		return "", 0, fmt.Errorf("error decoding response while getting meta for video [%v]: %v", videoID, err)
	}
	if len(respstruct.Items) == 0 {
		return "", 0, fmt.Errorf("no meta found for video [%v]", videoID)
	}

	durationStr := respstruct.Items[0].ContentDetails.Duration
	duration, err := parseDuration(durationStr)
	if err != nil {
		return "", 0, fmt.Errorf("error parsing duration [%v] for video [%v]: %v", durationStr, videoID, err)
	}

	return langTT(respstruct.Items[0].Snippet.AudioLang, respstruct.Items[0].Snippet.Title), duration, nil
}

func parseDuration(duration string) (time.Duration, error) {
	if duration == "P0D" {
		return 0, nil
	}

	return time.ParseDuration(strings.ToLower(duration[2:]))
}

func langTT(lang, title string) string {
	titleHasHindi := containsHindi(title)
	if titleHasHindi {
		return "hi-IN"
	}

	switch lang {
	case "hi", "hi-IN":
		return "hi-IN"

	case "en-GB", "en-US", "en":
		return "en-US"

	default:
		return lang
	}
}

func containsHindi(title string) bool {
	for _, r := range title {
		// Devanagari characters are generally in the range 0x0900 to 0x097F
		if r >= '\u0900' && r <= '\u097F' {
			return true
		}
	}
	return false
}
