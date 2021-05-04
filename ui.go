package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/icza/gox/osx"
	"github.com/rivo/tview"
	"os"
	"runtime"
	"strings"
)

const feedPage = "feedsPage"
const helpPage = "helpPage"
const quitPage = "quitPage"
const refreshPage = "refreshPage"
const openBrowserPage = "open"
const refreshMenuRegion = "refresh"
const helpMenuRegion = "help"
const quitMenuRegion = "quit"

// Using QueueUpdateDraw as its a threadsafe way to update tview primitives
func LoadFeedDataIntoLists() {
	app.QueueUpdateDraw(func() {
		feedList.Clear()
		// add items to feed list
		for _, feed := range allFeeds.Feeds {
			feedData := safeFeedData.GetEntries(feed.URL)
			feedList.AddItem(feedData.name, feed.URL, 0, func() {
				// handle user selecting item by moving focus to entry list
				SwitchAppFocus(entriesList.Box, feedList.Box, entriesList)
			})
		}

		// handle user changing selected feed item by loading entries list
		feedList.SetChangedFunc(func(i int, feedName string, url string, shortcut rune) {
			LoadEntriesIntoList(url)
		})

		// handle user changing selected item of entries list by loading entry text view
		entriesList.SetChangedFunc(func(i int, entryName string, secondaryText string, shortcut rune) {
			LoadEntryTextView(i)
		})

		// when user hits escape in entries list, move focus back to feed list
		entriesList.SetDoneFunc(func() {
			SwitchAppFocus(feedList.Box, entriesList.Box, feedList)
		})

		// load initial state of interface
		LoadEntriesIntoList(GetSelectedFeedUrl())
		//make sure there's at least one entry in selected
		if len(safeFeedData.GetEntries(GetSelectedFeedUrl()).entries) > 0 {
			LoadEntryTextView(0)
		}
	})
}

func SwitchAppFocus(newBox *tview.Box, oldBox *tview.Box, newFocus tview.Primitive) {
	oldBox.SetBorderColor(tcell.ColorWhite)
	newBox.SetBorderColor(tcell.ColorBlue)
	app.SetFocus(newFocus)
}

// Looks up the text of the corresponding entry and sets it on the text view
func LoadEntryTextView(i int) {
	entryTextView.Clear()
	feedData := safeFeedData.GetEntries(GetSelectedFeedUrl())
	if feedData.entries != nil {
		entryTextView.SetText(feedData.entries[i].content)
	}
}

// Urls of feeds are stored as secondary text on list items, uses that to look up selected feed
func GetSelectedFeedUrl() string {
	_, url := feedList.GetItemText(feedList.GetCurrentItem())
	return url
}

// Send the entries for the selected feed into the entry list
func LoadEntriesIntoList(url string) {
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

			openBrowserModal := CreateOverlayModal(openBrowserPage, "Open entry in browser?", []string{"Yes", "No"},
				func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Yes" {
						err := osx.OpenDefault(url)
						if err != nil {
							panic(err)
						}
					}
					pages.SwitchToPage(feedPage)
					pages.RemovePage(openBrowserPage)
					SwitchAppFocus(entriesList.Box, entryTextView.Box, entriesList)
					app.SetFocus(entriesList)
				})
			SwitchAppFocus(entryTextView.Box, entriesList.Box, openBrowserModal)
		})
	}
}

// create text view of the menu
func InitMenu() *tview.TextView {
	menu := tview.NewTextView()
	menu.SetRegions(true).SetDynamicColors(true).SetBorder(false)
	menu.SetText(`["` + refreshMenuRegion + `"][white:blue](r) Refresh [""][:black] ["` + helpMenuRegion + `"][white:blue](h) Help [""][:black] ["` + quitMenuRegion + `"][white:blue](q) Quit [""]`)
	return menu
}

// handle user pressing menu shortcuts
func HandleMenuKeyPresses() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				CreateQuitPage()
				menuTextView.Highlight(quitMenuRegion)
			case 'h':
				CreateHelpPage()
				menuTextView.Highlight(helpMenuRegion)
			case 'r':
				CreateRefreshPage()
				menuTextView.Highlight(refreshMenuRegion)
			}
		}
		return event
	})
}

// create modal box displaying error message after panic and quit app
func CreateErrorPage(errorMessage string) {
	errorBox := CreateOverlayModal(helpPage, errorMessage, []string{"Okay"},
		func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Okay" {
				app.Stop()
			}
		})
	app.SetInputCapture(nil)
	app.SetFocus(errorBox)
}

// create the modal box describing application's functions, embed in a page and display
func CreateHelpPage() {
	previousFocus = app.GetFocus()

	var stringBuilder strings.Builder
	_, _ = fmt.Fprint(&stringBuilder, "\nUse Arrow keys to navigate list items\n")
	_, _ = fmt.Fprint(&stringBuilder, "\nUse Enter and Esc to move between feed and entries lists\n")
	_, _ = fmt.Fprint(&stringBuilder, "\nHit Enter on an entry to open it in your default browser\n")
	_, _ = fmt.Fprint(&stringBuilder, "\nFeeds config are loaded from feeds.json\n")
	_, _ = fmt.Fprint(&stringBuilder, "\nCtrl-C or q to exit\n")

	helpBox := CreateOverlayModal(helpPage, stringBuilder.String(), []string{"Done"},
		func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Done" {
				pages.SwitchToPage(feedPage)
				pages.RemovePage(helpPage)
				app.SetFocus(previousFocus)
				menuTextView.Highlight()
			}
		})
	app.SetFocus(helpBox)
}

// create modal box asking user if they want to quit application
func CreateQuitPage() {
	previousFocus = app.GetFocus()

	quitBox := CreateOverlayModal(quitPage, "Are you sure you want to quit?", []string{"Yes", "No"},
		func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				app.Stop()
				os.Exit(0)
			} else if buttonLabel == "No" {
				pages.SwitchToPage(feedPage)
				pages.RemovePage(quitPage)
				app.SetFocus(previousFocus)
				menuTextView.Highlight()
			}
		})

	app.SetFocus(quitBox)
}

// create modal box asking if user wants to refresh page
func CreateRefreshPage() {
	previousFocus = app.GetFocus()

	refreshBox := CreateOverlayModal(refreshPage, "Do you want to refresh feed data?", []string{"Yes", "No"},
		func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				SetUITextToFetchingData()
				go LoadAllFeedDataAndUpdateInterface()
			}
			pages.SwitchToPage(feedPage)
			pages.RemovePage(refreshPage)
			app.SetFocus(previousFocus)
			menuTextView.Highlight()
		})
	app.SetFocus(refreshBox)
}

func SetUITextToFetchingData() {
	feedList.Clear()
	feedList.AddItem("Fetching Feed Data", "", 0, nil)
	entriesList.Clear()
	entriesList.AddItem("Fetching Feed Data", "", 0, nil)
	entryTextView.Clear()
	entryTextView.SetText("Fetching Feed Data")

}

func CreateOverlayModal(pageName, modalText string, buttons []string, buttonPressedHandler func(buttonIndex int, buttonLabel string)) *tview.Modal {
	modalBox := tview.NewModal()
	modalBox.SetText(modalText)
	modalBox.AddButtons(buttons)
	modalBox.SetDoneFunc(buttonPressedHandler)
	pages.AddPage(pageName, modalBox, true, true)
	return modalBox
}

// setup the main application layout and embed in a page
func CreateFeedLayout() *tview.Flex {
	// set up flex layout
	flexLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(feedList, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(entriesList, 0, 1, false).
				AddItem(entryTextView, 0, 2, false), 0, 2, false),
			0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(menuTextView, len(menuTextView.GetText(true)), 1, false).
			AddItem(tview.NewBox(), 0, 1, false),
			1, 0, false)

	return flexLayout

}