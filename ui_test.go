package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
	assert.Equal(t, "fail", err.Error())
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
	data := createTestData()
	simScreen, ui := setupWithSimScreen(data)
	defer simScreen.Fini()

	assert.Equal(t, ui.feedList, ui.app.GetFocus())
	assert.Equal(t, tcell.ColorBlue, ui.feedList.GetBorderColor())
	assert.Equal(t, tcell.ColorWhite, ui.entriesList.GetBorderColor())

	ui.switchAppFocus(ui.entriesList.Box, ui.feedList.Box, ui.entriesList)

	assert.Equal(t, ui.entriesList, ui.app.GetFocus())
	assert.Equal(t, tcell.ColorBlue, ui.entriesList.GetBorderColor())
	assert.Equal(t, tcell.ColorWhite, ui.feedList.GetBorderColor())
}

func TestLoadEntryTextView(t *testing.T) {
	data := createTestData()
	app := CreateStubbedApp(true)
	ui := CreateUI(app)

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
	data := createTestData()
	app := CreateStubbedApp(true)
	ui := CreateUI(app)

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
	data := createTestData()
	simScreen, ui := setupWithSimScreen(data)

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

	contents, _, _ := simScreen.GetContents()
	assert.NotNil(t, contents)

	keyEvent = tcell.NewEventKey(tcell.KeyEnter, rune(0), 0)
	quitModal.InputHandler()(keyEvent, nil)

	contents, _, _ = simScreen.GetContents()
	assert.Nil(t, contents)
}

func TestQuitMenuCancelButtonPress(t *testing.T) {
	data := createTestData()
	simScreen, ui := setupWithSimScreen(data)
	defer simScreen.Fini()

	keyEvent := tcell.NewEventKey(tcell.KeyRune, 'q', 0)
	ui.handleKeyboardPressEvents(keyEvent)

	quitFrontPage, quitModal := ui.pages.GetFrontPage()
	assert.Equal(t, quitPage, quitFrontPage)

	castModal, ok := quitModal.(*tview.Modal)
	assert.True(t, ok)

	keyEvent = tcell.NewEventKey(tcell.KeyTab, rune(0), 0)
	castModal.InputHandler()(keyEvent, nil)

	keyEvent = tcell.NewEventKey(tcell.KeyEnter, rune(0), 0)
	castModal.InputHandler()(keyEvent, nil)

	quitFrontPage, quitModal = ui.pages.GetFrontPage()
	assert.Equal(t, feedPage, quitFrontPage)
}

func TestHandleKeyboardPressHelpEvents(t *testing.T) {
	data := createTestData()
	simScreen, ui := setupWithSimScreen(data)
	defer simScreen.Fini()

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

	keyEvent = tcell.NewEventKey(tcell.KeyEnter, rune(0), 0)
	helpModal.InputHandler()(keyEvent, nil)

	currentFrontPage, _ = ui.pages.GetFrontPage()
	assert.Equal(t, feedPage, currentFrontPage)
	assert.Nil(t, ui.menuTextView.GetHighlights())
	assert.Equal(t, ui.feedList, ui.app.GetFocus())
}

func TestHandleKeyboardPressRefreshEvents(t *testing.T) {
	data := createTestData()
	simScreen, ui := setupWithSimScreen(data)
	defer simScreen.Fini()

	currentFrontPage, _ := ui.pages.GetFrontPage()
	assert.Equal(t, feedPage, currentFrontPage)

	currentlyHighlighted := ui.menuTextView.GetHighlights()
	assert.Nil(t, currentlyHighlighted)

	keyEvent := tcell.NewEventKey(tcell.KeyRune, 'r', 0)
	ui.handleKeyboardPressEvents(keyEvent)

	refreshFrontPage, refreshModal := ui.pages.GetFrontPage()
	assert.Equal(t, refreshPage, refreshFrontPage)

	currentlyHighlighted = ui.menuTextView.GetHighlights()
	assert.Equal(t, refreshMenuRegion, currentlyHighlighted[0])

	assert.NotNil(t, refreshModal)
	refreshModal.Draw(simScreen)
	simScreen.Show()

	screenContentsAsString := getScreenContents(simScreen)
	assert.Contains(t, screenContentsAsString, "Do you want to refresh feed data?")
}

func TestUserNavigatesLists(t *testing.T) {
	data := createTestData()
	simScreen, ui := setupWithSimScreen(data)
	defer simScreen.Fini()

	feedName, _ := ui.feedList.GetItemText(0)
	assert.Equal(t, data.safeFeedData.GetEntries(testUrlOne).name, feedName)

	ui.feedList.SetCurrentItem(1)
	entryName, _ := ui.entriesList.GetItemText(ui.entriesList.GetCurrentItem())
	assert.Equal(t, data.safeFeedData.GetEntries(testUrlTwo).entries[0].title, entryName)

	ui.entriesList.SetCurrentItem(1)
	assert.Equal(t, data.safeFeedData.GetEntries(testUrlTwo).entries[1].content, ui.entryTextView.GetText(true))
}

func TestUserSelectingItemInFeedListAndLeavingEntriesList(t *testing.T){
	data := createTestData()
	simScreen, ui := setupWithSimScreen(data)
	defer simScreen.Fini()

	assert.Equal(t, ui.feedList, ui.app.GetFocus())

	keyEvent := tcell.NewEventKey(tcell.KeyEnter, rune(0), 0)
	ui.feedList.InputHandler()(keyEvent, nil)

	assert.Equal(t, ui.entriesList, ui.app.GetFocus())

	keyEvent = tcell.NewEventKey(tcell.KeyESC, rune(0), 0)
	ui.entriesList.InputHandler()(keyEvent, nil)

	assert.Equal(t, ui.feedList, ui.app.GetFocus())
}

func TestUserSelectingItemInEntriesList(t *testing.T) {
	data := createTestData()
	simScreen, ui := setupWithSimScreen(data)
	defer simScreen.Fini()

	keyEvent := tcell.NewEventKey(tcell.KeyEnter, rune(0), 0)
	ui.entriesList.InputHandler()(keyEvent, nil)

	ui.pages.Draw(simScreen)
	simScreen.Show()

	screenContentsAsString := getScreenContents(simScreen)
	assert.Contains(t, screenContentsAsString, "Open entry in browser?")

	button, ok := ui.app.GetFocus().(*tview.Button)
	assert.True(t, ok)

	ui.browserLauncher = StubbedBrowserLauncher{}

	assert.Equal(t, "Yes", button.GetLabel())
	keyEvent = tcell.NewEventKey(tcell.KeyEnter, rune(0), 0)
	button.InputHandler()(keyEvent, nil)

	assert.Equal(t, ui.entriesList, ui.app.GetFocus())
}

func TestCreateErrorPage(t *testing.T) {
	data := createTestData()
	simScreen, ui := setupWithSimScreen(data)
	defer simScreen.Fini()

	errorModal := ui.createErrorPage("test error")

	errorModal.Draw(simScreen)
	simScreen.Show()

	screenContentsAsString := getScreenContents(simScreen)
	assert.Contains(t, screenContentsAsString, "test error")
}

func setupWithSimScreen(data * Data) (tcell.SimulationScreen, *UI) {
	app, simScreen := CreateTestAppWithSimScreen(150, 150)
	ui := CreateUI(app)
	if data != nil {
		ui.loadInitialDataAndListNavigationFunctions(data)()
	}
	return simScreen, ui
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

func createTestData() *Data {
	safeFeedData := &SafeFeedData{feedData: make(map[string]FeedDataModel)}
	safeFeedData.SetSiteData(testUrlOne, createFakeFeedDataModel("registry", testUrlOne))
	safeFeedData.SetSiteData(testUrlTwo, createFakeFeedDataModel("google", testUrlTwo))

	allFeeds := &AllFeeds{[]Feed{{testUrlOne}, {testUrlTwo}}}
	return &Data{safeFeedData: safeFeedData, allFeeds: allFeeds}
}

func createFakeFeedDataModel(name, url string) FeedDataModel {
	fakeEntryOne := Entry{
		title:   name + " fake title one",
		url:     url + "/one",
		content: name + " fake content one",
	}

	fakeEntryTwo := Entry{
		title:   name + " fake title two",
		url:     url + "/two",
		content: name + " fake content two",
	}

	fakeFeedDataModelOne := FeedDataModel{
		name:    name,
		entries: []Entry{fakeEntryOne, fakeEntryTwo},
	}
	return fakeFeedDataModelOne
}