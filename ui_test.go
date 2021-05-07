package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testUrlOne = "theregistry.com"
const testUrlTwo = "google.com"

func TestAppRun(t *testing.T) {
	app := CreateStubbedApp(false)
	ui := CreateUI(app)

	err := ui.startUILoop()
	assert.Nil(t, err)
}

func TestAppFail(t *testing.T) {
	app := CreateStubbedApp(true)
	ui := CreateUI(app)

	err := ui.startUILoop()
	assert.NotNil(t, err)
	assert.Equal(t, "Fail", err.Error())
}

func TestCreateUIInitialUIState(t *testing.T) {
	app := CreateStubbedApp(true)
	ui := CreateUI(app)

	assert.NotNil(t, ui.feedList)
	assert.NotNil(t, ui.entriesList)
	assert.NotNil(t, ui.entryTextView)

	assert.Equal(t, 1, ui.feedList.GetItemCount())
	assert.Equal(t, 1, ui.entriesList.GetItemCount())

	feedText, _ := ui.feedList.GetItemText(0)
	assert.Equal(t, "Fetching Feed Data", feedText)

	entriesText, _ := ui.entriesList.GetItemText(0)
	assert.Equal(t, "Fetching Feed Data", entriesText)

	textViewText := ui.entryTextView.GetText(true)
	assert.Equal(t, "Fetching Feed Data", textViewText)
}

func TestCreateUISetupPages(t *testing.T) {
	app := CreateStubbedApp(true)
	ui := CreateUI(app)

	frontPage, primitive := ui.pages.GetFrontPage()
	assert.Equal(t, feedPage, frontPage)

	assert.Equal(t, ui.flex, primitive)
}

func TestFeedsListContainsFetchingDataOnStartup(t *testing.T){
	app, simScreen := CreateTestAppWithSimScreen(100, 100)
	defer simScreen.Fini()

	ui := CreateUI(app)
	ui.feedList.Box.SetRect(50,50,50,50)
	ui.feedList.Draw(simScreen)

	screenContentsAsString := getScreenContents(simScreen)

	assert.Contains(t, screenContentsAsString, "Fetching Feed Data")
}

func getScreenContents(simScreen tcell.SimulationScreen) string {
	simScreen.Show()

	simCells, _, _ := simScreen.GetContents()

	screenContents := make([]rune, len(simCells))
	for i := 0; i < len(simCells); i++ {
		screenContents[i] = simCells[i].Runes[0]
	}

	screenContentsAsString := string(screenContents)
	return screenContentsAsString
}

func TestFeedsListContainsFeedNameAfterLoadingData(t *testing.T) {
	app := CreateStubbedApp(true)
	ui := CreateUI(app)

	data := createTestData()

	ui.updateInterface(data)

	for _, f := range ui.app.(*StubbedApp).UpdateDraws {
		f()
	}

	listItemOne, _ := ui.feedList.GetItemText(0)
	assert.Equal(t, data.safeFeedData.GetEntries(testUrlOne).name, listItemOne)

	listItemTwo, _ := ui.feedList.GetItemText(1)
	assert.Equal(t, data.safeFeedData.GetEntries(testUrlTwo).name, listItemTwo)
}

func TestSwitchUIFocus(t *testing.T) {
	app, simScreen := CreateTestAppWithSimScreen(100, 100)
	defer simScreen.Fini()

	ui := CreateUI(app)

	assert.Equal(t, ui.feedList, ui.app.GetFocus())
	assert.Equal(t, tcell.ColorBlue, ui.feedList.GetBorderColor())
	assert.Equal(t, tcell.ColorWhite, ui.entriesList.GetBorderColor())

	ui.switchAppFocus(ui.entriesList.Box, ui.feedList.Box, ui.entriesList)

	assert.Equal(t, ui.entriesList, ui.app.GetFocus())
	assert.Equal(t, tcell.ColorBlue, ui.entriesList.GetBorderColor())
	assert.Equal(t, tcell.ColorWhite, ui.feedList.GetBorderColor())
}

func TestLoadEntryTextView(t *testing.T) {
	app := CreateStubbedApp(true)
	ui := CreateUI(app)

	data := createTestData()
	ui.updateInterface(data)

	for _, f := range ui.app.(*StubbedApp).UpdateDraws {
		f()
	}

	assert.Equal(t, data.safeFeedData.GetEntries(testUrlOne).entries[0].content,
		ui.entryTextView.GetText(true))

	ui.loadEntryTextView(data, 1)

	assert.Equal(t, data.safeFeedData.GetEntries(testUrlOne).entries[1].content,
		ui.entryTextView.GetText(true))

	ui.feedList.SetCurrentItem(1)
	ui.loadEntryTextView(data, 0)

	assert.Equal(t, data.safeFeedData.GetEntries(testUrlTwo).entries[0].content,
		ui.entryTextView.GetText(true))
}

func TestLoadEntriesIntoList(t *testing.T) {
	app := CreateStubbedApp(true)
	ui := CreateUI(app)

	data := createTestData()
	ui.loadEntriesIntoList(data, testUrlOne)

	for _, f := range ui.app.(*StubbedApp).UpdateDraws {
		f()
	}

	firstItemText, _ := ui.entriesList.GetItemText(ui.entriesList.GetCurrentItem())
	assert.Equal(t, data.safeFeedData.GetEntries(testUrlOne).entries[0].title, firstItemText)

	secondItemText, _ := ui.entriesList.GetItemText(1)
	assert.Equal(t, data.safeFeedData.GetEntries(testUrlOne).entries[1].title, secondItemText)
}

func TestHandleKeyboardPressQuitEvents(t *testing.T) {
	app, simScreen := CreateTestAppWithSimScreen(100, 100)
	defer simScreen.Fini()

	ui := CreateUI(app)

	currentFrontPage, _ := ui.pages.GetFrontPage()

	assert.Equal(t, feedPage, currentFrontPage)

	currentlyHighlighted := ui.menuTextView.GetHighlights()
	assert.Nil(t, currentlyHighlighted)

	keyEvent := tcell.NewEventKey(tcell.KeyRune, 'q', 0)
	ui.handleKeyboardPressEvents(keyEvent)

	quitFrontPage, quitModal := ui.pages.GetFrontPage()
	assert.Equal(t, quitPage, quitFrontPage)

	currentlyHighlighted = ui.menuTextView.GetHighlights()
	assert.Equal(t, quitMenuRegion, currentlyHighlighted[0])

	assert.NotNil(t, quitModal)
	ui.feedList.Box.SetRect(0,0,100,100)
	quitModal.Draw(simScreen)
	simScreen.Show()

	screenContentsAsString := getScreenContents(simScreen)

	assert.Contains(t, screenContentsAsString, "Are you sure you want to quit?")
}

func TestHandleKeyboardPressHelpEvents(t *testing.T) {
	app, simScreen := CreateTestAppWithSimScreen(150, 150)
	defer simScreen.Fini()

	ui := CreateUI(app)

	currentFrontPage, _ := ui.pages.GetFrontPage()

	assert.Equal(t, feedPage, currentFrontPage)

	currentlyHighlighted := ui.menuTextView.GetHighlights()
	assert.Nil(t, currentlyHighlighted)

	keyEvent := tcell.NewEventKey(tcell.KeyRune, 'h', 0)
	ui.handleKeyboardPressEvents(keyEvent)

	helpFrontPage, helpModal := ui.pages.GetFrontPage()
	assert.Equal(t, helpPage, helpFrontPage)

	currentlyHighlighted = ui.menuTextView.GetHighlights()
	assert.Equal(t, helpMenuRegion, currentlyHighlighted[0])

	assert.NotNil(t, helpModal)
	ui.feedList.Box.SetRect(0,0,150,150)
	helpModal.Draw(simScreen)
	simScreen.Show()

	screenContentsAsString := getScreenContents(simScreen)

	assert.Contains(t, screenContentsAsString, "Use Arrow keys to navigate list items")
}

func TestHandleKeyboardPressRefreshEvents(t *testing.T) {
	app, simScreen := CreateTestAppWithSimScreen(100, 100)
	defer simScreen.Fini()

	ui := CreateUI(app)

	currentFrontPage, _ := ui.pages.GetFrontPage()

	assert.Equal(t, feedPage, currentFrontPage)

	currentlyHighlighted := ui.menuTextView.GetHighlights()
	assert.Nil(t, currentlyHighlighted)

	keyEvent := tcell.NewEventKey(tcell.KeyRune, 'r', 0)
	ui.handleKeyboardPressEvents(keyEvent)

	refreshFrontPage, helpModal := ui.pages.GetFrontPage()
	assert.Equal(t, refreshPage, refreshFrontPage)

	currentlyHighlighted = ui.menuTextView.GetHighlights()
	assert.Equal(t, refreshMenuRegion, currentlyHighlighted[0])

	assert.NotNil(t, helpModal)
	ui.feedList.Box.SetRect(0,0,100,100)
	helpModal.Draw(simScreen)
	simScreen.Show()

	screenContentsAsString := getScreenContents(simScreen)

	assert.Contains(t, screenContentsAsString, "Do you want to refresh feed data?")
}

func createTestData() *Data {
	safeFeedData := &SafeFeedData{feedData: make(map[string]FeedDataModel)}
	safeFeedData.SetSiteData(testUrlOne, createFakeFeedDataModel("registry"))
	safeFeedData.SetSiteData(testUrlTwo, createFakeFeedDataModel("google"))

	allFeeds := &AllFeeds{[]Feed{Feed{testUrlOne}, {testUrlTwo}}}
	return &Data{safeFeedData: safeFeedData, allFeeds: allFeeds}
}

func createFakeFeedDataModel(name string) FeedDataModel {
	fakeEntryOne := Entry{
		title:   name + " fake title one",
		url:     testUrlOne + "/one",
		content: name + " fake content one",
	}

	fakeEntryTwo := Entry{
		title:   name + " fake title two",
		url:     testUrlOne + "/two",
		content: name + " fake content two",
	}

	fakeFeedDataModelOne := FeedDataModel{
		name:    name,
		entries: []Entry{fakeEntryOne, fakeEntryTwo},
	}
	return fakeFeedDataModelOne
}