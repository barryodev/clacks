package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// tview ui elements, config and feed data
var (
	app 			*tview.Application
	feedList		*tview.List
	entriesList 	*tview.List
	entryTextView 	*tview.TextView
	flex			*tview.Flex
	menuTextView	*tview.TextView
	previousFocus	tview.Primitive
	pages			*tview.Pages
	safeFeedData	*SafeFeedData
	allFeeds		*AllFeeds
)

func main() {
	// init threadsafe feed data
	safeFeedData = &SafeFeedData{feedData: make(map[string]FeedDataModel)}

	// init ui elements
	app = tview.NewApplication()
	pages = tview.NewPages()

	feedList = tview.NewList().ShowSecondaryText(false)
	feedList.SetBorder(true).SetTitle("Feeds")
	feedList.SetBorderColor(tcell.ColorBlue)
	feedList.AddItem("Fetching Feed Data", "", 0, nil)

	entriesList = tview.NewList().ShowSecondaryText(false)
	entriesList.SetBorder(true).SetTitle("Entries")
	entriesList.AddItem("Fetching Feed Data", "", 0, nil)

	entryTextView = tview.NewTextView().
		SetChangedFunc(func() {
			app.Draw()
		})

	entryTextView.SetBorder(true)
	entryTextView.SetWordWrap(true)
	entryTextView.SetTitle("Description")
	entryTextView.SetText("Fetching Feed Data")

	menuTextView = InitMenu()

	flex = CreateFeedLayout()
	pages.AddPage(feedPage, flex, true, true)

	HandleMenuKeyPresses()

	// async call to load feed data
	go LoadAllFeedDataAndUpdateInterface()

	// call to run tview app
	if uiError := app.SetRoot(pages, true).SetFocus(feedList).Run(); uiError != nil {
		panic(uiError)
	}

}

