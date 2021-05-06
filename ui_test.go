package main

import (
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

	simScreen.Show()

	simCells, _, _ := simScreen.GetContents()

	screenContents := make([]rune, len(simCells))
	for i := 0; i < len(simCells); i++ {
		screenContents[i] = simCells[i].Runes[0]
	}

	screenContentsAsString := string(screenContents)

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
	assert.Equal(t, "registry", listItemOne)

	listItemTwo, _ := ui.feedList.GetItemText(1)
	assert.Equal(t, "google", listItemTwo)
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