package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.NewWithID("com.logitech.mxmaster.config")
	w := a.NewWindow("Logitech MX Master Configuration")

	app := NewApp(w)
	app.Run()
}
