package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"runtime"
	"strings"
)

const feedPage = "feedsPage"
const helpPage = "helpPage"
const quitPage = "quitPage"
const refreshPage = "refreshPage"
const errorPage = "errorPage"
const openBrowserPage = "open"
const refreshMenuRegion = "refresh"
const helpMenuRegion = "help"
const quitMenuRegion = "quit"

// UI tview ui elements
type UI struct {
	app             TermApplication
	data            *Data
	browserLauncher BrowserLauncherInterface
	feedList        *tview.List
	entriesList     *tview.List
	entryTextView   *tview.TextView
	flex            *tview.Flex
	menuTextView    *tview.TextView
	previousFocus   tview.Primitive
	pages           *tview.Pages
}

// CreateUI create and configure all ui elements for app start up
func CreateUI(app TermApplication, data *Data) *UI {
	ui := &UI{}
	ui.app = app
	ui.data = data

	ui.feedList = tview.NewList().ShowSecondaryText(false)
	ui.feedList.SetBorder(true).SetTitle("Feeds")
	ui.feedList.SetBorderColor(tcell.ColorBlue)
	ui.feedList.AddItem("Fetching Feed Data", "", 0, nil)

	ui.entriesList = tview.NewList().ShowSecondaryText(false)
	ui.entriesList.SetBorder(true).SetTitle("Entries")
	ui.entriesList.AddItem("Fetching Feed Data", "", 0, nil)

	ui.entryTextView = tview.NewTextView().
		SetChangedFunc(func() {
			app.Draw()
		})

	ui.entryTextView.SetBorder(true)
	ui.entryTextView.SetWordWrap(true)
	ui.entryTextView.SetTitle("Description")
	ui.entryTextView.SetText("Fetching Feed Data")

	ui.menuTextView = tview.NewTextView()
	ui.menuTextView.SetRegions(true).SetDynamicColors(true).SetBorder(false)
	ui.menuTextView.SetText(`["` + refreshMenuRegion + `"][white:blue](r) Refresh [""][:black] ["` + helpMenuRegion + `"][white:blue](h) Help [""][:black] ["` + quitMenuRegion + `"][white:blue](q) Quit [""]`)

	ui.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(ui.feedList, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(ui.entriesList, 0, 1, false).
				AddItem(ui.entryTextView, 0, 2, false), 0, 2, false),
			0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(ui.menuTextView, len(ui.menuTextView.GetText(true)), 1, false).
			AddItem(tview.NewBox(), 0, 1, false),
			1, 0, false)

	ui.pages = tview.NewPages()
	ui.pages.AddPage(feedPage, ui.flex, true, true)

	ui.app.SetRoot(ui.pages, true)
	ui.app.SetFocus(ui.feedList)

	return ui
}

// load
func (ui *UI) loadAllFeedDataAndUpdateInterface() {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			ui.app.QueueUpdateDraw(func() {
				ui.createErrorPage(err.Error())
			})
		}
	}()

	ui.data.safeFeedData.Clear()
	err := ui.data.loadDataFromFeeds()
	if err != nil {
		panic(err)
	}
	ui.updateInterface()
}

// start ui run loop
func (ui *UI) startUILoop() error {
	err := ui.app.Run()
	if err != nil {
		return err
	}
	return nil
}

// Using QueueUpdateDraw as its a threadsafe way to update tview primitives
func (ui *UI) updateInterface() {
	ui.app.QueueUpdateDraw(ui.setupLists)
}

// load data into list and setup functions to handle user navigating list
func (ui *UI) setupLists() {
	ui.feedList.Clear()
	// add items to feed list
	for _, feed := range ui.data.configData.Feeds {
		feedData := ui.data.safeFeedData.GetEntries(feed.URL)
		ui.feedList.AddItem(feedData.name, feed.URL, 0, func() {
			// handle user selecting item by moving focus to entry list
			ui.switchAppFocus(ui.entriesList.Box, ui.feedList.Box, ui.entriesList)
		})
	}

	// handle user changing selected feed item by loading entries list
	ui.feedList.SetChangedFunc(func(i int, feedName string, url string, shortcut rune) {
		ui.loadEntriesIntoList(url)
	})

	// handle user changing selected item of entries list by loading entry text view
	ui.entriesList.SetChangedFunc(func(i int, entryName string, secondaryText string, shortcut rune) {
		ui.loadEntryTextView(i)
	})

	// when user hits escape in entries list, move focus back to feed list
	ui.entriesList.SetDoneFunc(func() {
		ui.switchAppFocus(ui.feedList.Box, ui.entriesList.Box, ui.feedList)
	})

	// load initial state of interface
	ui.loadEntriesIntoList(ui.getSelectedFeedURL())
	//make sure there's at least one entry in selected
	if len(ui.data.safeFeedData.GetEntries(ui.getSelectedFeedURL()).entries) > 0 {
		ui.loadEntryTextView(0)
	}

}

func (ui *UI) switchAppFocus(newBox *tview.Box, oldBox *tview.Box, newFocus tview.Primitive) {
	oldBox.SetBorderColor(tcell.ColorWhite)
	newBox.SetBorderColor(tcell.ColorBlue)
	ui.app.SetFocus(newFocus)
}

// Looks up the text of the corresponding entry and sets it on the text view
func (ui *UI) loadEntryTextView(i int) {
	ui.entryTextView.Clear()
	feedData := ui.data.safeFeedData.GetEntries(ui.getSelectedFeedURL())
	if feedData.entries != nil {
		ui.entryTextView.SetText(feedData.entries[i].content)
	}
}

// Urls of feeds are stored as secondary text on list items, uses that to look up selected feed
func (ui *UI) getSelectedFeedURL() string {
	_, url := ui.feedList.GetItemText(ui.feedList.GetCurrentItem())
	return url
}

// Send the entries for the selected feed into the entry list
func (ui *UI) loadEntriesIntoList(url string) {
	ui.entriesList.Clear()
	feedData := ui.data.safeFeedData.GetEntries(url)
	for _, entry := range feedData.entries {
		ui.entriesList.AddItem(entry.title, entry.url, 0, func() {
			// when an item in the entry list is selected, open the link in the browser
			_, url = ui.entriesList.GetItemText(ui.entriesList.GetCurrentItem())
			// if on windows escape &
			if runtime.GOOS == "windows" {
				strings.ReplaceAll(url, "&", "^&")
			}
			//use gox library to make platform specific call to open url in browser

			openBrowserModal := ui.createOverlayModal(openBrowserPage, "Open entry in browser?", []string{"Yes", "No"},
				func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Yes" {
						err := ui.browserLauncher.OpenDefault(url)
						if err != nil {
							panic(err)
						}
					}
					ui.pages.SwitchToPage(feedPage)
					ui.pages.RemovePage(openBrowserPage)
					ui.switchAppFocus(ui.entriesList.Box, ui.entryTextView.Box, ui.entriesList)
					ui.app.SetFocus(ui.entriesList)
				})
			ui.switchAppFocus(ui.entryTextView.Box, ui.entriesList.Box, openBrowserModal)
		})
	}
}

// handle user pressing menu shortcuts
func (ui *UI) setInputCaptureHandler() {
	ui.app.SetInputCapture(ui.handleKeyboardPressEvents)
}

func (ui *UI) handleKeyboardPressEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'q':
			ui.createQuitPage()
			ui.menuTextView.Highlight(quitMenuRegion)
		case 'h':
			ui.createHelpPage()
			ui.menuTextView.Highlight(helpMenuRegion)
		case 'r':
			ui.createRefreshPage()
			ui.menuTextView.Highlight(refreshMenuRegion)
		}
	}
	return event
}

// create modal box displaying error message after panic and quit app
func (ui *UI) createErrorPage(errorMessage string) *tview.Modal {
	errorBox := ui.createOverlayModal(errorPage, errorMessage, []string{"Okay"},
		func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Okay" {
				ui.app.Stop()
			}
		})
	ui.app.SetInputCapture(nil)
	ui.app.SetFocus(errorBox)

	return errorBox
}

// create the modal box describing application's functions, embed in a page and display
func (ui *UI) createHelpPage() {
	ui.previousFocus = ui.app.GetFocus()

	var stringBuilder strings.Builder
	_, _ = fmt.Fprint(&stringBuilder, "\nUse Arrow keys to navigate list items\n")
	_, _ = fmt.Fprint(&stringBuilder, "\nUse Enter and Esc to move between feed and entries lists\n")
	_, _ = fmt.Fprint(&stringBuilder, "\nHit Enter on an entry to open it in your default browser\n")
	_, _ = fmt.Fprint(&stringBuilder, "\nFeeds config are loaded from feeds.json\n")
	_, _ = fmt.Fprint(&stringBuilder, "\nCtrl-C or q to exit\n")

	helpBox := ui.createOverlayModal(helpPage, stringBuilder.String(), []string{"Done"},
		func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Done" {
				ui.pages.SwitchToPage(feedPage)
				ui.pages.RemovePage(helpPage)
				ui.app.SetFocus(ui.previousFocus)
				ui.menuTextView.Highlight()
			}
		})
	ui.app.SetFocus(helpBox)
}

// create modal box asking user if they want to quit application
func (ui *UI) createQuitPage() {
	ui.previousFocus = ui.app.GetFocus()

	quitBox := ui.createOverlayModal(quitPage, "Are you sure you want to quit?", []string{"Yes", "No"},
		func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				ui.app.Stop()
			} else if buttonLabel == "No" {
				ui.pages.SwitchToPage(feedPage)
				ui.pages.RemovePage(quitPage)
				ui.app.SetFocus(ui.previousFocus)
				ui.menuTextView.Highlight()
			}
		})

	ui.app.SetFocus(quitBox)
}

// create modal box asking if user wants to refresh page
func (ui *UI) createRefreshPage() {
	ui.previousFocus = ui.app.GetFocus()

	refreshBox := ui.createOverlayModal(refreshPage, "Do you want to refresh feed data?", []string{"Yes", "No"},
		func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				ui.setUITextToFetchingData()
				go ui.loadAllFeedDataAndUpdateInterface()
			}
			ui.pages.SwitchToPage(feedPage)
			ui.pages.RemovePage(refreshPage)
			ui.app.SetFocus(ui.previousFocus)
			ui.menuTextView.Highlight()
		})
	ui.app.SetFocus(refreshBox)
}

func (ui *UI) setUITextToFetchingData() {
	ui.feedList.Clear()
	ui.feedList.AddItem("Fetching Feed Data", "", 0, nil)
	ui.entriesList.Clear()
	ui.entriesList.AddItem("Fetching Feed Data", "", 0, nil)
	ui.entryTextView.Clear()
	ui.entryTextView.SetText("Fetching Feed Data")

}

func (ui *UI) createOverlayModal(pageName, modalText string, buttons []string, buttonPressedHandler func(buttonIndex int, buttonLabel string)) *tview.Modal {
	modalBox := tview.NewModal()
	modalBox.SetText(modalText)
	modalBox.AddButtons(buttons)
	modalBox.SetDoneFunc(buttonPressedHandler)
	ui.pages.AddPage(pageName, modalBox, true, true)
	return modalBox
}
