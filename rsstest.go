// rss-notifier: A lightweight RSS feed watcher that sends desktop notifications.
// Great for learning Go! Each section is commented to explain what's happening.
package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	// gofeed: parses RSS, Atom, and JSON feeds
	"github.com/mmcdole/gofeed"

	// beeep: sends cross-platform desktop notifications
	"github.com/gen2brain/beeep"
)

// -----------------------------------------------------------------------------
// Configuration
// Edit this section to add/remove feeds and change settings.
// -----------------------------------------------------------------------------

const custom_time_format string = "2006/01/02 15:04:05"

// How often to check feeds
const checkInterval = 45 * time.Minute

// Feed represents a single RSS/Atom feed to watch.
type Feed struct {
	Name   string // Friendly name shown in notifications
	URL    string // The RSS/Atom feed URL
	Filter string // Optional: only notify if title contains this string (case-insensitive). Leave "" for all.
	Icon   string // Optional: Icon for notify, local file only
}

// feeds is your list of subscribed feeds.
// Add or remove entries here to customize what you follow.
// Please put icons in ~/.local/share/rss-notifier/icons
var feeds = []Feed{
	{
		Name:   "LTS Kernel",
		URL:    "https://www.kernel.org/feeds/kdist.xml",
		Filter: "6.18", // Only notify for LTS 6.18 releases
		Icon:   "tux.png",
	},
	{
		Name:   "ArchLinux Latest News",
		URL:    "https://archlinux.org/feeds/news/",
		Filter: "",
		Icon:   "arch.png",
	},
}

// -----------------------------------------------------------------------------
// Core logic: fetching feeds and sending notifications
// -----------------------------------------------------------------------------

// checkFeed fetches a single feed and sends notifications for new entries.
// It takes a pointer to the seen map so it can update it in place.
func checkFeed(feed Feed, seen map[string]bool) {
	// gofeed.NewParser() creates a new feed parser
	fp := gofeed.NewParser()

	// fp.ParseURL fetches and parses the feed from the internet
	// In Go, functions often return (value, error) — always check the error!
	parsedFeed, err := fp.ParseURL(feed.URL)
	if err != nil {
		log.Printf("Error fetching feed '%s': %v\n", feed.Name, err)
		return // return early on error (like continue in other languages)
	}

	iconPath := checkIcon(feed.Icon)

	// Loop over each item in the feed
	// In Go, range gives you (index, value) for slices
	for _, item := range parsedFeed.Items {
		// Build a unique ID for this entry.
		// We prefer the feed's own GUID, but fall back to the link URL.
		id := item.GUID
		if id == "" {
			id = item.Link
		}

		// Skip if we've already seen this entry
		if seen[id] {
			continue
		}

		// Apply filter: skip if title doesn't contain the filter string
		// strings.Contains is case-sensitive, so we use ToLower on both sides
		if feed.Filter != "" && !strings.Contains(strings.ToLower(item.Title), strings.ToLower(feed.Filter)) {
			// Mark as seen so we don't recheck it every time
			seen[id] = true
			continue
		}

		// Format it nicely for your notification body
		date := ""
		if item.PublishedParsed != nil {
			date = item.PublishedParsed.Format(custom_time_format)
		}

		// Build the notification message
		title := fmt.Sprintf("%s", item.Title)
		// body := item.Link + item.UpdatedParsed// Show the URL as the notification body
		body := fmt.Sprintf("更新时间：%s\n", date)

		// beeep.Notify sends a desktop notification via libnotify (notify-send)
		// Arguments: title, message, icon path (empty = default)

		beeep.AppName = feed.Name
		if feed.Name == "LTS Kernel" {
			beeep.AppName = feed.Name + " " + feed.Filter
			title = strings.Replace(item.Title, ": longterm", "", 1)
		}

		if err := beeep.Alert(title, body, iconPath); err != nil {
			log.Printf("Error sending notification for '%s': %v\n", item.Title, err)
		} else {
			fmt.Printf("Notified: [%s] %s\n", feed.Name, item.Title)
		}

		// Mark this entry as seen
		seen[id] = true
	}
}

// checkAllFeeds iterates over all configured feeds.
func checkAllFeeds() {
	fmt.Println("Checking feeds at", time.Now().Format(custom_time_format))

	// Load seen entries from disk
	seen := loadSeen()

	// Check each feed
	for _, feed := range feeds {
		checkFeed(feed, seen)
	}

	// Save updated seen entries back to disk
	saveSeen(seen)
}

// -----------------------------------------------------------------------------
// Entry point
// -----------------------------------------------------------------------------

func main() {
	fmt.Println("rss-notifier started.")
	fmt.Printf("Checking every %s\n\n", checkInterval)

	// Check immediately on startup
	checkAllFeeds()

	// time.NewTicker creates a channel that sends a value every `checkInterval`
	// This is Go's idiomatic way to do periodic tasks
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop() // defer runs this when main() exits (cleanup)

	// range over a channel blocks until a value arrives, then loops
	for range ticker.C {
		checkAllFeeds()
	}
}
