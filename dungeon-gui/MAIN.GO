package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Dungeon Explorer 🏰")

	title := widget.NewLabel("Welcome to Dungeon Explorer!")
	sub := widget.NewLabel("A GUI-based dungeon adventure built in Go ⚔️")

	startBtn := widget.NewButton("Start Game", func() {
		startGame(w)
	})

	content := container.NewVBox(
		title,
		sub,
		startBtn,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()
}

func startGame(w fyne.Window) {
	w.SetContent(widget.NewLabel("The dungeon awaits..."))
}
