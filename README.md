# clacks
A Terminal Atom Reader - reads a list of atom feeds from a json file and asynchronously fetches them.

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
    }
  ]
}
```

## Instructions
- Add feeds name/url to feeds.json. 
- Navigate lists using arrow keys. 
- Hit enter/esc to select and deselect list items.
- Ctrl-C to quit.

## Dependiences 
- [gofeed](https://github.com/mmcdole/gofeed) - feed parser
- [tview](https://github.com/rivo/tview) - terminal ui library
- [html-strip-tags-go](https://github.com/grokify/html-strip-tags-go) - strips html tags

## Credits
- The tview [postgres](https://github.com/rivo/tview/wiki/Postgres) example.
- Terry Pratchett for the name [clacks](https://discworld.fandom.com/wiki/Clacks).

## TODO
- Add a menu that has help instructions, add new feed, remove feed, reload feeds, quit
- Add functionality to open browser with a selected entry 
- Handle RSS feeds, shouldn't be much extra work as [gofeed](https://github.com/mmcdole/gofeed) supports detecting feed type and handling both.
