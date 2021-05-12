package main

import (
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const badConfigFile = "malformed_config.json"

func TestLoadFeedData(t *testing.T) {
	fakeFeed := CreateTestFeed()
	parser := createStubbedParser(&fakeFeed, false)

	testFeed := Feed{URL: testURLOne}
	testConfig := ConfigData{[]Feed{testFeed}}

	data := NewData(parser)
	data.configData = &testConfig
	err := data.loadDataFromFeeds()

	assert.Nil(t, err)
	feedDataModel := data.safeFeedData.GetEntries(testURLOne)
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
	err := data.loadFeedData(testURLOne)

	assert.NotNil(t, err)
	assert.Equal(t, "error loading feed: stubbed parser error", err.Error())
}

func TestLoadFeedDataWithoutConfig(t *testing.T) {
	parser := createStubbedParser(nil, false)
	data := NewData(parser)

	err := data.loadDataFromFeeds()

	assert.NotNil(t, err)
	assert.Equal(t, "error attempted to load feed data with no config set", err.Error())
}

func TestLoadFeedDataWithFeedWithNoEntries(t *testing.T) {
	fakeFeedWithNoEntries := gofeed.Feed{
		Title: "Test Feed Title From Parser",
		Items: []*gofeed.Item{},
	}
	testFeed := Feed{URL: testURLOne}
	testConfig := ConfigData{[]Feed{testFeed}}

	parser := createStubbedParser(&fakeFeedWithNoEntries, false)
	data := NewData(parser)

	data.configData = &testConfig

	err := data.loadDataFromFeeds()

	assert.NotNil(t, err)
	assert.Equal(t, "error feed at url: "+testURLOne+" has no entries", err.Error())
}

func TestLoadJsonConfig(t *testing.T) {
	parser := createStubbedParser(nil, false)
	data := NewData(parser)

	err := data.loadJSONConfig(configFileName)

	assert.Nil(t, err)
	assert.Equal(t, 4, len(data.configData.Feeds))
	assert.Equal(t, "https://www.theregister.com/offbeat/bofh/headlines.atom",
		data.configData.Feeds[0].URL)
}

func TestLoadJsonConfigWithMissingFile(t *testing.T) {
	parser := createStubbedParser(nil, false)
	data := NewData(parser)

	err := data.loadJSONConfig("badfilename.json")

	assert.NotNil(t, err)
	assert.Equal(t, "error: could not find feeds.json config file", err.Error())
}

func TestLoadJsonConfigReturnsErrorFromParse(t *testing.T) {
	parser := createStubbedParser(nil, true)
	data := NewData(parser)

	err := data.loadJSONConfig(badConfigFile)

	assert.NotNil(t, err)
	assert.Equal(t, "error reading feeds.json, please check structure", err.Error())
}

func TestParseConfigWithGoodJson(t *testing.T) {
	goodJSON := `{
				  "feeds": [
					{
					  "url": "https://blah.com/rss"
					},
					{
					  "url": "https://example.com/rss"
					}
				  ]
				}`

	r := strings.NewReader(goodJSON)

	configData, err := parseConfig(r)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(configData.Feeds))
	assert.Equal(t, "https://blah.com/rss", configData.Feeds[0].URL)
	assert.Equal(t, "https://example.com/rss", configData.Feeds[1].URL)
}

func TestParseConfigWithMalformedJson(t *testing.T) {
	badJSON := `{
				  "feasdfasdfds": [
					{
					  "uasdfasdfrl": "https://blah.com/rss"
					},
					{
					  "url": "https://example.com/rss"
					}
				  ]
				asdfasdff}`

	r := strings.NewReader(badJSON)

	configData, err := parseConfig(r)

	assert.Nil(t, configData)
	assert.NotNil(t, err)
	assert.Equal(t, "error reading feeds.json, please check structure", err.Error())
}

func TestParseConfigWithEmptyReader(t *testing.T) {
	r := StubbedBuffer{}
	configData, err := parseConfig(r)

	assert.Nil(t, configData)
	assert.NotNil(t, err)
	assert.Equal(t, "error reading from feeds.json file stubbed buffer error", err.Error())
}

func CreateTestFeed() gofeed.Feed {
	fakeItemOne := &gofeed.Item{
		Title:       "Test Entry Title One",
		Description: "Fake Entry Description One",
		Link:        testURLOne + "/one",
	}
	fakeItemTwo := &gofeed.Item{
		Title:       "Test Entry Title Two",
		Description: "Fake Entry Description Two",
		Link:        testURLOne + "/two",
	}
	fakeFeed := gofeed.Feed{
		Title: "Test Feed Title From Parser",
		Items: []*gofeed.Item{fakeItemOne, fakeItemTwo},
	}
	return fakeFeed
}
