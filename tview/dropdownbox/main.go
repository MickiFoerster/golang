// Template code for building on
package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var dropdown = tview.NewDropDown()

func main() {
	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Name() == "Rune[A]":
			dropdown.SetOptions([]string{"Aachen", "Albertville"}, nil)
		case event.Name() == "Rune[B]":
			dropdown.SetOptions([]string{"Berlin", "Baden Baden"}, nil)
		}
		return event
	})
	dropdown.SetLabel("Select an option (hit Enter): ")
	if err := app.SetRoot(dropdown, true).Run(); err != nil {
		panic(err)
	}
}
