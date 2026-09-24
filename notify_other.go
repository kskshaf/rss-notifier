//go:build !linux && !windows

package main

import (
	"time"

	"github.com/gen2brain/beeep"
)

func sendArticleNotification(appName, title, body, iconPath, articleURL string, _ time.Duration, _ string, _ bool, _ bool) error {
	beeep.AppName = appName
	return beeep.Alert(title, body, iconPath)
}
