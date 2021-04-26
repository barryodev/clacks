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

## Dependiences 
- [gofeed](https://github.com/mmcdole/gofeed) - feed parser
- [tview](https://github.com/rivo/tview) - terminal ui library
- [html-strip-tags-go](https://github.com/grokify/html-strip-tags-go) - strips html tags

## Credits
- The tview [postgres](https://github.com/rivo/tview/wiki/Postgres) example.
- Terry Pratchett for the name [clacks](https://discworld.fandom.com/wiki/Clacks).
