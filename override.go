package main

import (
	"log"
	"time"
)

// customizeCache applies manual corrections to cached video data after it has
// been downloaded from external sources, allowing overrides of specific fields
// (such as language) before the cache is used elsewhere in the application.
func customizeCache(cache *videoCache) error {
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
	cache.set(videoMeta{
		VideoID: "01vCqZoMnyE",
		Name:    "Sahitya Aaj Tak 2025: स्वयं से साक्षात्कार | Prem Rawat | Sahitya Aaj Tak | Aaj Tak",
		Description: "दिल्ली की गुलाबी सर्दी के बीच मेजर ध्यानचंद स्टेडियम में आजतक के बेहद चर्चित " +
			"कार्यक्रम साहित्य आजतक 2025 का आगाज हो चुका है." +
			"शुक्रवार को कार्यक्रम के पहले दिन न्यूयॉर्क टाइम्स के बेस्टसेलर लेखकों में शुमार प्रेम रावत ने " +
			"आज की भागती-दौड़ती जिंदगी में शांति और आनंद का मतलब समझाया.",
		VideoDuration: time.Minute*43 + time.Second*49,
		Language:      hindiLang,
		ClickURL:      "https://www.youtube.com/watch?v=01vCqZoMnyE",
		PublishYear:   2025,
		PublishMonth:  time.November,
		PublishDay:    21,
		ThumbnailURL:  "https://i.ytimg.com/vi/01vCqZoMnyE/hq720.jpg",
		AudioOnly:     false,
	})

	// Add https://www.youtube.com/watch?v=_J_KLm4kj-Y
	cache.set(videoMeta{
		VideoID: "_J_KLm4kj-Y",
		Name:    "ये 3 कानून ज़िंदगी बदल देंगे... Prem Rawat से जानें सत्य क्या है?. देखें पूरा वीडियो सिर्फ साहित्य तक पर.",
		Description: "Sahitya Tak Podcast, ये 3 कानून ज़िंदगी बदल देंगे... Prem Rawat से जानें सत्य क्या है?. " +
			"देखें पूरा वीडियो सिर्फ साहित्य तक पर.",
		VideoDuration: time.Minute*10 + time.Second*18,
		Language:      hindiLang,
		ClickURL:      "https://www.youtube.com/watch?v=_J_KLm4kj-Y",
		PublishYear:   2025,
		PublishMonth:  time.March,
		PublishDay:    4,
		ThumbnailURL:  "https://i.ytimg.com/vi/_J_KLm4kj-Y/hq720.jpg",
		AudioOnly:     false,
	})

	// Add https://www.youtube.com/watch?v=JW9W31HLiH0
	cache.set(videoMeta{
		VideoID:       "JW9W31HLiH0",
		Name:          "जीवन! बस इन 3 नियमों पर टिका | Prem Rawat से Breath: Wake Up to Life पर बतकही | EP 109 | Sahitya Tak",
		Description:   "",
		VideoDuration: time.Minute*40 + time.Second*23,
		Language:      hindiLang,
		ClickURL:      "https://www.youtube.com/watch?v=JW9W31HLiH0",
		PublishYear:   2025,
		PublishMonth:  time.February,
		PublishDay:    23,
		ThumbnailURL:  "https://i.ytimg.com/vi/JW9W31HLiH0/hq720.jpg",
		AudioOnly:     false,
	})

	// Add https://www.youtube.com/watch?v=bLzvopaMLwk
	cache.set(videoMeta{
		VideoID: "bLzvopaMLwk",
		Name:    "Prem Rawat | Peace Education Keynote | Global Peace Education Day",
		Description: `What skills and knowledge do we need to build a culture of peace on a ` +
			`healthy planet? A keynote address from Prem Rawat, Author; Founder of ` +
			`The Prem Rawat Foundation and the Peace Education Program.`,
		VideoDuration: time.Minute*12 + time.Second*24,
		Language:      englishLang,
		ClickURL:      "https://www.youtube.com/watch?v=bLzvopaMLwk",
		PublishYear:   2022,
		PublishMonth:  time.October,
		PublishDay:    22,
		ThumbnailURL:  "https://i.ytimg.com/vi/bLzvopaMLwk/hqdefault.jpg",
		AudioOnly:     false,
	})

	// Add https://www.youtube.com/watch?v=djd5THkx7Hs
	cache.set(videoMeta{
		VideoID:       "djd5THkx7Hs",
		Name:          "विश्वास की बजाय अनुभव को चुनिए... पथ प्रदर्शक Prem Rawat | 'स्वयं की आवाज़' पर चर्चा | Sahitya Tak",
		Description:   "",
		VideoDuration: time.Minute*21 + time.Second*15,
		Language:      hindiLang,
		ClickURL:      "https://www.youtube.com/watch?v=djd5THkx7Hs",
		PublishYear:   2023,
		PublishMonth:  time.April,
		PublishDay:    1,
		ThumbnailURL:  "https://i.ytimg.com/vi/djd5THkx7Hs/hq720.jpg",
		AudioOnly:     false,
	})

	// Add https://www.youtube.com/watch?v=tMe7_9GSXEM
	cache.set(videoMeta{
		VideoID:       "tMe7_9GSXEM",
		Name:          "SPECIAL INTERVIEW WITH GLOBAL PEACE AMBASSADOR PREM RAWAT",
		Description:   "",
		VideoDuration: time.Minute*27 + time.Second*48,
		Language:      hindiLang,
		ClickURL:      "https://www.youtube.com/watch?v=tMe7_9GSXEM",
		PublishYear:   2018,
		PublishMonth:  time.November,
		PublishDay:    2,
		ThumbnailURL:  "https://i.ytimg.com/vi/tMe7_9GSXEM/hqdefault.jpg",
		AudioOnly:     false,
	})

	// Add https://www.youtube.com/watch?v=zCuKz6M-hTo
	cache.set(videoMeta{
		VideoID:       "zCuKz6M-hTo",
		Name:          "आखिर किसकी सुनें…दिल की या मन की ? अंतर्राष्ट्रीय वक्ता और शांति दूत Prem Rawat EXCLUSIVE | Asha Jha",
		Description:   "",
		VideoDuration: time.Minute*25 + time.Second*41,
		Language:      hindiLang,
		ClickURL:      "https://www.youtube.com/watch?v=zCuKz6M-hTo",
		PublishYear:   2025,
		PublishMonth:  time.April,
		PublishDay:    13,
		ThumbnailURL:  "https://i.ytimg.com/vi/zCuKz6M-hTo/hq720.jpg",
		AudioOnly:     false,
	})

	// Add https://www.youtube.com/watch?v=-vyRZwCsn9I
	cache.set(videoMeta{
		VideoID:       "-vyRZwCsn9I",
		Name:          "Jail की सज़ा काट रही इस औरत के लिए कोई Hope है? Prem Rawat ने क्या कहा | Prem Rawat Interview",
		Description:   "",
		VideoDuration: time.Minute*9 + time.Second*57,
		Language:      hindiLang,
		ClickURL:      "https://www.youtube.com/watch?v=-vyRZwCsn9I",
		PublishYear:   2022,
		PublishMonth:  time.March,
		PublishDay:    9,
		ThumbnailURL:  "https://i.ytimg.com/vi/-vyRZwCsn9I/hq720.jpg",
		AudioOnly:     false,
	})

	// Add https://www.youtube.com/watch?v=4TVaZCbEpWs
	cache.set(videoMeta{
		VideoID:       "4TVaZCbEpWs",
		Name:          "Prem Rawat Life Story | 4 वर्ष की उम्र में जिन्होंने रोक दी भीड़, 12 की उम्र में England को लिया लुभा",
		Description:   "",
		VideoDuration: time.Minute*36 + time.Second*18,
		Language:      hindiLang,
		ClickURL:      "https://www.youtube.com/watch?v=4TVaZCbEpWs",
		PublishYear:   2023,
		PublishMonth:  time.May,
		PublishDay:    28,
		ThumbnailURL:  "https://i.ytimg.com/vi/4TVaZCbEpWs/hq720.jpg",
		AudioOnly:     false,
	})

	// Add https://www.youtube.com/watch?v=4xN4SDjbpjI
	cache.set(videoMeta{
		VideoID:       "4xN4SDjbpjI",
		Name:          "लेखक , मानवतावादी 'प्रेम रावत' से खास बातचीत संजय गिरि गोस्वामी के साथ |Network10|PREM RAWAT PODCAST",
		Description:   "",
		VideoDuration: time.Minute*41 + time.Second*17,
		Language:      hindiLang,
		ClickURL:      "https://www.youtube.com/watch?v=4xN4SDjbpjI",
		PublishYear:   2023,
		PublishMonth:  time.July,
		PublishDay:    25,
		ThumbnailURL:  "https://i.ytimg.com/vi/4xN4SDjbpjI/hq720.jpg",
		AudioOnly:     true,
	})

	// Add https://www.youtube.com/watch?v=MBnMKUE8bFo
	cache.set(videoMeta{
		VideoID:       "MBnMKUE8bFo",
		Name:          "Aaj Savere - An interview with - Sh. Prem Rawat, International Peace Speaker",
		Description:   "",
		VideoDuration: time.Minute*48 + time.Second*55,
		Language:      englishLang,
		ClickURL:      "https://www.youtube.com/watch?v=MBnMKUE8bFo",
		PublishYear:   2017,
		PublishMonth:  time.November,
		PublishDay:    24,
		ThumbnailURL:  "https://i.ytimg.com/vi/MBnMKUE8bFo/hq720.jpg",
		AudioOnly:     false,
	})

	// Add https://www.youtube.com/watch?v=k845byCwFWg
	cache.set(videoMeta{
		VideoID: "k845byCwFWg",
		Name:    "Prem Rawat | Hear Yourself: How to Find Peace in a Noisy World | Talks at Google",
		Description: `Renowned teacher and author Prem Rawat discusses his book "Hear Yourself: How to Find Peace ` +
			`in a Noisy World", where he teaches us how to turn down the noise to “hear ourselves”—to listen to ` +
			`the subtle song of peace that sings inside each of us. Once we learn to truly “hear ourselves” and ` +
			`the voice of peace within, then we can hold on to that as we face all the noise of the world.\n\n` +
			`Prem Rawat is the founder of the Prem Rawat Foundation, where he works with people from all walks of ` +
			`life, showing them how to experience the source of peace within themselves. His work spans six decades ` +
			`of international effort to bring a practical message of hope, happiness, and peace to all, one person ` +
			`at a time. He is the internationally bestselling author of Peace Is Possible and is also a pilot, ` +
			`photographer, classic car restorer, father of four children and grandfather of four.`,
		VideoDuration: time.Minute*51 + time.Second*35,
		Language:      englishLang,
		ClickURL:      "https://www.youtube.com/watch?v=k845byCwFWg",
		PublishYear:   2022,
		PublishMonth:  time.May,
		PublishDay:    30,
		ThumbnailURL:  "https://i.ytimg.com/vi/k845byCwFWg/hq720.jpg",
		AudioOnly:     false,
	})

	return nil
}
