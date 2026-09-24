//go:build linux

package main

import (
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/esiqveland/notify"
	"github.com/gen2brain/beeep"
	"github.com/godbus/dbus/v5"
)

var (
	linuxNotifierOnce     sync.Once
	linuxNotifier         notify.Notifier
	linuxNotifierErr      error
	linuxNotificationURLs sync.Map
)

func getLinuxNotifier() (notify.Notifier, error) {
	linuxNotifierOnce.Do(func() {
		conn, err := dbus.SessionBus()
		if err != nil {
			linuxNotifierErr = err
			return
		}
		linuxNotifier, linuxNotifierErr = notify.New(conn, notify.WithOnAction(func(action *notify.ActionInvokedSignal) {
			value, ok := linuxNotificationURLs.Load(action.ID)
			if !ok {
				return
			}
			data := value.(linuxNotificationData)
			switch action.ActionKey {
			case "default", "open":
				linuxNotificationURLs.Delete(action.ID)
				if !data.dismissOnCopy {
					if _, err := linuxNotifier.CloseNotification(action.ID); err != nil {
						log.Printf("Could not close notification: %v", err)
					}
				}
				if data.openLinkOnClick {
					if err := exec.Command("xdg-open", data.url).Start(); err != nil {
						log.Printf("Could not open article URL: %v", err)
					}
				}
			case "copy":
				if data.dismissOnCopy {
					linuxNotificationURLs.Delete(action.ID)
				}
				if err := copyURLToClipboard(data.url); err != nil {
					log.Printf("Could not copy article URL: %v", err)
				}
			}
		}), notify.WithOnClosed(func(closed *notify.NotificationClosedSignal) {
			linuxNotificationURLs.Delete(closed.ID)
		}))
	})
	return linuxNotifier, linuxNotifierErr
}

type linuxNotificationData struct {
	url             string
	openLinkOnClick bool
	dismissOnCopy   bool
}

func sendArticleNotification(appName, title, body, iconPath, articleURL string, timeout time.Duration, notificationUrgency string, openLinkOnClick, dismissOnCopy bool) error {
	notifier, err := getLinuxNotifier()
	if err != nil {
		beeep.AppName = appName
		return beeep.Alert(title, body, iconPath)
	}
	expireTimeout := timeout
	if notificationUrgency == "critical" {
		// Critical notifications should not expire automatically.
		expireTimeout = notify.ExpireTimeoutNever
	}
	n := notify.Notification{
		AppName:       appName,
		AppIcon:       iconPath,
		Summary:       title,
		Body:          body,
		Actions:       []notify.Action{{Key: "copy", Label: "Copy link"}},
		ExpireTimeout: expireTimeout,
	}
	switch notificationUrgency {
	case "low":
		n.SetUrgency(notify.UrgencyLow)
	case "critical":
		n.SetUrgency(notify.UrgencyCritical)
	default:
		n.SetUrgency(notify.UrgencyNormal)
	}
	if openLinkOnClick {
		n.Actions = append([]notify.Action{notify.NewDefaultAction("Open article")}, n.Actions...)
	}
	if !dismissOnCopy {
		n.AddHint(notify.Hint{ID: "resident", Variant: dbus.MakeVariant(true)})
	}
	id, err := notifier.SendNotification(n)
	if err != nil {
		beeep.AppName = appName
		return beeep.Alert(title, body, iconPath)
	}
	linuxNotificationURLs.Store(id, linuxNotificationData{url: articleURL, openLinkOnClick: openLinkOnClick, dismissOnCopy: dismissOnCopy})
	return nil
}
