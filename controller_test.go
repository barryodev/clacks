package main

import (
	"github.com/mmcdole/gofeed"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewController(t *testing.T) {
	controller := NewController()

	assert.IsType(t, &tview.Application{}, controller.app)
	assert.IsType(t, &gofeed.Parser{}, controller.feedParser)
	assert.IsType(t, BrowserLauncher{}, controller.browserLauncher)

	castParser, ok := controller.feedParser.(*gofeed.Parser)
	assert.True(t, ok)
	assert.Equal(t, "Clacks - Terminal Atom/RSS Reader", castParser.UserAgent)
}

func TestStartUPLoopWithStubbedInterfaces(t *testing.T) {
	app := CreateStubbedApp(false)

	fakeFeed := CreateTestFeed()
	parser := createStubbedParser(&fakeFeed, false)

	browserLauncherStub := StubbedBrowserLauncher{withError: false}

	controllerWithStubs := Controller{app, parser, browserLauncherStub, nil, configFileName}

	controllerWithStubs.setupAndLaunchUILoop()

	assert.Eventually(t, func() bool { return controllerWithStubs.ui != nil }, time.Second, 10*time.Millisecond)

	assert.Eventually(t, func() bool { return len(controllerWithStubs.ui.data.configData.Feeds) == 4 }, time.Second, 10*time.Millisecond)

	for _, f := range controllerWithStubs.ui.app.(*StubbedApp).UpdateDraws {
		f()
	}

	assert.Equal(t, 4, controllerWithStubs.ui.feedList.GetItemCount())
}

func TestStartUPLoopWithConfigError(t *testing.T) {
	app := CreateStubbedApp(false)

	fakeFeed := CreateTestFeed()
	parser := createStubbedParser(&fakeFeed, true)

	browserLauncherStub := StubbedBrowserLauncher{withError: false}

	controllerWithStubs := Controller{app, parser, browserLauncherStub, nil, ""}

	assert.PanicsWithError(t, "error: could not find feeds.json config file", func() { controllerWithStubs.setupAndLaunchUILoop() })

}

func TestStartUPLoopWithUIError(t *testing.T) {
	app := CreateStubbedApp(true)

	fakeFeed := CreateTestFeed()
	parser := createStubbedParser(&fakeFeed, true)

	browserLauncherStub := StubbedBrowserLauncher{withError: false}

	controllerWithStubs := Controller{app, parser, browserLauncherStub, nil, configFileName}

	assert.PanicsWithError(t, "stubbed ui error", func() { controllerWithStubs.setupAndLaunchUILoop() })

}
