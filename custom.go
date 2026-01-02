package main

import (
	"log"
)

func customizeCache(cache *videoCache) {
	// language for this video is English
	if video, exists := cache.get("1FVPtXv2pWU"); exists {
		video.Language = "en"
		cache.set(video)
		log.Printf("Updated video %s language to %s", video.VideoID, video.Language)
	}

	// deletes
	// Example: remove videos older than a certain date
	// cutoffDate := time.Now().AddDate(0, -6, 0) // 6 months ago
	// for id, video := range cache.Videos {
	//     if video.PublishedAt.Before(cutoffDate) {
	//         delete(cache.Videos, id)
	//     }
	// }
}
