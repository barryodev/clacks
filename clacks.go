package main

import (
	"github.com/mmcdole/gofeed"
	"github.com/rivo/tview"
)

func LoadAllFeedDataAndUpdateInterface(ui *UI, data *Data){
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			ui.app.QueueUpdateDraw(func() {
				ui.createErrorPage(err.Error())
			})
		}
	}()

	fp := gofeed.NewParser()
	fp.UserAgent = "Clacks - Terminal Atom Reader"

	data.loadAllFeedData(fp)
	ui.updateInterface(data)
}

func main() {
	// init threadsafe feed data
	safeFeedData := &SafeFeedData{feedData: make(map[string]FeedDataModel)}
	allFeeds := &AllFeeds{}

	data := &Data{safeFeedData: safeFeedData, allFeeds: allFeeds}

	// init ui elements
	app := tview.NewApplication()

	ui := CreateUI(app)

	ui.handleMenuKeyPresses()

	// async call to load feed data
	go LoadAllFeedDataAndUpdateInterface(ui, data)

	err := ui.startUILoop()
	if err != nil {
		panic(err)
	}

}

