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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Player struct {
	X, Y     int
	HP       int
	MaxHP    int
	Gold     int
	Potions  int
	Visited  map[string]bool
}

type Enemy struct {
	Name   string
	HP     int
	Attack int
}

const gridSize = 4

func main() {
	rand.Seed(time.Now().UnixNano())

	a := app.NewWithID("dungeon.gui")
	a.Settings().SetTheme(theme.DarkTheme())

	w := a.NewWindow("Dungeon Explorer 🏰")

	title := widget.NewLabelWithStyle("Dungeon Explorer 🏰", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	sub := widget.NewLabel("A GUI Dungeon Adventure in Go ⚔️")

	startBtn := widget.NewButton("Start Game", func() {
		startGame(w)
	})
	quitBtn := widget.NewButton("Quit", func() {
		w.Close()
	})

	menu := container.NewVBox(
		title,
		sub,
		startBtn,
		quitBtn,
	)

	w.SetContent(container.NewCenter(menu))
	w.Resize(fyne.NewSize(600, 500))
	w.ShowAndRun()
}

func startGame(w fyne.Window) {
	player := &Player{
		X: 1, Y: 1,
		HP: 100, MaxHP: 100,
		Gold: 0, Potions: 1,
		Visited: map[string]bool{"1,1": true},
	}

	infoLabel := widget.NewLabel("")
	logLabel := widget.NewLabel("You enter the dungeon...")

	grid := drawDungeon(player)

	btnUp := widget.NewButton("⬆️", func() {
		movePlayer(player, 0, -1, infoLabel, logLabel, grid, w)
	})
	btnDown := widget.NewButton("⬇️", func() {
		movePlayer(player, 0, 1, infoLabel, logLabel, grid, w)
	})
	btnLeft := widget.NewButton("⬅️", func() {
		movePlayer(player, -1, 0, infoLabel, logLabel, grid, w)
	})
	btnRight := widget.NewButton("➡️", func() {
		movePlayer(player, 1, 0, infoLabel, logLabel, grid, w)
	})
	btnPotion := widget.NewButton("🧪 Use Potion", func() {
		usePotion(player, logLabel, infoLabel)
	})

	// Menu bar
	menu := container.NewHBox(
		widget.NewButton("🔁 Restart", func() { startGame(w) }),
		widget.NewButton("🚪 Quit", func() { w.Close() }),
	)

	movement := container.NewGridWithColumns(3,
		btnLeft,
		btnDown,
		btnRight,
	)

	statusBox := container.NewVBox(
		drawHPBar(player),
		infoLabel,
		btnPotion,
		logLabel,
	)

	updateInfo(player, infoLabel)

	content := container.NewVBox(
		menu,
		widget.NewLabelWithStyle("Dungeon Explorer 🏰", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		grid,
		container.NewCenter(btnUp),
		container.NewCenter(movement),
		statusBox,
	)

	w.SetContent(container.NewPadded(content))
}

func drawDungeon(player *Player) *fyne.Container {
	grid := container.NewGridWithColumns(gridSize)

	for y := 0; y < gridSize; y++ {
		row := container.NewGridWithColumns(gridSize)
		for x := 0; x < gridSize; x++ {
			key := fmt.Sprintf("%d,%d", x, y)
			rect := canvas.NewRectangle(color.NRGBA{50, 50, 50, 255})
			rect.SetMinSize(fyne.NewSize(60, 60))

			if x == player.X && y == player.Y {
				rect.FillColor = color.NRGBA{0, 255, 0, 255} // Player
			} else if player.Visited[key] {
				rect.FillColor = color.NRGBA{70, 70, 150, 255} // Visited
			} else {
				rect.FillColor = color.NRGBA{100, 50, 100, 255} // Unvisited
			}
			row.Add(rect)
		}
		grid.Add(row)
	}
	return grid
}

func movePlayer(p *Player, dx, dy int, infoLabel, logLabel *widget.Label, grid *fyne.Container, w fyne.Window) {
	newX, newY := p.X+dx, p.Y+dy
	if newX < 0 || newX >= gridSize || newY < 0 || newY >= gridSize {
		logLabel.SetText("You bump into a wall! 🧱")
		return
	}

	p.X, p.Y = newX, newY
	p.Visited[fmt.Sprintf("%d,%d", newX, newY)] = true

	event := randomEvent(p, logLabel, w)
	if event != "" {
		logLabel.SetText(event)
	} else {
		logLabel.SetText(fmt.Sprintf("You moved to room (%d, %d).", p.X, p.Y))
	}

	updateInfo(p, infoLabel)
	newGrid := drawDungeon(p)
	grid.Objects = newGrid.Objects
	grid.Refresh()
}

func drawHPBar(p *Player) *canvas.Rectangle {
	ratio := float64(p.HP) / float64(p.MaxHP)
	width := float32(300 * ratio)
	bar := canvas.NewRectangle(color.NRGBA{255, 0, 0, 255})
	bar.SetMinSize(fyne.NewSize(width, 15))
	return bar
}

func updateInfo(p *Player, label *widget.Label) {
	label.SetText(fmt.Sprintf("❤️ HP: %d/%d | 💰 Gold: %d | 🧪 Potions: %d | 📍 (%d,%d)",
		p.HP, p.MaxHP, p.Gold, p.Potions, p.X, p.Y))
}

func randomEvent(p *Player, logLabel *widget.Label, w fyne.Window) string {
	r := rand.Intn(100)

	switch {
	case r < 35:
		enemy := Enemy{Name: randomEnemy(), HP: 20 + rand.Intn(20), Attack: 5 + rand.Intn(5)}
		combat(p, enemy, logLabel, w)
		return ""
	case r < 65:
		goldFound := 5 + rand.Intn(15)
		p.Gold += goldFound
		return fmt.Sprintf("💰 You found %d gold!", goldFound)
	case r < 85:
		p.Potions++
		return "🧪 You found a potion!"
	default:
		return "The room is eerily quiet..."
	}
}

func combat(p *Player, e Enemy, logLabel *widget.Label, w fyne.Window) {
	logLabel.SetText(fmt.Sprintf("⚔️ A wild %s appears!", e.Name))

	for e.HP > 0 && p.HP > 0 {
		damageToEnemy := 10 + rand.Intn(5)
		damageToPlayer := e.Attack + rand.Intn(3)

		e.HP -= damageToEnemy
		if e.HP <= 0 {
			logLabel.SetText(fmt.Sprintf("💀 You defeated the %s! +10 gold", e.Name))
			p.Gold += 10
			return
		}

		p.HP -= damageToPlayer
		if p.HP <= 0 {
			showGameOver(w)
			return
		}
	}
}

func usePotion(p *Player, logLabel, infoLabel *widget.Label) {
	if p.Potions <= 0 {
		logLabel.SetText("❌ You have no potions left!")
		return
	}

	p.Potions--
	heal := 25
	p.HP += heal
	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}

	logLabel.SetText(fmt.Sprintf("💖 You drink a potion and recover %d HP.", heal))
	updateInfo(p, infoLabel)
}

func showGameOver(w fyne.Window) {
	label := widget.NewLabelWithStyle("☠️ You died! Game Over.", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	restartBtn := widget.NewButton("Restart", func() {
		startGame(w)
	})
	content := container.NewVBox(label, restartBtn)
	w.SetContent(container.NewCenter(content))
}

func randomEnemy() string {
	enemies := []string{"Goblin", "Skeleton", "Orc", "Bat"}
	return enemies[rand.Intn(len(enemies))]
}
