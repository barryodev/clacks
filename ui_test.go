package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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