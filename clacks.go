package main

import (
	"github.com/rivo/tview"
)

// tview ui elements, config and feed data
var (
	safeFeedData	*SafeFeedData
	allFeeds		*AllFeeds
)

func main() {
	// init threadsafe feed data
	safeFeedData = &SafeFeedData{feedData: make(map[string]FeedDataModel)}

	// init ui elements
	app := tview.NewApplication()

	ui := CreateUI(app)

	ui.handleMenuKeyPresses()

	// async call to load feed data
	go LoadAllFeedDataAndUpdateInterface(ui)

	err := ui.startUILoop()
	if err != nil {
		panic(err)
	}

}

