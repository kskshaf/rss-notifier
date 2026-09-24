# RSS Notifier

**English** | [简体中文](README.zh-CN.md)

A lightweight RSS/Atom update notifier written in Go. It checks feeds immediately at startup and then at the configured interval. New entries that match a feed's filter trigger desktop notifications, and processed entries are stored locally to prevent duplicate notifications. Notifications include a link to the article.

This project was originally written with help from Claude Web.

## Features

- Parses RSS, Atom, and JSON Feed using `gofeed`.
- Supports an optional, case-insensitive title filter for each feed.
- Sends desktop notifications with `beeep`.
- Stores seen entries in `~/.local/share/rss-notifier/seen.json`.
- Supports local notification icons in `~/.local/share/rss-notifier/icons/`.

## Feed configuration

The configuration file is `~/.local/share/rss-notifier/feeds.json`. On first run, the program creates it with a 45-minute interval and an empty feed list. Edit the file to add feeds or change the interval:

```json
{
  "check_interval": "45m",
  "feeds": [
    {
      "name": "Example Feed",
      "url": "https://example.com/feed.xml",
      "filter": "release",
      "icon": "example.png"
    }
  ]
}
```

`check_interval` uses Go duration syntax, such as `30m` or `1h`, and must be greater than zero. The previous bare-array feed configuration is still supported, using a 45-minute interval. Omit `filter` or set it to an empty string to match every title. Omit `icon` to use the system default. Place icon files in `~/.local/share/rss-notifier/icons/`.

## Run

Install Go (the required version is specified in `go.mod`) and make sure a desktop notification service is available. From the project directory, run:

```sh
go run .
```

To build and run a binary:

```sh
go build -o rss-notifier .
./rss-notifier
```

On the first run, the program quietly records existing feed entries as its baseline. It notifies you about new entries after initialization. It then checks feeds at the configured `check_interval`. Press `Ctrl+C` to stop it.

## Data and icons

- Feed configuration: `~/.local/share/rss-notifier/feeds.json`
- Seen entries: `~/.local/share/rss-notifier/seen.json`
- Icons: `~/.local/share/rss-notifier/icons/`
