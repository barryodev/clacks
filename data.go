package main

import (
	"encoding/json"
	"errors"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"
	"io"

	"html"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

// Data struct holding config and feed data
type Data struct {
	safeFeedData	*SafeFeedData
	configData		*ConfigData
	parser			FeedParser
}

type FeedParser interface {
	ParseURL(feedURL string) (feed *gofeed.Feed, err error)
}

const configFileName = "feeds.json"

// ConfigData & Feed struct to unmarshall json config
type ConfigData struct {
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

// Clear clear all feed data, use before a refresh
func (c *SafeFeedData) Clear() {
	c.mu.Lock()
	// Lock so only one goroutine at a time can access the map
	c.feedData = make(map[string]FeedDataModel)
	c.mu.Unlock()
}

// LoadJsonConfig load feeds from json file
func (data *Data) loadJsonConfig(fileName string) error {
	//open config file
	configFile, fileError := openConfigFile(fileName)
	if fileError != nil {
		return fileError
	}

	config, parseError := parseConfig(configFile)
	if parseError != nil {
		return parseError
	} else {
		data.configData = config
		return nil
	}
}

func parseConfig(r io.Reader) (*ConfigData, error) {
	//load data from file
	byteValue, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.New("error reading from feeds.json file " + err.Error())
	}

	//unmarshall json
	var loadedFeeds ConfigData
	err = json.Unmarshal(byteValue, &loadedFeeds)
	if err != nil {
		return nil, errors.New("error reading feeds.json, please check structure")
	}

	return &loadedFeeds, nil
}

func openConfigFile(fileName string) (*os.File, error) {
	configFile, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("error: could not find feeds.json config file")
	}
	return configFile, err
}

// LoadFeedData use gofeed library to load data from atom feed
func (data *Data) loadFeedData(url string) error {
	feedData, err := data.parser.ParseURL(url)
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
	} else {
		return errors.New("error feed at url: " + url + " has no entries")
	}
}

// asynchronously fetch atom feeds and load the data into the interface
func (data *Data) loadDataFromFeeds() error {
	if data.configData == nil || len(data.configData.Feeds) == 0 {
		return errors.New("error attempted to load feed data with no config set")
	}

	for i := 0; i < len(data.configData.Feeds); i++ {
		atomFeedError := data.loadFeedData(data.configData.Feeds[i].URL)
		if atomFeedError != nil {
			return atomFeedError
		}
	}

	return nil
}

// NewData factory method for data objects
func NewData(parser FeedParser) *Data {
	safeFeedData := &SafeFeedData{feedData: make(map[string]FeedDataModel)}
	allFeeds := &ConfigData{}
	data := &Data{safeFeedData: safeFeedData, configData: allFeeds, parser: parser}
	return data
}
