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
func LoadJsonConfig() (*AllFeeds, error){
	//open config file
	configFile, err := os.Open("feeds.json")
	if err != nil {
		return nil, errors.New("error: could not find feeds.json config file")
	}

	//load data from file
	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, errors.New("Error Reading from feeds.json file: " + err.Error())
	}

	//unmarshall json
	var loadedFeeds AllFeeds
	err = json.Unmarshal(byteValue, &loadedFeeds)
	if err != nil {
		return nil, errors.New("error: reading feeds.json, please check structure")
	}

	return &loadedFeeds, nil
}

// LoadFeedData use gofeed library to load data from atom feed
func LoadFeedData(url string) (FeedDataModel, error) {
	fp := gofeed.NewParser()
	fp.UserAgent = "Clacks - Terminal Atom Reader"
	feedData, err := fp.ParseURL(url)
	if err != nil {
		return FeedDataModel{}, errors.New("error loading feed: " + err.Error())
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
		return FeedDataModel{name: feedName, entries: entrySlice}, nil
	}

	return FeedDataModel{}, nil
}

// This will run asynchronously to fetch atom feeds and load the data into the interface
func LoadAllFeedDataAndUpdateInterface() {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			app.QueueUpdateDraw(func() {
				CreateErrorPage(err.Error())
			})
		}
	}()
	var configError error
	allFeeds, configError = LoadJsonConfig()
	if configError != nil {
		panic(configError)
	}

	for i := 0; i < len(allFeeds.Feeds); i++ {
		feedData, atomFeedError := LoadFeedData(allFeeds.Feeds[i].URL)
		if atomFeedError != nil {
			panic(atomFeedError)
		} else {
			safeFeedData.SetSiteData(allFeeds.Feeds[i].URL, feedData)
		}
	}

	LoadFeedDataIntoLists()
}
