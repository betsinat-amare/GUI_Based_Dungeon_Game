package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// --- Game Data Structures ---

type Room struct {
	HasEnemy   bool
	HasPotion  bool
	HasTreasure bool
	Visited    bool
}

type Player struct {
	X, Y        int
	HP, Gold    int
	Visited     map[string]bool
}

const gridSize = 5

var dungeon [gridSize][gridSize]Room
var gameLog *widget.Label

// --- Helper Functions ---

func key(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func initDungeon() {
	rand.Seed(time.Now().UnixNano())

	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			// Randomized dungeon content
			hasEnemy := rand.Float64() < 0.25     // 25% chance
			hasPotion := rand.Float64() < 0.15    // 15% chance
			hasTreasure := rand.Float64() < 0.20  // 20% chance

			dungeon[y][x] = Room{
				HasEnemy:   hasEnemy,
				HasPotion:  hasPotion,
				HasTreasure: hasTreasure,
				Visited:    false,
			}
		}
	}
}

func drawDungeon(player *Player) *fyne.Container {
	grid := container.NewGridWithColumns(gridSize)
	roomSize := fyne.NewSize(70, 70)

	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			rKey := key(x, y)
			room := dungeon[y][x]

			// Base room color
			roomColor := color.NRGBA{30, 30, 30, 255} // Unexplored: dark
			if player.Visited[rKey] {
				roomColor = color.NRGBA{100, 80, 60, 255} // Visited: stone brown
			}
			if x == player.X && y == player.Y {
				roomColor = color.NRGBA{40, 120, 40, 255} // Current: green
			}

			bg := canvas.NewRectangle(roomColor)
			bg.SetMinSize(roomSize)

			// Room symbol
			symbol := ""
			if x == player.X && y == player.Y {
				symbol = "🧙‍♂️"
			}

			label := widget.NewLabelWithStyle(symbol, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
			roomBox := container.NewMax(bg, container.NewCenter(label))

			grid.Add(roomBox)
		}
	}

	return grid
}

func exploreRoom(player *Player) string {
	room := &dungeon[player.Y][player.X]
	room.Visited = true
	player.Visited[key(player.X, player.Y)] = true

	if room.HasEnemy {
		damage := rand.Intn(20) + 5
		player.HP -= damage
		room.HasEnemy = false
		if player.HP <= 0 {
			return fmt.Sprintf("💀 You fought bravely but took %d damage and died!", damage)
		}
		return fmt.Sprintf("⚔️ You fought an enemy and took %d damage!", damage)
	}

	if room.HasPotion {
		heal := rand.Intn(15) + 10
		player.HP += heal
		room.HasPotion = false
		return fmt.Sprintf("🧪 You found a potion and restored %d HP!", heal)
	}

	if room.HasTreasure {
		gold := rand.Intn(30) + 10
		player.Gold += gold
		room.HasTreasure = false
		return fmt.Sprintf("💰 You found %d gold pieces!", gold)
	}

	return "The room is empty..."
}

// --- GUI Game Logic ---

func main() {
	a := app.New()
	w := a.NewWindow("Dungeon Explorer 🏰")
	w.Resize(fyne.NewSize(600, 600))

	initDungeon()

	player := &Player{X: 0, Y: 0, HP: 100, Gold: 0, Visited: make(map[string]bool)}
	player.Visited[key(player.X, player.Y)] = true

	gameLog = widget.NewLabel("Welcome to the Dungeon! Explore carefully...")

	grid := drawDungeon(player)
	status := widget.NewLabel(fmt.Sprintf("HP: %d | Gold: %d", player.HP, player.Gold))

	movePlayer := func(dx, dy int) {
		newX, newY := player.X+dx, player.Y+dy
		if newX < 0 || newX >= gridSize || newY < 0 || newY >= gridSize {
			gameLog.SetText("🚫 You hit a wall! Try another direction.")
			return
		}

		player.X, player.Y = newX, newY
		message := exploreRoom(player)
		gameLog.SetText(message)

		status.SetText(fmt.Sprintf("HP: %d | Gold: %d", player.HP, player.Gold))
		w.SetContent(buildUI(w, player, grid, status))
	}

	// Movement buttons
	btnUp := widget.NewButton("↑", func() { movePlayer(0, -1) })
	btnDown := widget.NewButton("↓", func() { movePlayer(0, 1) })
	btnLeft := widget.NewButton("←", func() { movePlayer(-1, 0) })
	btnRight := widget.NewButton("→", func() { movePlayer(1, 0) })

	movement := container.NewGridWithColumns(3,
		btnLeft,
		btnDown,
		btnRight,
	)

	content := container.NewVBox(
		widget.NewLabelWithStyle("Dungeon Explorer 🏰", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		grid,
		container.NewCenter(btnUp),
		container.NewCenter(movement),
		widget.NewLabel("Legend: 🧙‍♂️ You | 🟩 Current | 🟫 Visited | ⬛ Unexplored"),
		gameLog,
		status,
	)

	w.SetContent(content)
	w.ShowAndRun()
}

func buildUI(w fyne.Window, player *Player, oldGrid *fyne.Container, status *widget.Label) *fyne.Container {
	grid := drawDungeon(player)
	return container.NewVBox(
		widget.NewLabelWithStyle("Dungeon Explorer 🏰", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		grid,
		widget.NewLabel("Legend: 🧙‍♂️ You | 🟩 Current | 🟫 Visited | ⬛ Unexplored"),
		gameLog,
		status,
	)
}
