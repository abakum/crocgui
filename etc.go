//go:build !android

package main

import "fyne.io/fyne/v2"

func setupIntentHandler() {}

func quit(a fyne.App) {
	a.Quit()
}
