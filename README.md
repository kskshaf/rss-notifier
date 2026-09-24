# RSS Notifier

**English** | [简体中文](README.zh-CN.md)

A lightweight RSS/Atom update notifier written in Go. It checks feeds immediately at startup and then at the configured interval. New entries that match a feed's filter trigger desktop notifications, and processed entries are stored locally to prevent duplicate notifications. Notifications include a link to the article.
On Windows and Linux, clicking a notification opens the article in your default browser; the “Copy link” action copies its URL to the clipboard.

This project was originally written with help from Claude Web.

## Features

- Parses RSS, Atom, and JSON Feed using `gofeed`
- Supports an optional, case-insensitive title filter for each feed
- Sends desktop notifications on Windows and Linux, with an action to copy the article URL
- Configures whether clicking a notification opens the article in the default browser
- Configures notification duration, including keeping notifications until dismissed where supported
- Configures whether a notification closes after copying its link
- On Linux, configures notification urgency (`low`, `normal`, or `critical`); the desktop environment determines the visual style
- Stores seen entries in `~/.local/share/rss-notifier/seen.json`
- Supports local notification icons in `~/.local/share/rss-notifier/icons/`

## Feed configuration

The configuration file is `~/.local/share/rss-notifier/feeds.json`. On first run, the program creates it with a 45-minute interval and an empty feed list. Edit the file to add feeds or change the interval:

```json
{
  "check_interval": "45m",
  "notification_timeout": "5s",
  "notification_urgency": "normal",
  "open_link_on_click": true,
  "dismiss_on_copy": true,
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

- Option Descriptions
   - `check_interval`: Uses Go duration syntax (e.g., `60s`, `30m`, or `1h`) and must be greater than zero.
   - `notification_timeout`: Accepts a non-negative Go duration (e.g., `5s`) or `"never"` to persist until dismissed (defaults to `5s`).
     - Linux applies the duration directly, except when `notification_urgency` is `critical`; critical notifications do not expire automatically, as required by the notification specification.
     - Windows only supports "short" and "long" presets, so finite values are mapped accordingly; `"never"` uses the Windows reminder scenario. Legacy bare-array feed configs remain supported (45-minute interval, default timeout). Omit `filter` or leave empty to match all titles. Omit `icon` to use the system default; place custom icons in `~/.local/share/rss-notifier/icons/`.
   - `notification_urgency`: Configurable on Linux as `low`, `normal`, or `critical`. Visual appearance depends on your desktop environment (e.g., critical notifications may use accent colors or persist until manually dismissed). Ignored on Windows.
   - `open_link_on_click`: Controls whether clicking the notification opens the article in your default browser; defaults to `true`. Set to `false` to disable this behavior while keeping the "Copy link" action.
   - `dismiss_on_copy`: Controls whether to dismiss the notification after copying the link; defaults to `true`. Set to `false` to keep the notification visible. Note that some Linux notification daemons may not support the `resident` hint. On Windows, toast button actions typically dismiss notifications automatically, so the app re-displays them after copying—setting this to `false` on Windows is therefore not recommended.


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

All runtime files are kept together in the application data directory:

```text
~/.local/share/rss-notifier/
├── feeds.json       # Feeds and notification settings
├── seen.json        # Articles already processed; prevents repeat notifications
├── initialized      # Marks completion of the first, silent feed scan
└── icons/           # Local notification icons referenced by feeds.json
```

On Windows, `~` means your user profile directory, for example `C:\Users\Your Name`. To use an icon, set its filename in a feed's `icon` field and put the image under `icons/`; for example, `icon: "example.png"` refers to `~/.local/share/rss-notifier/icons/example.png`.
