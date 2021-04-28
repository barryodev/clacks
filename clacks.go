package main

import (
	"encoding/json"
	"errors"
	"github.com/icza/gox/osx"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"
	"github.com/rivo/tview"
	"golang.org/x/net/html"
)

// tview ui elements, config and feed data
var (
	app 			*tview.Application
	feedList		*tview.List
	entriesList 	*tview.List
	entryTextView 	*tview.TextView
	flex			*tview.Flex
	safeFeedData	*SafeFeedData
	allFeeds		*AllFeeds
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

// setSiteData Adding stat strings to the map
func (c *SafeFeedData) setSiteData(url string, data FeedDataModel) {
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

// loadJsonConfig load feeds from json file
func loadJsonConfig() (*AllFeeds, error){
	//open config file
	configFile, err := os.Open("feeds.json")
	if err != nil {
		return nil, errors.New("Error loading feeds.json file: " + err.Error())
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
		return nil, errors.New("Error unmarshalling json: " + err.Error())
	}

	return &loadedFeeds, nil
}

// loadFeedData use gofeed library to load data from atom feed
func loadFeedData(url string) (FeedDataModel, error) {
	fp := gofeed.NewParser()
	fp.UserAgent = "Clacks - Terminal Atom Reader"
	feedData, err := fp.ParseURL(url)
	if err != nil {
		return FeedDataModel{}, errors.New("Error processing atom feed: " + err.Error())
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
func loadAllFeedDataAndUpdateInterface() {
	var configError error
	allFeeds, configError = loadJsonConfig()
	if configError != nil {
		panic(configError)
	}

	for i := 0; i < len(allFeeds.Feeds); i++ {
		feedData, atomFeedError := loadFeedData(allFeeds.Feeds[i].URL)
		if atomFeedError != nil {
			panic(atomFeedError)
		} else {
			safeFeedData.setSiteData(allFeeds.Feeds[i].URL, feedData)
		}
	}

	// Using QueueUpdateDraw as its a threadsafe way to update tview primitives
	app.QueueUpdateDraw(func() {
		feedList.Clear()
		// add items to feed list
		for _, feed := range allFeeds.Feeds {
			feedData := safeFeedData.GetEntries(feed.URL)
			feedList.AddItem(feedData.name, feed.URL, 0, func() {
				// handle user selecting item by moving focus to entry list
				app.SetFocus(entriesList)
			})
		}

		// handle user changing selected feed item by loading entries list
		feedList.SetChangedFunc(func(i int, feedName string, url string, shortcut rune) {
			loadEntriesIntoList(url)
		})

		// handle user changing selected item of entries list by loading entry text view
		entriesList.SetChangedFunc(func(i int, entryName string, secondaryText string, shortcut rune) {
			loadEntryTextView(i)
		})

		// when user hits escape in entries list, move focus back to feed list
		entriesList.SetDoneFunc(func() {
			app.SetFocus(feedList)
		})

		// load initial state of interface
		loadEntriesIntoList(getSelectedFeedUrl())
		//make sure there's at least one entry in selected
		if len(safeFeedData.GetEntries(getSelectedFeedUrl()).entries) > 0 {
			loadEntryTextView(0)
		}
	})
}

// Looks up the text of the corresponding entry and sets it on the text view
func loadEntryTextView(i int) {
	entryTextView.Clear()
	feedData := safeFeedData.GetEntries(getSelectedFeedUrl())
	entryTextView.SetText(feedData.entries[i].content)
}

// Urls of feeds are stored as secondary text on list items, uses that to look up selected feed
func getSelectedFeedUrl() string {
	_, url := feedList.GetItemText(feedList.GetCurrentItem())
	return url
}

// Send the entries for the selected feed into the entry list
func loadEntriesIntoList(url string) {
	entriesList.Clear()
	feedData := safeFeedData.GetEntries(url)
	for _, entry := range feedData.entries {
		entriesList.AddItem(entry.title, entry.url, 0, func() {
			// when an item in the entry list is selected, open the link in the browser
			_, url = entriesList.GetItemText(entriesList.GetCurrentItem())
			// if on windows escape &
			if runtime.GOOS == "windows" {
				strings.ReplaceAll(url, "&", "^&")
			}
			//use gox library to make platform specific call to open url in browser

			openBrowserModal := tview.NewModal()
			openBrowserModal.SetText("Open entry in browser?").
				AddButtons([]string{"Yes", "Cancel"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Yes" {
						osx.OpenDefault(url)
					}
					app.SetRoot(flex, false)
					app.SetFocus(entriesList)
				})
			app.SetRoot(openBrowserModal, false)

		})
	}
}

func main() {
	// init threadsafe feed data
	safeFeedData = &SafeFeedData{feedData: make(map[string]FeedDataModel)}

	// init ui elements
	app = tview.NewApplication()

	feedList = tview.NewList().ShowSecondaryText(false)
	feedList.SetBorder(true).SetTitle("Feeds")
	feedList.AddItem("Fetching Feed Data", "", 0, nil)

	entriesList = tview.NewList().ShowSecondaryText(false)
	entriesList.SetBorder(true).SetTitle("Entries")
	entriesList.AddItem("Fetching Feed Data", "", 0, nil)

	entryTextView = tview.NewTextView().
		SetChangedFunc(func() {
			app.Draw()
		})

	entryTextView.SetBorder(true)
	entryTextView.SetTitle("Description")
	entryTextView.SetText("Fetching Feed Data")

	// set up flex layout
	flex = tview.NewFlex().
		AddItem(feedList, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(entriesList, 0, 1, false).
			AddItem(entryTextView, 0, 2, false), 0, 2, false)

	// async call to load feed data
	go loadAllFeedDataAndUpdateInterface()

	// call to run tview app
	if uiError := app.SetRoot(flex, true).SetFocus(feedList).Run(); uiError != nil {
		panic(uiError)
	}

}