package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateWaitWindow(a fyne.App) fyne.Window {
	win := a.NewWindow("Ожидайте")
	win.Resize(fyne.Size{
		Width:  200,
		Height: 100,
	})

	label := widget.NewLabel("Идет загрузка...")

	content := container.NewVBox(
		label)
	win.SetContent(content)
	win.SetFixedSize(true)
	win.Show()

	return win
}
