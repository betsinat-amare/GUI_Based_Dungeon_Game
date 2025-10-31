package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Game state
type Player struct {
	X, Y  int
	HP    int
	MaxHP int
}

const gridSize = 3 // 3x3 dungeon grid

func main() {
	a := app.New()
	w := a.NewWindow("Dungeon Explorer 🏰")

	title := widget.NewLabelWithStyle("Dungeon Explorer 🏰", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	sub := widget.NewLabel("A GUI Dungeon Adventure in Go ⚔️")

	startBtn := widget.NewButton("Start Game", func() {
		startGame(w)
	})

	menu := container.NewVBox(
		title,
		sub,
		startBtn,
	)
	w.SetContent(menu)
	w.Resize(fyne.NewSize(500, 400))
	w.ShowAndRun()
}

// Core game UI
func startGame(w fyne.Window) {
	player := &Player{X: 1, Y: 1, HP: 100, MaxHP: 100}

	infoLabel := widget.NewLabel("")
	logLabel := widget.NewLabel("You enter the dungeon...")

	// Draw dungeon grid (3x3)
	grid := drawDungeon(player)

	// Movement buttons
	btnUp := widget.NewButton("⬆️", func() {
		movePlayer(player, 0, -1, infoLabel, logLabel, grid)
	})
	btnDown := widget.NewButton("⬇️", func() {
		movePlayer(player, 0, 1, infoLabel, logLabel, grid)
	})
	btnLeft := widget.NewButton("⬅️", func() {
		movePlayer(player, -1, 0, infoLabel, logLabel, grid)
	})
	btnRight := widget.NewButton("➡️", func() {
		movePlayer(player, 1, 0, infoLabel, logLabel, grid)
	})

	movement := container.NewGridWithColumns(3,
		btnLeft,
		btnDown,
		btnRight,
	)

	statusBox := container.NewVBox(
		infoLabel,
		canvas.NewRectangle(color.RGBA{200, 200, 200, 255}),
		logLabel,
	)

	updateInfo(player, infoLabel)

	content := container.NewVBox(
		widget.NewLabelWithStyle("Dungeon Explorer 🏰", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		grid,
		btnUp,
		movement,
		statusBox,
	)

	w.SetContent(content)
}

// Draw grid of rectangles to represent rooms
func drawDungeon(player *Player) *fyne.Container {
	grid := container.NewGridWithColumns(gridSize)

	for y := 0; y < gridSize; y++ {
		row := container.NewGridWithColumns(gridSize)
		for x := 0; x < gridSize; x++ {
			rect := canvas.NewRectangle(color.NRGBA{100, 100, 100, 255})
			rect.SetMinSize(fyne.NewSize(60, 60))
			if x == player.X && y == player.Y {
				rect.FillColor = color.NRGBA{0, 255, 0, 255} // Player
			}
			row.Add(rect)
		}
		grid.Add(row)
	}
	return grid
}

// Move player and update grid/log
func movePlayer(p *Player, dx, dy int, infoLabel, logLabel *widget.Label, grid *fyne.Container) {
	newX, newY := p.X+dx, p.Y+dy

	if newX < 0 || newX >= gridSize || newY < 0 || newY >= gridSize {
		logLabel.SetText("You bump into a wall! 🧱")
		return
	}

	p.X, p.Y = newX, newY
	logLabel.SetText(fmt.Sprintf("You moved to room (%d, %d).", p.X, p.Y))
	updateInfo(p, infoLabel)

	// Redraw grid
	newGrid := drawDungeon(p)
	grid.Objects = newGrid.Objects
	grid.Refresh()
}

// Update info label (HP, position)
func updateInfo(p *Player, label *widget.Label) {
	label.SetText(fmt.Sprintf("HP: %d/%d  |  Position: (%d,%d)", p.HP, p.MaxHP, p.X, p.Y))
}
