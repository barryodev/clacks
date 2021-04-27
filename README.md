# clacks
A Terminal Atom Reader

## Demo
![demo](img/clacksDemo.gif)

## Config
Json file called feeds.json in the same directory with the following format:
```json
{
  "feeds": [
    {
      "name": "The Register - BOFH",
      "url": "https://www.theregister.com/offbeat/bofh/headlines.atom"
    },
    {
      "name": "My Blog",
      "url": "https://barryodriscoll.net/feed/atom/"
    },
    {
      "name": "Boing Boing",
      "url": "https://boingboing.net/feed/atom"
    }
  ]
}
```

## Instructions
Add feeds name/url to feeds.json. Navigate lists using arrow keys. Hit enter/esc to select and deselect list items.

## Dependiences 
- [gofeed](https://github.com/mmcdole/gofeed) - feed parser
- [tview](https://github.com/rivo/tview) - terminal ui library
- [html-strip-tags-go](https://github.com/grokify/html-strip-tags-go) - strips html tags

## Credits
- The tview [postgres](https://github.com/rivo/tview/wiki/Postgres) example.
- Terry Pratchett for the name [clacks](https://discworld.fandom.com/wiki/Clacks).

## TODO
- Currently the feeds are loaded first and than the ui is loaded, loading the ui first and asynchronously fetching the feed data would be better.
- Add a menu that has help instructions, add new feed, remove feed, reload feeds, quit
- Add functionality to open browser with a selected entry 
