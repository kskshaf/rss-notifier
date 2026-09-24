//go:build windows

package main

import (
	"encoding/xml"
	"log"
	"strings"
	"sync"
	"time"

	"git.sr.ht/~jackmordaunt/go-toast/wintoast"
	"github.com/gen2brain/beeep"
)

var (
	windowsToastCallbackOnce sync.Once
	windowsToastCopies       sync.Map
)

type windowsToast struct {
	XMLName        xml.Name            `xml:"toast"`
	ActivationType string              `xml:"activationType,attr"`
	Launch         string              `xml:"launch,attr"`
	Duration       string              `xml:"duration,attr"`
	Scenario       string              `xml:"scenario,attr,omitempty"`
	Visual         windowsToastVisual  `xml:"visual"`
	Audio          windowsToastAudio   `xml:"audio"`
	Actions        windowsToastActions `xml:"actions"`
}

type windowsToastVisual struct {
	Binding windowsToastBinding `xml:"binding"`
}

type windowsToastBinding struct {
	Template string             `xml:"template,attr"`
	Image    *windowsToastImage `xml:"image,omitempty"`
	Texts    []string           `xml:"text"`
}

type windowsToastImage struct {
	Placement string `xml:"placement,attr"`
	Source    string `xml:"src,attr"`
}

type windowsToastAudio struct {
	Silent bool `xml:"silent,attr"`
}

type windowsToastActions struct {
	Items []windowsToastAction `xml:"action"`
}

type windowsToastAction struct {
	ActivationType string `xml:"activationType,attr"`
	Content        string `xml:"content,attr"`
	Arguments      string `xml:"arguments,attr"`
}

type windowsToastCopy struct {
	url           string
	payload       string
	dismissOnCopy bool
}

func sendArticleNotification(appName, title, body, iconPath, articleURL string, timeout time.Duration, notificationUrgency string, openLinkOnClick, dismissOnCopy bool) error {
	windowsToastCallbackOnce.Do(func() {
		wintoast.SetActivationCallback(func(_ string, args string, _ []wintoast.UserData) {
			const copyPrefix = "copy:"
			if !strings.HasPrefix(args, copyPrefix) {
				return
			}
			value, ok := windowsToastCopies.Load(args)
			if !ok {
				return
			}
			copy := value.(windowsToastCopy)
			if err := copyURLToClipboard(copy.url); err != nil {
				log.Printf("Could not copy article URL: %v", err)
				return
			}
			if !copy.dismissOnCopy {
				if err := wintoast.Push(copy.payload, wintoast.PowershellFallback); err != nil {
					log.Printf("Could not keep notification open after copying: %v", err)
				}
			}
		})
	})

	duration := "short"
	scenario := ""
	if timeout == 0 {
		duration = "long"
		scenario = "reminder"
	} else if timeout > 7*time.Second {
		duration = "long"
	}

	binding := windowsToastBinding{Template: "ToastGeneric", Texts: []string{title, body}}
	if iconPath != "" {
		binding.Image = &windowsToastImage{Placement: "appLogoOverride", Source: iconPath}
	}
	activationType := "foreground"
	launch := ""
	if openLinkOnClick {
		activationType = "protocol"
		launch = articleURL
	}
	copyArgs := "copy:" + articleURL
	notification := windowsToast{
		ActivationType: activationType,
		Launch:         launch,
		Duration:       duration,
		Scenario:       scenario,
		Visual:         windowsToastVisual{Binding: binding},
		Audio:          windowsToastAudio{Silent: true},
		Actions: windowsToastActions{Items: []windowsToastAction{{
			ActivationType: "background",
			Content:        "Copy link",
			Arguments:      copyArgs,
		}}},
	}
	payload, err := xml.Marshal(notification)
	if err == nil {
		err = wintoast.SetAppData(wintoast.AppData{AppID: "RSS Notifier", IconPath: iconPath})
	}
	if err == nil {
		fullPayload := xml.Header + string(payload)
		windowsToastCopies.Store(copyArgs, windowsToastCopy{url: articleURL, payload: fullPayload, dismissOnCopy: dismissOnCopy})
		err = wintoast.Push(fullPayload, wintoast.PowershellFallback)
	}
	if err != nil {
		beeep.AppName = appName
		return beeep.Alert(title, body, iconPath)
	}
	return nil
}
