package main

import (
	"log"
)

// customizeCache applies manual corrections to cached video data after it has
// been downloaded from external sources, allowing overrides of specific fields
// (such as language) before the cache is used elsewhere in the application.
func customizeCache(cache *videoCache) {
	// Language for this video is English.
	if video, exists := cache.get("1FVPtXv2pWU"); exists {
		video.Language = "en"
		cache.set(video)
		log.Printf("Updated video %s language to %s", video.VideoID, video.Language)
	}

	// Delete videos from the cache.
	// Example: remove videos older than a certain date
	// cutoffDate := time.Now().AddDate(0, -6, 0) // 6 months ago
	// for id, video := range cache.Videos {
	//     if video.PublishedAt.Before(cutoffDate) {
	//         delete(cache.Videos, id)
	//     }
	// }
}
