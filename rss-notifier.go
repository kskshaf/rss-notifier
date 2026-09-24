// rss-notifier periodically checks RSS feeds and sends desktop notifications.
package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

const customTimeFormat = "2006/01/02 15:04:05"
const defaultCheckInterval = 45 * time.Minute
const defaultNotificationTimeout = 5 * time.Second
const defaultNotificationUrgency = "normal"

type AppConfig struct {
	CheckInterval       time.Duration
	NotificationTimeout time.Duration
	NotificationUrgency string
	OpenLinkOnClick     bool
	DismissOnCopy       bool
	Feeds               []Feed
}

// Feed is one RSS/Atom feed to watch. Icon is a filename in the icons directory.
type Feed struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Filter string `json:"filter,omitempty"`
	Icon   string `json:"icon,omitempty"`
}

func checkFeed(feed Feed, seen map[string]bool, notify bool, notificationTimeout time.Duration, notificationUrgency string, openLinkOnClick, dismissOnCopy bool) bool {
	parsedFeed, err := gofeed.NewParser().ParseURL(feed.URL)
	if err != nil {
		log.Printf("Error fetching feed '%s': %v", feed.Name, err)
		return false
	}

	iconPath := checkIcon(feed.Icon)
	for _, item := range parsedFeed.Items {
		id := item.GUID
		if id == "" {
			id = item.Link
		}
		if id == "" || seen[id] {
			continue
		}
		if feed.Filter != "" && !strings.Contains(strings.ToLower(item.Title), strings.ToLower(feed.Filter)) {
			seen[id] = true
			continue
		}
		if !notify {
			seen[id] = true
			continue
		}

		date := ""
		if item.PublishedParsed != nil {
			date = item.PublishedParsed.Format(customTimeFormat)
		}
		body := fmt.Sprintf("Time: %s\nLink: %s\n", date, item.Link)
		title := item.Title
		appName := feed.Name
		if feed.Name == "LTS Kernel" {
			appName = feed.Name + " " + feed.Filter
			title = strings.Replace(item.Title, ": longterm", "", 1)
		}
		if err := sendArticleNotification(appName, title, body, iconPath, item.Link, notificationTimeout, notificationUrgency, openLinkOnClick, dismissOnCopy); err != nil {
			log.Printf("Error sending notification for '%s': %v", item.Title, err)
		} else {
			fmt.Printf("Notified: [%s] %s\n", feed.Name, item.Title)
		}
		seen[id] = true
	}
	return true
}

func checkAllFeeds(feeds []Feed, notificationTimeout time.Duration, notificationUrgency string, openLinkOnClick, dismissOnCopy bool) {
	fmt.Println("Checking feeds at", time.Now().Format(customTimeFormat))
	seen := loadSeen()
	initialized := isInitialized()
	allFetched := true
	for _, feed := range feeds {
		if !checkFeed(feed, seen, initialized, notificationTimeout, notificationUrgency, openLinkOnClick, dismissOnCopy) {
			allFetched = false
		}
	}
	saveSeen(seen)
	if !initialized && allFetched {
		markInitialized()
		fmt.Println("Initial feed scan complete; existing entries were recorded without notifications.")
	}
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Could not load feed configuration:", err)
	}
	fmt.Println("rss-notifier started.")
	fmt.Printf("Checking every %s; notification timeout %s; open link on click %t; dismiss on copy %t; notification urgency %s\n\n", config.CheckInterval, formatNotificationTimeout(config.NotificationTimeout), config.OpenLinkOnClick, config.DismissOnCopy, config.NotificationUrgency)
	checkAllFeeds(config.Feeds, config.NotificationTimeout, config.NotificationUrgency, config.OpenLinkOnClick, config.DismissOnCopy)
	ticker := time.NewTicker(config.CheckInterval)
	defer ticker.Stop()
	for range ticker.C {
		checkAllFeeds(config.Feeds, config.NotificationTimeout, config.NotificationUrgency, config.OpenLinkOnClick, config.DismissOnCopy)
	}
}
