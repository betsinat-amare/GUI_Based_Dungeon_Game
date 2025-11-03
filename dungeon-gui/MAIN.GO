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

type Player struct {
	X, Y    int
	HP      int
	MaxHP   int
	Gold    int
	Potions int
}

type Enemy struct {
	Name   string
	HP     int
	Attack int
}

const gridSize = 3

func main() {
	rand.Seed(time.Now().UnixNano())

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

func startGame(w fyne.Window) {
	player := &Player{X: 1, Y: 1, HP: 100, MaxHP: 100, Gold: 0, Potions: 1}

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

	movement := container.NewGridWithColumns(3,
		btnLeft,
		btnDown,
		btnRight,
	)

	statusBox := container.NewVBox(
		infoLabel,
		btnPotion,
		canvas.NewRectangle(color.NRGBA{200, 200, 200, 255}),
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

func drawDungeon(player *Player) *fyne.Container {
	grid := container.NewGridWithColumns(gridSize)

	for y := 0; y < gridSize; y++ {
		row := container.NewGridWithColumns(gridSize)
		for x := 0; x < gridSize; x++ {
			rect := canvas.NewRectangle(color.NRGBA{80, 80, 80, 255})
			rect.SetMinSize(fyne.NewSize(60, 60))
			if x == player.X && y == player.Y {
				rect.FillColor = color.NRGBA{0, 255, 0, 255}
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

func updateInfo(p *Player, label *widget.Label) {
	label.SetText(fmt.Sprintf("❤️ HP: %d/%d  |  💰 Gold: %d  |  🧪 Potions: %d  |  Position: (%d,%d)", p.HP, p.MaxHP, p.Gold, p.Potions, p.X, p.Y))
}

// -------------------------
// Random events and combat
// -------------------------
func randomEvent(p *Player, logLabel *widget.Label, w fyne.Window) string {
	r := rand.Intn(100)

	switch {
	case r < 40:
		enemy := Enemy{Name: randomEnemy(), HP: 20 + rand.Intn(20), Attack: 5 + rand.Intn(5)}
		combat(p, enemy, logLabel, w)
		return ""
	case r < 70:
		goldFound := 5 + rand.Intn(15)
		p.Gold += goldFound
		return fmt.Sprintf("💰 You found %d gold!", goldFound)
	case r < 90:
		p.Potions++
		return "🧪 You found a potion!"
	default:
		return "The room is empty and silent..."
	}
}

func combat(p *Player, e Enemy, logLabel *widget.Label, w fyne.Window) {
	logLabel.SetText(fmt.Sprintf("⚔️ A wild %s appears!", e.Name))

	for e.HP > 0 && p.HP > 0 {
		damageToEnemy := 10 + rand.Intn(5)
		damageToPlayer := e.Attack + rand.Intn(3)

		e.HP -= damageToEnemy
		if e.HP <= 0 {
			logLabel.SetText(fmt.Sprintf("💀 You defeated the %s!", e.Name))
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
	w.SetContent(content)
}

func randomEnemy() string {
	enemies := []string{"Goblin", "Skeleton", "Orc", "Bat"}
	return enemies[rand.Intn(len(enemies))]
}
