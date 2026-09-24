package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func dataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not find home directory:", err)
	}
	dir := filepath.Join(home, ".local", "share", "rss-notifier")
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal("Could not create data directory:", err)
	}
	return dir
}

func feedsFilePath() string { return filepath.Join(dataDir(), "feeds.json") }

type configFile struct {
	CheckInterval string `json:"check_interval"`
	Feeds         []Feed `json:"feeds"`
}

func loadConfig() (AppConfig, error) {
	path := feedsFilePath()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		initial := configFile{CheckInterval: defaultCheckInterval.String(), Feeds: []Feed{}}
		data, err = json.MarshalIndent(initial, "", "  ")
		if err != nil {
			return AppConfig{}, err
		}
		data = append(data, '\n')
		if err := os.WriteFile(path, data, 0644); err != nil {
			return AppConfig{}, fmt.Errorf("create config %s: %w", path, err)
		}
		log.Printf("Created feed configuration at %s", path)
	} else if err != nil {
		return AppConfig{}, fmt.Errorf("read config %s: %w", path, err)
	}

	interval := defaultCheckInterval
	var feeds []Feed
	trimmed := strings.TrimSpace(string(data))
	if strings.HasPrefix(trimmed, "[") {
		// Keep reading the previous bare-array format with the default interval.
		if err := json.Unmarshal(data, &feeds); err != nil {
			return AppConfig{}, fmt.Errorf("parse config %s: %w", path, err)
		}
	} else {
		var config configFile
		if err := json.Unmarshal(data, &config); err != nil {
			return AppConfig{}, fmt.Errorf("parse config %s: %w", path, err)
		}
		feeds = config.Feeds
		if config.CheckInterval != "" {
			interval, err = time.ParseDuration(config.CheckInterval)
			if err != nil || interval <= 0 {
				return AppConfig{}, fmt.Errorf("config %s: check_interval must be a positive duration such as \"45m\"", path)
			}
		}
	}

	for i, feed := range feeds {
		if feed.Name == "" || feed.URL == "" {
			return AppConfig{}, fmt.Errorf("config %s: feed %d must include name and url", path, i+1)
		}
	}
	return AppConfig{CheckInterval: interval, Feeds: feeds}, nil
}

// -----------------------------------------------------------------------------
// State: remembering which entries we've already seen
// We store seen entry IDs in a JSON file so we don't re-notify after a restart.
// -----------------------------------------------------------------------------

// seenFile is where we save seen entry IDs between runs.
// It lives in ~/.local/share/rss-notifier/seen.json
func seenFilePath() string {
	return filepath.Join(dataDir(), "seen.json")
}

func isInitialized() bool {
	_, err := os.Stat(filepath.Join(dataDir(), "initialized"))
	return err == nil
}

func markInitialized() {
	if err := os.WriteFile(filepath.Join(dataDir(), "initialized"), []byte("ok\n"), 0644); err != nil {
		log.Printf("Warning: could not save initialization state: %v", err)
	}
}

// loadSeen reads the set of already-seen entry IDs from disk.
// It returns a map where keys are IDs and values are true (a Go "set" pattern).
func loadSeen() map[string]bool {
	path := seenFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		// File doesn't exist yet — that's fine, return empty map
		return make(map[string]bool)
	}

	// json.Unmarshal parses JSON bytes into a Go value
	var seen map[string]bool
	if err := json.Unmarshal(data, &seen); err != nil {
		log.Println("Warning: could not parse seen file, starting fresh:", err)
		return make(map[string]bool)
	}
	return seen
}

// saveSeen writes the seen map back to disk.
func saveSeen(seen map[string]bool) {
	path := seenFilePath()

	// json.MarshalIndent produces human-readable JSON (nice if you want to inspect the file)
	data, err := json.MarshalIndent(seen, "", "  ")
	if err != nil {
		log.Println("Warning: could not serialize seen entries:", err)
		return
	}

	// os.WriteFile writes bytes to a file, creating or overwriting it
	if err := os.WriteFile(path, data, 0644); err != nil {
		log.Println("Warning: could not save seen file:", err)
	}
}

// checkIcon takes an icon string (local path) and always returns
// a local file path that beeep.Notify() can use.
// Returns an empty string if the icon is empty or cannot be resolved.
func checkIcon(icon string) string {
	// No icon specified — return empty, beeep will use the default system icon
	if icon == "" {
		return ""
	}

	// os.UserHomeDir() returns the current user's home directory
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not find home directory:", err)
	}
	dir := filepath.Join(home, ".local", "share", "rss-notifier", "icons")

	os.MkdirAll(dir, 0755)
	icon = filepath.Join(dir, icon)

	// Check the file actually exists (after ~ expansion, so the path is fully resolved)
	if _, err := os.Stat(icon); err != nil {
		// Log the fully expanded path so the user knows exactly what was checked
		log.Printf("Warning: icon file not found: %s\n", icon)
		return "" // fall back to default system icon
	}

	// log.Printf("Using local icon: %s\n", icon)
	return icon
}
