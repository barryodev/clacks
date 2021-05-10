package main

import (
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadFeedData(t *testing.T) {
	fakeFeed := CreateTestFeed()
	parser := createStubbedParser(&fakeFeed, false)

	data := NewData(parser)
	err := data.loadFeedData(testUrlOne)

	assert.Nil(t, err)
	feedDataModel := data.safeFeedData.GetEntries(testUrlOne)
	assert.Equal(t, fakeFeed.Title, feedDataModel.name)
	assert.Equal(t, 2, len(feedDataModel.entries))

	assert.Equal(t, fakeFeed.Items[0].Title, feedDataModel.entries[0].title)
	assert.Equal(t, fakeFeed.Items[0].Description, feedDataModel.entries[0].content)
	assert.Equal(t, fakeFeed.Items[0].Link, feedDataModel.entries[0].url)

	assert.Equal(t, fakeFeed.Items[1].Title, feedDataModel.entries[1].title)
	assert.Equal(t, fakeFeed.Items[1].Description, feedDataModel.entries[1].content)
	assert.Equal(t, fakeFeed.Items[1].Link, feedDataModel.entries[1].url)
}

func TestLoadFeedDataWithError(t *testing.T) {
	parser := createStubbedParser(nil, true)

	data := NewData(parser)
	err := data.loadFeedData(testUrlOne)

	assert.NotNil(t, err)
	assert.Equal(t, "error loading feed: parser error", err.Error())
}

func CreateTestFeed() gofeed.Feed {
	fakeItemOne := &gofeed.Item{
		Title:       "Test Entry Title One",
		Description: "Fake Entry Description One",
		Link:        testUrlOne+"/one",
	}
	fakeItemTwo := &gofeed.Item{
		Title:       "Test Entry Title Two",
		Description: "Fake Entry Description Two",
		Link:        testUrlOne+"/two",
	}
	fakeFeed := gofeed.Feed{
		Title: "Test Feed Title From Parser",
		Items: []*gofeed.Item{fakeItemOne, fakeItemTwo},
	}
	return fakeFeed
}
