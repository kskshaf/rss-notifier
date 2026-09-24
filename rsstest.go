// rss-notifier periodically checks RSS feeds and sends desktop notifications.
package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/mmcdole/gofeed"
)

const customTimeFormat = "2006/01/02 15:04:05"
const checkInterval = 45 * time.Minute

// Feed is one RSS/Atom feed to watch. Icon is a filename in the icons directory.
type Feed struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Filter string `json:"filter,omitempty"`
	Icon   string `json:"icon,omitempty"`
}

func checkFeed(feed Feed, seen map[string]bool, notify bool) bool {
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
		body := fmt.Sprintf("更新时间：%s\n链接：%s\n", date, item.Link)
		title := item.Title
		beeep.AppName = feed.Name
		if feed.Name == "LTS Kernel" {
			beeep.AppName = feed.Name + " " + feed.Filter
			title = strings.Replace(item.Title, ": longterm", "", 1)
		}
		if err := beeep.Alert(title, body, iconPath); err != nil {
			log.Printf("Error sending notification for '%s': %v", item.Title, err)
		} else {
			fmt.Printf("Notified: [%s] %s\n", feed.Name, item.Title)
		}
		seen[id] = true
	}
	return true
}

func checkAllFeeds(feeds []Feed) {
	fmt.Println("Checking feeds at", time.Now().Format(customTimeFormat))
	seen := loadSeen()
	initialized := isInitialized()
	allFetched := true
	for _, feed := range feeds {
		if !checkFeed(feed, seen, initialized) {
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
	feeds, err := loadFeeds()
	if err != nil {
		log.Fatal("Could not load feed configuration:", err)
	}
	fmt.Println("rss-notifier started.")
	fmt.Printf("Checking every %s\n\n", checkInterval)
	checkAllFeeds(feeds)
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	for range ticker.C {
		checkAllFeeds(feeds)
	}
}
