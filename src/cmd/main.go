package main

import (
	"logitech-mx-master-customization-linux/src"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.NewWithID("com.logitech.mxmaster.config")
	w := a.NewWindow("Logitech MX Master Configuration")

	app := src.NewApp(w)
	app.Run()
}
