package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/icza/gox/osx"
	"github.com/mmcdole/gofeed"
	"github.com/rivo/tview"
)



type Controller struct {
	app 			TermApplication
	feedParser		FeedParser
	browserLauncher BrowserLauncherInterface
	ui				*UI
	configFileName	string
}

// TermApplication interface for the terminal UI app
type TermApplication interface {
	Run() error
	Draw() *tview.Application
	Stop()
	SetRoot(root tview.Primitive, fullscreen bool) *tview.Application
	SetFocus(p tview.Primitive) *tview.Application
	SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *tview.Application
	QueueUpdateDraw(f func()) *tview.Application
	GetFocus() tview.Primitive
	GetInputCapture() func(event *tcell.EventKey) *tcell.EventKey
}

// BrowserLauncherInterface interface to library for launching browser
type BrowserLauncherInterface interface {
	OpenDefault(fileOrURL string) error
}

type BrowserLauncher struct {}

func (BrowserLauncher) OpenDefault(fileOrURL string) error {
	return osx.OpenDefault(fileOrURL)
}

// FeedParser interface to gofeed library for parsing atom/rss feeds
type FeedParser interface {
	ParseURL(feedURL string) (feed *gofeed.Feed, err error)

}

func NewController() *Controller {
	feedParser := gofeed.NewParser()
	feedParser.UserAgent = "Clacks - Terminal Atom/RSS Reader"
	return &Controller{
		app: tview.NewApplication(),
		feedParser: feedParser,
		browserLauncher: BrowserLauncher{},
		configFileName: configFileName}
}

func (controller *Controller) setupAndLaunchUILoop() {
	// init threadsafe feed data
	data := NewData(controller.feedParser)
	err := data.loadJsonConfig(controller.configFileName)
	if err != nil {
		panic(err)
	}

	// init ui elements
	controller.ui = CreateUI(controller.app, data)

	// Set Browser launcher
	controller.ui.browserLauncher = controller.browserLauncher

	controller.ui.setInputCaptureHandler()

	// async call to load feed data
	go controller.ui.loadAllFeedDataAndUpdateInterface()

	err = controller.ui.startUILoop()
	if err != nil {
		panic(err)
	}

}
