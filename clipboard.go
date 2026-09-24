package main

import (
	"context"
	"sync"

	"golang.design/x/clipboard"
)

var (
	clipboardInitOnce sync.Once
	clipboardInitErr  error
)

func copyURLToClipboard(url string) error {
	clipboardInitOnce.Do(func() {
		clipboardInitErr = clipboard.Init()
	})
	if clipboardInitErr != nil {
		return clipboardInitErr
	}
	_, err := clipboard.Write(context.Background(), clipboard.FmtText, []byte(url))
	return err
}
