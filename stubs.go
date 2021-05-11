package main

import (
	"errors"
	"github.com/gdamore/tcell/v2"
	"github.com/mmcdole/gofeed"
	"github.com/rivo/tview"
	"sync"
)

// CreateTestAppWithSimScreen returns app with simulation screen for tests
func CreateTestAppWithSimScreen(width, height int) (*tview.Application, tcell.SimulationScreen) {
	screen := tcell.NewSimulationScreen("UTF-8")
	err := screen.Init()
	if err != nil {
		return nil, nil
	}
	screen.SetSize(width, height)

	app := tview.NewApplication()
	app.SetScreen(screen)

	return app, screen
}

// StubbedApp is tview.Application with mocked methods
type StubbedApp struct {
	FailRun     bool
	UpdateDraws []func()
	focus		tview.Primitive
	mutex       *sync.Mutex
}

// CreateStubbedApp returns app with simulation screen for tests
func CreateStubbedApp(failRun bool) TermApplication {
	app := &StubbedApp{
		FailRun:		failRun,
		UpdateDraws:	make([]func(), 0, 1),
		mutex:			&sync.Mutex{},
	}
	return app
}

// Run does nothing
func (app *StubbedApp) Run() error {
	if app.FailRun {
		return errors.New("fail")
	}
	return nil
}

// Stop does nothing
func (app *StubbedApp) Stop() {}

// SetRoot does nothing
func (app *StubbedApp) SetRoot(_ tview.Primitive, _ bool) *tview.Application {
	return nil
}

// Draw does nothing
func (app *StubbedApp) Draw() *tview.Application {
	return nil
}

// GetFocus does nothing
func (app *StubbedApp) GetFocus() tview.Primitive {
	return app.focus
}

// SetFocus does nothing
func (app *StubbedApp) SetFocus(p tview.Primitive) *tview.Application {
	app.focus = p
	return nil
}

// SetInputCapture does nothing
func (app *StubbedApp) SetInputCapture(_ func(event *tcell.EventKey) *tcell.EventKey) *tview.Application {
	return nil
}

// QueueUpdateDraw does nothing
func (app *StubbedApp) QueueUpdateDraw(f func()) *tview.Application {
	app.mutex.Lock()
	app.UpdateDraws = append(app.UpdateDraws, f)
	app.mutex.Unlock()
	return nil
}

// GetInputCapture does nothing
func (app *StubbedApp) GetInputCapture() func(event *tcell.EventKey) *tcell.EventKey {
	return nil
}

type StubbedParser struct {
	fakeFeed	*gofeed.Feed
	withError	bool
}

func createStubbedParser(fakeFeed *gofeed.Feed, withError bool) FeedParser {
	parser := &StubbedParser{
		fakeFeed: fakeFeed,
		withError: withError,
	}

	return parser
}

// ParseURL stub
func (parser *StubbedParser) ParseURL(_ string) (feed *gofeed.Feed, err error) {
	if parser.withError {
		return nil, errors.New("parser error")
	} else {
		return parser.fakeFeed, nil
	}
}

type StubbedBrowserLauncher struct {}

func (StubbedBrowserLauncher) OpenDefault(_ string) error {
	return nil
}




