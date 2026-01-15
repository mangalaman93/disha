package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSpotify(t *testing.T) {
	cache = videoCache{Videos: make(map[string]videoMeta)}
	if err := customizeSpotifyCache(&cache); err != nil {
		t.Fatal(err)
	}

	v, ok := cache.get("5PSCnndWS27XzNv43djH0g")
	assert.True(t, ok)
	assert.Equal(t, v.Name, "Tired of your hidden load?")
	assert.Equal(t, v.PublishYear, 2026)
	assert.Equal(t, v.PublishMonth, time.January)
	assert.Equal(t, v.PublishDay, 5)

	v, ok = cache.get("47qLeSG40eHeAXSsh6IhsH")
	assert.True(t, ok)
	assert.Equal(t, v.Name, "From Worrying to Thriving?")
	assert.Equal(t, v.PublishYear, 2026)
	assert.Equal(t, v.PublishMonth, time.January)
	assert.Equal(t, v.PublishDay, 13)
}
