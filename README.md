# clacks
[![Build Status](https://api.travis-ci.com/barryodev/clacks.svg?branch=main)](https://travis-ci.com/github/barryodev/clacks)
[![codecov](https://codecov.io/gh/barryodev/clacks/branch/main/graph/badge.svg)](https://app.codecov.io/gh/barryodev/clacks)
[![Go Report Card](https://goreportcard.com/badge/github.com/barryodev/clacks)](https://goreportcard.com/report/github.com/barryodev/clacks)

A Terminal Based Atom/RSS Reader that reads a list of feeds from a json file and displays them, once selected a feed entry can be opened in the browser.

## Demo
![demo](img/clacksDemo.gif)

## Config
Json file called feeds.json in the same directory with the following format:
```json
{
  "feeds": [
    {
      "url": "https://www.theregister.com/offbeat/bofh/headlines.atom"
    },
    {
      "url": "https://barryodriscoll.net/feed/atom/"
    },
    {
      "url": "https://boingboing.net/feed/atom"
    },
    {
      "url": "https://lifehacker.com/rss"
    }
  ]
}
```

## Instructions
- Add feeds name/url to feeds.json. 
- Navigate lists using arrow keys. 
- Hit enter/esc to select and deselect list items.
- Hit enter on an entry to open in default system browser.
- Use menu shortcuts to perform related tasks.
- Ctrl-C to quit.

## Dependiences 
- [gofeed](https://github.com/mmcdole/gofeed) - feed parser
- [tview](https://github.com/rivo/tview) - terminal ui library
- [html-strip-tags-go](https://github.com/grokify/html-strip-tags-go) - strips html tags
- [gox](https://github.com/icza/gox) - utility library used to open browser cross platform manner

## Credits
- The tview [postgres](https://github.com/rivo/tview/wiki/Postgres) example.
- The tview [presentation](https://github.com/rivo/tview/tree/master/demos/presentation) demo.  
- This stackoverflow [answer](https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang) on how to open a browser in go and escape urls on windows by [icza](https://stackoverflow.com/users/1705598/icza).
- The project [gdu](https://github.com/dundee/gdu) provided a great example on how to unit test code that calls the tview library with stubs.
- Terry Pratchett for the name [clacks](https://discworld.fandom.com/wiki/Clacks).

## Alternatives
- [GORSS](https://github.com/lallassu/gorss) is a terminal atom/rss reader also written in go that is built directly on tcell, the library tview is built on.