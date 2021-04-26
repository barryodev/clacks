package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"
	"github.com/rivo/tview"
	"golang.org/x/net/html"
)

// structs to decode json config
type AllFeeds struct {
	Feeds 	[]Feed	`json:"feeds"`
}

type Feed struct {
	Name 	string	`json:"name"`
	URL 	string 	`json:"url"`
	entries []Entry
}

type Entry struct {
	title string
	content string
}

// load feeds from json file
func loadJsonConfig() (AllFeeds, error){
	//open config file
	configFile, err := os.Open("feeds.json")
	if err != nil {
		return AllFeeds{}, errors.New("Error loading feeds.json file: " + err.Error())
	}

	//load data from file
	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		return AllFeeds{}, errors.New("Error Reading from feeds.json file: " + err.Error())
	}

	//unmarshall json
	var loadedFeeds AllFeeds
	err = json.Unmarshal(byteValue, &loadedFeeds)
	if err != nil {
		return AllFeeds{}, errors.New("Error unmarshalling json: " + err.Error())
	}

	return loadedFeeds, nil
}

// use gofeed library to load data from atom feed
func loadFeedData(url string) ([]Entry, error) {
	fp := gofeed.NewParser()
	fp.UserAgent = "Clacks - Terminal Atom Reader"
	feedData, err := fp.ParseURL(url)
	if err != nil {
		return []Entry{}, errors.New("Error processing atom feed: " + err.Error())
	}

	if len(feedData.Items) > 0 {
		entries := make([]Entry, len(feedData.Items))
		for i, item := range feedData.Items {
			entries[i] = Entry{
				title: html.UnescapeString(strip.StripTags(item.Title)),
				content: strings.TrimSpace(html.UnescapeString(strip.StripTags(item.Description))),
			}
		}
		return entries, nil
	}

	return []Entry{}, nil

}

// handle a feed item being selected


func main() {
	app := tview.NewApplication()

	allFeeds, configError := loadJsonConfig()
	if configError != nil {
		panic(configError)
	}

	for i := 0; i < len(allFeeds.Feeds); i++ {
		entries, atomFeedError := loadFeedData(allFeeds.Feeds[i].URL)
		if atomFeedError != nil {
			panic(atomFeedError)
		} else {
			allFeeds.Feeds[i].entries = entries
		}
	}

	feedList := tview.NewList().ShowSecondaryText(false)
	feedList.SetBorder(true).SetTitle("Feeds")

	entriesList := tview.NewList().ShowSecondaryText(false)
	entriesList.SetBorder(true).SetTitle("Entries")

	entryTextView := tview.NewTextView().
		SetChangedFunc(func() {
			app.Draw()
		})
	entryTextView.SetBorder(true)
	entryTextView.SetTitle("Description")

	if len(allFeeds.Feeds) > 0 && len(allFeeds.Feeds[0].entries) > 0 {
		entryTextView.SetText(allFeeds.Feeds[0].entries[0].content)
	}

	feedList.SetChangedFunc(func(i int, feedName string, secondaryText string, shortcut rune) {
		entriesList.Clear()
		selectedFeed := allFeeds.Feeds[i]
		for _, entry := range selectedFeed.entries {
			entriesList.AddItem(entry.title, "", 0, nil)
		}
	})

	for _, feed := range allFeeds.Feeds {
		feedList.AddItem(feed.Name, feed.URL, 0, func() {
			app.SetFocus(entriesList)
		})
	}

	entriesList.SetChangedFunc(func(i int, entryName string, secondaryText string, shortcut rune) {
		entryTextView.Clear()
		entryTextView.SetText(allFeeds.Feeds[feedList.GetCurrentItem()].entries[i].content)
	})

	entriesList.SetDoneFunc(func() {
		app.SetFocus(feedList)
	})

	flex := tview.NewFlex().
		AddItem(feedList, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(entriesList, 0, 1, false).
			AddItem(entryTextView, 0, 3, false), 0, 2, false)

	if uiError := app.SetRoot(flex, true).SetFocus(feedList).Run(); uiError != nil {
		panic(uiError)
	}

}