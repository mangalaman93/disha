package main

import (
	"log"
	"time"
)

// customizeCache applies manual corrections to cached video data after it has
// been downloaded from external sources, allowing overrides of specific fields
// (such as language) before the cache is used elsewhere in the application.
func customizeCache(cache *videoCache) {
	// Language for this video is English.
	// UXV4hcudGo0 & 1FVPtXv2pWU videos are in english
	toEnglish := []string{"UXV4hcudGo0", "1FVPtXv2pWU"}
	for _, videoID := range toEnglish {
		if video, exists := cache.get(videoID); exists {
			video.Language = englishLang
			cache.set(video)
			log.Printf("Updated video %s language to %s", video.VideoID, video.Language)
		}
	}

	// Add AajTak video
	video := videoMeta{
		VideoID: "01vCqZoMnyE",
		Name:    "Sahitya Aaj Tak 2025: स्वयं से साक्षात्कार | Prem Rawat | Sahitya Aaj Tak | Aaj Tak",
		Description: "दिल्ली की गुलाबी सर्दी के बीच मेजर ध्यानचंद स्टेडियम में आजतक के बेहद चर्चित कार्यक्रम साहित्य आजतक 2025 का आगाज हो चुका है." +
			"शुक्रवार को कार्यक्रम के पहले दिन न्यूयॉर्क टाइम्स के बेस्टसेलर लेखकों में शुमार प्रेम रावत ने आज की भागती-दौड़ती जिंदगी में शांति और आनंद का मतलब समझाया.",
		VideoDuration: time.Minute*43 + time.Second*49,
		Language:      hindiLang,
		ClickURL:      "https://www.youtube.com/watch?v=01vCqZoMnyE",
		PublishYear:   2025,
		PublishMonth:  time.November,
		ThumbnailURL:  "https://i.ytimg.com/vi/01vCqZoMnyE/hq720.jpg",
		AudioOnly:     false,
	}
	cache.set(video)

	// Delete videos from the cache.
	// Example: remove videos older than a certain date
	// cutoffDate := time.Now().AddDate(0, -6, 0) // 6 months ago
	// for id, video := range cache.Videos {
	//     if video.PublishedAt.Before(cutoffDate) {
	//         delete(cache.Videos, id)
	//     }
	// }
}
