package main

import (
	"encoding/json"
	"errors"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"

	"html"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

// Data struct holding config and feed data
type Data struct {
	safeFeedData	*SafeFeedData
	allFeeds		*AllFeeds
}

type FeedParser interface {
	ParseURL(feedURL string) (feed *gofeed.Feed, err error)
}

const configFileName = "feeds.json"

// AllFeeds & Feed struct to unmarshall json config
type AllFeeds struct {
	Feeds []Feed `json:"feeds"`
}

type Feed struct {
	URL string `json:"url"`
}

// Entry struct describes a single item in an atom feed
type Entry struct {
	title 	string
	content string
	url		string
}

// FeedDataModel struct of an feed title and slice of entries
type FeedDataModel struct {
	name string
	entries []Entry
}

// SafeFeedData Map of Urls Strings to []Entry with mutex for thread safety
type SafeFeedData struct {
	mu       sync.Mutex
	feedData map[string]FeedDataModel
}

// SetSiteData Adding stat strings to the map
func (c *SafeFeedData) SetSiteData(url string, data FeedDataModel) {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map
	c.feedData[url] = data
	c.mu.Unlock()
}

// GetEntries get a value from map
func (c *SafeFeedData) GetEntries(url string) FeedDataModel {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map
	defer c.mu.Unlock()
	return c.feedData[url]
}

// LoadJsonConfig load feeds from json file
func (data *Data) loadJsonConfig(config string) error {
	//open config file
	configFile, err := os.Open(config)
	if err != nil {
		return errors.New("error: could not find feeds.json config file")
	}

	//load data from file
	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		return errors.New("Error Reading from feeds.json file: " + err.Error())
	}

	//unmarshall json
	var loadedFeeds AllFeeds
	err = json.Unmarshal(byteValue, &loadedFeeds)
	if err != nil {
		return errors.New("error: reading feeds.json, please check structure")
	}

	data.allFeeds = &loadedFeeds
	return nil
}

// LoadFeedData use gofeed library to load data from atom feed
func (data *Data) loadFeedData(url string, parser FeedParser) error {
	feedData, err := parser.ParseURL(url)
	if err != nil {
		return errors.New("error loading feed: " + err.Error())
	}

	if len(feedData.Items) > 0 {
		feedName := feedData.Title
		entrySlice := make([]Entry, len(feedData.Items))
		for i, item := range feedData.Items {
			entrySlice[i] = Entry{
				title: html.UnescapeString(strip.StripTags(item.Title)),
				content: strings.TrimSpace(html.UnescapeString(strip.StripTags(item.Description))),
				url: item.Link,
			}
		}
		feedDataModel := FeedDataModel{name: feedName, entries: entrySlice}
		data.safeFeedData.SetSiteData(url, feedDataModel)
		return nil
	}

	//TODO handle feed having no entries
	return nil
}

// asynchronously fetch atom feeds and load the data into the interface
func (data *Data) loadAllFeedData(parser FeedParser) error {
	var configError error
	configError = data.loadJsonConfig(configFileName)
	if configError != nil {
		return configError
	}

	for i := 0; i < len(data.allFeeds.Feeds); i++ {
		atomFeedError := data.loadFeedData(data.allFeeds.Feeds[i].URL, parser)
		if atomFeedError != nil {
			return atomFeedError
		}
	}

	return nil
}
