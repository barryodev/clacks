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

func main() {
	// init threadsafe feed data
	fp := gofeed.NewParser()
	fp.UserAgent = "Clacks - Terminal Atom Reader"
	data := NewData(fp)
	err := data.loadJsonConfig(configFileName)
	if err != nil {
		panic(err)
	}

	// init ui elements
	app := tview.NewApplication()
	ui := CreateUI(app, data)

	// Set Browser launcher
	ui.browserLauncher = BrowserLauncher{}

	ui.setInputCaptureHandler()

	// async call to load feed data
	go ui.loadAllFeedDataAndUpdateInterface()

	err = ui.startUILoop()
	if err != nil {
		panic(err)
	}

}

