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

// --- Game Structures ---

type Room struct {
	Type        string
	HasEnemy    bool
	HasPotion   bool
	HasTreasure bool
	Visited     bool
}

type Player struct {
	X, Y    int
	HP      int
	Gold    int
	Visited map[string]bool
}

const gridSize = 5

var dungeon [gridSize][gridSize]Room
var gameLog, status, roomLabel *widget.Label

// --- Helper Functions ---

func key(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

// Generate random dungeon
func initDungeon() {
	rand.Seed(time.Now().UnixNano())
	roomTypes := []string{"Cave", "Library", "Armory", "Shrine", "Treasury"}

	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			roomType := roomTypes[rand.Intn(len(roomTypes))]
			dungeon[y][x] = Room{
				Type:        roomType,
				HasEnemy:    rand.Float64() < 0.35,
				HasPotion:   rand.Float64() < 0.15,
				HasTreasure: rand.Float64() < 0.25,
				Visited:     false,
			}
		}
	}
}

// Map visuals
func drawDungeon(player *Player) *fyne.Container {
	grid := container.NewGridWithColumns(gridSize)
	roomSize := fyne.NewSize(70, 70)

	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			rKey := key(x, y)

			var roomColor color.NRGBA
			if x == player.X && y == player.Y {
				roomColor = color.NRGBA{0, 120, 0, 255} // Current: green
			} else if player.Visited[rKey] {
				roomColor = color.NRGBA{90, 70, 50, 255} // Visited
			} else {
				roomColor = color.NRGBA{25, 25, 25, 255} // Unexplored
			}

			bg := canvas.NewRectangle(roomColor)
			bg.SetMinSize(roomSize)
			bg.StrokeWidth = 2
			bg.StrokeColor = color.NRGBA{60, 60, 60, 255}

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

// Exploration mechanics
func exploreRoom(player *Player) string {
	room := &dungeon[player.Y][player.X]
	room.Visited = true
	player.Visited[key(player.X, player.Y)] = true

	roomLabel.SetText(fmt.Sprintf("📍 You entered a %s", room.Type))

	switch room.Type {
	case "Cave":
		roomLabel.SetText("🌑 You step into a dark, echoing cave.")
	case "Library":
		roomLabel.SetText("📚 Ancient books whisper secrets in this library.")
	case "Armory":
		roomLabel.SetText("🗡️ Rusty weapons line the walls of the armory.")
	case "Shrine":
		roomLabel.SetText("🙏 The shrine glows faintly — a place of power.")
	case "Treasury":
		roomLabel.SetText("💎 You see golden glimmers all around — a treasury!")
	}

	if room.HasEnemy {
		damage := rand.Intn(15) + 10
		player.HP -= damage
		room.HasEnemy = false
		if player.HP <= 0 {
			return fmt.Sprintf("💀 You were slain by a monster after taking %d damage!", damage)
		}
		return fmt.Sprintf("⚔️ You fought bravely and took %d damage.", damage)
	}

	if room.HasPotion {
		heal := rand.Intn(15) + 10
		player.HP += heal
		room.HasPotion = false
		return fmt.Sprintf("🧪 You found a potion and restored %d HP!", heal)
	}

	if room.HasTreasure {
		gold := rand.Intn(25) + 15
		player.Gold += gold
		room.HasTreasure = false
		return fmt.Sprintf("💰 You found %d gold in a chest!", gold)
	}

	return "The room is empty but eerie..."
}

// --- Main Function ---
func main() {
	a := app.New()
	w := a.NewWindow("Dungeon Explorer 🏰")
	w.Resize(fyne.NewSize(620, 680))

	initDungeon()
	player := &Player{X: 0, Y: 0, HP: 100, Gold: 0, Visited: make(map[string]bool)}
	player.Visited[key(player.X, player.Y)] = true

	gameLog = widget.NewLabel("Welcome, adventurer! Use arrows to explore.")
	status = widget.NewLabel(fmt.Sprintf("❤️ HP: %d | 💰 Gold: %d", player.HP, player.Gold))
	roomLabel = widget.NewLabel("📍 You start at the dungeon entrance.")

	grid := drawDungeon(player)

	// Movement logic
	movePlayer := func(dx, dy int) {
		newX, newY := player.X+dx, player.Y+dy
		if newX < 0 || newX >= gridSize || newY < 0 || newY >= gridSize {
			gameLog.SetText("🚫 You hit a stone wall. Try another direction.")
			return
		}
		player.X, player.Y = newX, newY
		message := exploreRoom(player)
		gameLog.SetText(message)
		status.SetText(fmt.Sprintf("❤️ HP: %d | 💰 Gold: %d", player.HP, player.Gold))

		// Refresh only grid
		grid.Objects = drawDungeon(player).Objects
		grid.Refresh()
	}

	// Movement buttons
	btnUp := widget.NewButton("↑", func() { movePlayer(0, -1) })
	btnDown := widget.NewButton("↓", func() { movePlayer(0, 1) })
	btnLeft := widget.NewButton("←", func() { movePlayer(-1, 0) })
	btnRight := widget.NewButton("→", func() { movePlayer(1, 0) })

	movement := container.NewVBox(
		container.NewCenter(btnUp),
		container.NewCenter(container.NewHBox(btnLeft, btnDown, btnRight)),
	)

	legend := widget.NewLabel("Legend: 🧙‍♂️ You | 🟩 Current | 🟫 Visited | ⬛ Unexplored")

	content := container.NewVBox(
		widget.NewLabelWithStyle("Dungeon Explorer 🏰", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		roomLabel,
		grid,
		movement,
		legend,
		gameLog,
		status,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
