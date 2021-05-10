package main

import (
	"github.com/icza/gox/osx"
	"github.com/mmcdole/gofeed"
	"github.com/rivo/tview"
)

type BrowserLauncherInterface interface {
	OpenDefault(fileOrURL string) error
}

type BrowserLauncher struct {}
func (BrowserLauncher) OpenDefault(fileOrURL string) error {
	return osx.OpenDefault(fileOrURL)
}

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

	data.loadDataFromFeeds(fp)
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

	// Set Browser launcher
	ui.browserLauncher = BrowserLauncher{}

	ui.setInputCaptureHandler()

	// async call to load feed data
	go LoadAllFeedDataAndUpdateInterface(ui, data)

	err := ui.startUILoop()
	if err != nil {
		panic(err)
	}

}

