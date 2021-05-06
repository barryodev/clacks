package main

import (
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testURL = "www.example.com"

func TestLoadFeedData(t *testing.T) {
	fakeFeed := createTestFeed()

	parser := createStubbedParser(&fakeFeed, false)

	safeFeedData := &SafeFeedData{feedData: make(map[string]FeedDataModel)}
	allFeeds := &AllFeeds{}

	data := &Data{safeFeedData: safeFeedData, allFeeds: allFeeds}

	err := data.loadFeedData(testURL, parser)

	assert.Nil(t, err)
	feedDataModel := data.safeFeedData.GetEntries(testURL)
	assert.Equal(t, "Test Feed Title", feedDataModel.name)
	assert.Equal(t, 2, len(feedDataModel.entries))

	assert.Equal(t, "Test Entry Title One", feedDataModel.entries[0].title)
	assert.Equal(t, "Fake Entry Description One", feedDataModel.entries[0].content)
	assert.Equal(t, testURL+"/one", feedDataModel.entries[0].url)

	assert.Equal(t, "Test Entry Title Two", feedDataModel.entries[1].title)
	assert.Equal(t, "Fake Entry Description Two", feedDataModel.entries[1].content)
	assert.Equal(t, testURL+"/two", feedDataModel.entries[1].url)
}

func TestLoadFeedDataWithError(t *testing.T) {
	parser := createStubbedParser(nil, true)

	safeFeedData := &SafeFeedData{feedData: make(map[string]FeedDataModel)}
	allFeeds := &AllFeeds{}

	data := &Data{safeFeedData: safeFeedData, allFeeds: allFeeds}

	err := data.loadFeedData(testURL, parser)

	assert.NotNil(t, err)
	assert.Equal(t, "error loading feed: parser error", err.Error())
}

func createTestFeed() gofeed.Feed {
	fakeItemOne := &gofeed.Item{
		Title:       "Test Entry Title One",
		Description: "Fake Entry Description One",
		Link:        testURL+"/one",
	}
	fakeItemTwo := &gofeed.Item{
		Title:       "Test Entry Title Two",
		Description: "Fake Entry Description Two",
		Link:        testURL+"/two",
	}
	fakeFeed := gofeed.Feed{
		Title: "Test Feed Title",
		Items: []*gofeed.Item{fakeItemOne, fakeItemTwo},
	}
	return fakeFeed
}
