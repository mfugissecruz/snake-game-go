package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/nsf/termbox-go"
)

type Point struct {
	X int
	Y int
}

type Snake struct {
	Body      []Point
	Direction string
}

type FoodType int

const (
	NormalFood FoodType = iota
	PowerUpFood
)

type Food struct {
	Position Point
	Type     FoodType
}

type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StateGameOver
)

type Game struct {
	Snake      Snake
	Food       Food
	Score      int
	HighScore  int
	GameOver   bool
	Width      int
	Height     int
	State      GameState
	Level      int
	Speed      time.Duration
	FrameCount int
	Obstacles  []Point
}

type ToneGenerator struct {
	freq float64
	pos  float64
	sr   beep.SampleRate
}

func NewTone(sr beep.SampleRate, freq float64) *ToneGenerator {
	return &ToneGenerator{
		freq: freq,
		sr:   sr,
	}
}

func (t *ToneGenerator) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		v := math.Sin(t.pos * 2 * math.Pi * t.freq / float64(t.sr))
		samples[i][0] = v
		samples[i][1] = v
		t.pos++
	}
	return len(samples), true
}

func (t *ToneGenerator) Err() error {
	return nil
}

var soundInitialized = false

func initSound() {
	if !soundInitialized {
		sr := beep.SampleRate(44100)
		speaker.Init(sr, sr.N(time.Second/10))
		soundInitialized = true
	}
}

func playTone(freq float64, duration time.Duration) {
	if !soundInitialized {
		return
	}

	sr := beep.SampleRate(44100)
	tone := NewTone(sr, freq)
	sound := beep.Take(sr.N(duration), tone)

	done := make(chan bool)
	speaker.Play(beep.Seq(sound, beep.Callback(func() {
		done <- true
	})))

	go func() {
		<-done
	}()
}

func soundEat() {
	go playTone(800, 50*time.Millisecond)
}

func soundPowerUp() {
	go func() {
		playTone(600, 100*time.Millisecond)
		time.Sleep(50 * time.Millisecond)
		playTone(800, 100*time.Millisecond)
		time.Sleep(50 * time.Millisecond)
		playTone(1000, 100*time.Millisecond)
	}()
}

func soundLevelUp() {
	go func() {
		playTone(1000, 100*time.Millisecond)
		time.Sleep(80 * time.Millisecond)
		playTone(1200, 100*time.Millisecond)
	}()
}

func soundGameOver() {
	go func() {
		playTone(400, 200*time.Millisecond)
		time.Sleep(100 * time.Millisecond)
		playTone(300, 200*time.Millisecond)
		time.Sleep(100 * time.Millisecond)
		playTone(200, 300*time.Millisecond)
	}()
}

func LoadHighScore() int {
	data, err := os.ReadFile("highscore.txt")
	if err != nil {
		return 0
	}

	score, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0
	}

	return score
}

func SaveHighScore(score int) error {
	return os.WriteFile("highscore.txt", []byte(fmt.Sprintf("%d", score)), 0644)
}

func NewGame() *Game {
	game := &Game{
		Snake: Snake{
			Body: []Point{
				{X: 10, Y: 10},
				{X: 9, Y: 10},
				{X: 8, Y: 10},
			},
			Direction: "right",
		},
		Score:      0,
		HighScore:  LoadHighScore(),
		GameOver:   false,
		Width:      40,
		Height:     20,
		State:      StateMenu,
		Level:      1,
		Speed:      150 * time.Millisecond,
		FrameCount: 0,
		Obstacles:  []Point{},
	}
	game.GenerateFood()
	game.GenerateObstacles()
	return game
}

func (g *Game) Reset() {
	g.Snake = Snake{
		Body: []Point{
			{X: 10, Y: 10},
			{X: 9, Y: 10},
			{X: 8, Y: 10},
		},
		Direction: "right",
	}
	g.Score = 0
	g.GameOver = false
	g.State = StatePlaying
	g.Level = 1
	g.Speed = 150 * time.Millisecond
	g.FrameCount = 0
	g.Obstacles = []Point{}
	g.GenerateFood()
	g.GenerateObstacles()
}

func (g *Game) UpdateLevel() {
	newLevel := (g.Score / 50) + 1

	if newLevel > g.Level {
		g.Level = newLevel
		g.Speed = time.Duration(150-((g.Level-1)*10)) * time.Millisecond
		if g.Speed < 50*time.Millisecond {
			g.Speed = 50 * time.Millisecond
		}
		g.GenerateObstacles()
	}
}

func (g *Game) CheckAndSaveHighScore() bool {
	if g.Score > g.HighScore {
		g.HighScore = g.Score
		SaveHighScore(g.HighScore)
		return true
	}
	return false
}

func (g *Game) IsPositionSafe(pos Point) bool {
	for _, chunk := range g.Snake.Body {
		if pos.X == chunk.X && pos.Y == chunk.Y {
			return false
		}
	}

	if pos.X == g.Food.Position.X && pos.Y == g.Food.Position.Y {
		return false
	}

	for _, obs := range g.Obstacles {
		if pos.X == obs.X && pos.Y == obs.Y {
			return false
		}
	}

	startX, startY := 10, 10
	if pos.X >= startX-2 && pos.X <= startX+2 &&
		pos.Y >= startY-2 && pos.Y <= startY+2 {
		return false
	}

	return true
}

func (g *Game) GenerateObstacles() {
	g.Obstacles = []Point{}

	numObstacles := g.Level * 2
	if numObstacles > 20 {
		numObstacles = 20
	}

	for i := 0; i < numObstacles; i++ {
		for attempts := 0; attempts < 50; attempts++ {
			pos := Point{
				X: rand.Intn(g.Width-2) + 1,
				Y: rand.Intn(g.Height-2) + 1,
			}

			if g.IsPositionSafe(pos) {
				g.Obstacles = append(g.Obstacles, pos)
				break
			}
		}
	}
}

func (g *Game) GenerateFood() {
	var position Point

	for attempts := 0; attempts < 100; attempts++ {
		position = Point{
			X: rand.Intn(g.Width-2) + 1,
			Y: rand.Intn(g.Height-2) + 1,
		}

		if g.IsPositionSafe(position) {
			break
		}
	}

	foodType := NormalFood
	if rand.Intn(100) < 20 {
		foodType = PowerUpFood
	}

	g.Food = Food{
		Position: position,
		Type:     foodType,
	}
}

func (g *Game) MoveSnake() {
	head := g.Snake.Body[0]
	newHead := Point{X: head.X, Y: head.Y}

	switch g.Snake.Direction {
	case "up":
		newHead.Y--
	case "down":
		newHead.Y++
	case "left":
		newHead.X--
	case "right":
		newHead.X++
	}

	if g.CheckWallCollision(newHead) ||
		g.CheckSelfCollision(newHead) ||
		g.CheckObstacleCollision(newHead) {
		g.GameOver = true
		g.State = StateGameOver
		g.CheckAndSaveHighScore()
		soundGameOver()
		return
	}

	g.Snake.Body = append([]Point{newHead}, g.Snake.Body...)

	if newHead.X == g.Food.Position.X && newHead.Y == g.Food.Position.Y {
		points := 10
		if g.Food.Type == PowerUpFood {
			points = 50
			soundPowerUp()
		} else {
			soundEat()
		}

		oldLevel := g.Level
		g.Score += points
		g.UpdateLevel()

		if g.Level > oldLevel {
			soundLevelUp()
		}

		g.GenerateFood()
	} else {
		g.Snake.Body = g.Snake.Body[:len(g.Snake.Body)-1]
	}
}

func (g *Game) CheckWallCollision(p Point) bool {
	return p.X <= 0 || p.X >= g.Width-1 || p.Y <= 0 || p.Y >= g.Height-1
}

func (g *Game) CheckSelfCollision(head Point) bool {
	for _, chunk := range g.Snake.Body {
		if head.X == chunk.X && head.Y == chunk.Y {
			return true
		}
	}
	return false
}

func (g *Game) CheckObstacleCollision(p Point) bool {
	for _, obs := range g.Obstacles {
		if p.X == obs.X && p.Y == obs.Y {
			return true
		}
	}
	return false
}

func (g *Game) DrawMenu() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	title := []string{
		"          ____  _   _    _    _  ________ ",
		"         / ___|| \\ | |  / \\  | |/ / ____| ",
		"         \\___ \\|  \\| | / _ \\ | ' /|  _|  ",
		"          ___) | |\\  |/ ___ \\| . \\| |___  ",
		"         |____/|_| \\_/_/   \\_\\_|\\_\\_____|",
	}

	menu := []string{
		"  ╔═════════════ ═╗",
		"  ║                                           ║",
		fmt.Sprintf("  ║         ★ RECORDE: %-21d║", g.HighScore),
		"  ║                                           ║",
		"  ║  CONTROLES:                               ║",
		"  ║    Setas : Movimentar                     ║",
		"  ║    ENTER : Iniciar jogo                   ║",
		"  ║    R     : Reiniciar                      ║",
		"  ║    ESC   : Sair                           ║",
		"  ║                                           ║",
		"  ║  REGRAS:                                  ║",
		"  ║    ◆ Comida normal ....... 10 pontos    ║",
		"  ║    ★ Power-up ............ 50 pontos    ║",
		"  ║    ▓ Obstaculos .......... Evite!       ║",
		"  ║                                           ║",
		"  ║  A cada 50 pontos = +1 nivel              ║",
		"  ║  Mais nivel = Mais rapido + obstaculos    ║",
		"  ║                                           ║",
		"  ║      Pressione ENTER para comecar         ║",
		"  ║                                           ║",
		"  ╚══════════════ ╝",
	}

	startY := 3
	startX := 2

	for i, line := range title {
		for j, char := range line {
			termbox.SetCell(startX+j, startY+i, char, termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault)
		}
	}

	menuStartY := startY + len(title) + 1
	for i, line := range menu {
		color := termbox.ColorCyan
		if i == 2 {
			color = termbox.ColorYellow
		}
		if i == len(menu)-2 {
			color = termbox.ColorYellow | termbox.AttrBold
		}
		for j, char := range line {
			termbox.SetCell(startX+j, menuStartY+i, char, color, termbox.ColorDefault)
		}
	}

	termbox.Flush()
}

func (g *Game) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for x := 0; x < g.Width; x++ {
		termbox.SetCell(x, 0, '═', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, g.Height-1, '═', termbox.ColorWhite, termbox.ColorDefault)
	}

	for y := 0; y < g.Height; y++ {
		termbox.SetCell(0, y, '║', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(g.Width-1, y, '║', termbox.ColorWhite, termbox.ColorDefault)
	}

	termbox.SetCell(0, 0, '╔', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(g.Width-1, 0, '╗', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(0, g.Height-1, '╚', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(g.Width-1, g.Height-1, '╝', termbox.ColorWhite, termbox.ColorDefault)

	for _, obs := range g.Obstacles {
		termbox.SetCell(obs.X, obs.Y, '▓', termbox.ColorWhite, termbox.ColorDefault)
	}

	for i, chunk := range g.Snake.Body {
		char := '█'
		color := termbox.ColorGreen

		if i == 0 {
			char = '●'
			color = termbox.ColorYellow
		}

		termbox.SetCell(chunk.X, chunk.Y, char, color, termbox.ColorDefault)
	}

	foodChar := '◆'
	foodColor := termbox.ColorRed

	if g.Food.Type == PowerUpFood {
		foodChar = '★'
		foodColor = termbox.ColorYellow
		if (g.FrameCount/5)%2 == 0 {
			foodColor = termbox.ColorMagenta
		}
	}

	termbox.SetCell(g.Food.Position.X, g.Food.Position.Y, foodChar, foodColor, termbox.ColorDefault)

	msg := fmt.Sprintf(" Pontos: %d | Recorde: %d | Nivel: %d | Tamanho: %d ",
		g.Score, g.HighScore, g.Level, len(g.Snake.Body))
	for i, char := range msg {
		termbox.SetCell(i+2, g.Height, char, termbox.ColorCyan, termbox.ColorDefault)
	}

	termbox.Flush()
}

func (g *Game) DrawGameOver() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	isNewRecord := g.Score >= g.HighScore && g.Score > 0

	var messages []string

	if isNewRecord {
		messages = []string{
			"╔═══════════════════════════╗",
			"║     GAME OVER!            ║",
			"║                           ║",
			"║  ★ NOVO RECORDE! ★        ║",
			"║                           ║",
			fmt.Sprintf("║  Pontos: %-16d║", g.Score),
			fmt.Sprintf("║  Nivel: %-17d║", g.Level),
			fmt.Sprintf("║  Tamanho: %-15d║", len(g.Snake.Body)),
			"║                           ║",
			"║  Pressione R - Reiniciar  ║",
			"║  Pressione ESC - Sair     ║",
			"╚═══════════════════════════╝",
		}
	} else {
		messages = []string{
			"╔═════════╗",
			"║     GAME OVER!            ║",
			"║                           ║",
			fmt.Sprintf("║  Pontos: %-16d ║", g.Score),
			fmt.Sprintf("║  Recorde: %-15d ║", g.HighScore),
			fmt.Sprintf("║  Nivel: %-17d ║", g.Level),
			fmt.Sprintf("║  Tamanho: %-15d ║", len(g.Snake.Body)),
			"║                           ║",
			"║  Pressione R - Reiniciar  ║",
			"║  Pressione ESC - Sair     ║",
			"╚═════════╝",
		}
	}

	startX := g.Width/2 - 14
	startY := g.Height/2 - len(messages)/2

	for i, msg := range messages {
		color := termbox.ColorRed
		if isNewRecord && (i == 3) {
			color = termbox.ColorYellow
		}
		for j, char := range msg {
			termbox.SetCell(startX+j, startY+i, char, color, termbox.ColorDefault)
		}
	}

	termbox.Flush()
}

func (g *Game) HandleInput(end chan bool) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				end <- true
				return
			}

			if ev.Key == termbox.KeyEnter && g.State == StateMenu {
				g.State = StatePlaying
			}

			if (ev.Ch == 'r' || ev.Ch == 'R') && g.State == StateGameOver {
				g.Reset()
			}

			if g.State == StatePlaying {
				switch ev.Key {
				case termbox.KeyArrowUp:
					if g.Snake.Direction != "down" {
						g.Snake.Direction = "up"
					}
				case termbox.KeyArrowDown:
					if g.Snake.Direction != "up" {
						g.Snake.Direction = "down"
					}
				case termbox.KeyArrowLeft:
					if g.Snake.Direction != "right" {
						g.Snake.Direction = "left"
					}
				case termbox.KeyArrowRight:
					if g.Snake.Direction != "left" {
						g.Snake.Direction = "right"
					}
				}
			}
		}
	}
}

func main() {
	initSound()

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	game := NewGame()
	end := make(chan bool)

	go game.HandleInput(end)

	ticker := time.NewTicker(game.Speed)
	defer ticker.Stop()

	lastSpeed := game.Speed

	for {
		select {
		case <-end:
			speaker.Close()
			return
		case <-ticker.C:
			if game.Speed != lastSpeed {
				ticker.Stop()
				ticker = time.NewTicker(game.Speed)
				lastSpeed = game.Speed
			}

			game.FrameCount++

			switch game.State {
			case StateMenu:
				game.DrawMenu()
			case StatePlaying:
				game.MoveSnake()
				game.Draw()
			case StateGameOver:
				game.DrawGameOver()
			}
		}
	}
}
