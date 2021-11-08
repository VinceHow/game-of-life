package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
	"time"
)

const (
	screenWidth  = 1000
	screenHeight = 700
	gridSize     = 20
)

func getSpeed() [8]float64 {
	possibleSpeeds := [8]float64{0, 0.25, 0.5, 1, 2, 3, 4, 5}
	return possibleSpeeds
}

func getPresets() [][]Cell {
	var presets [][]Cell
	set1 :=  make([]Cell,0)
	presets = append(presets, set1) // default is empty board
	set2 :=  []Cell{[2]int{17,24},[2]int{17,25},[2]int{17,26}}
	presets = append(presets, set2) // set 2 is a basic oscillator
	return presets
}

func getBoard(cells []Cell) board {
	b := InitializeBoard()
	for _, c := range cells {
		b[c[0]][c[1]] = true
	}
	return b
}

type Cell [2]int // [2]int{1,2} means cell in row 1, col 2 is alive

type Game struct {
	board  					board
	liveCells 				[]Cell
	generations 			int64
	evolutionSpeed 			int // tiers of speeds taken, with 0 being first tier
	gamePaused				bool
	presetId				int
}

func (g *Game) getLiveCells() []Cell {
	// converts a board into a list of alive cell positions
	var liveCells []Cell
	for r, row := range g.board {
		for c, v := range row {
			if v {
				liveCells = append(liveCells, [2]int{r,c})
			}
		}
	}
	return liveCells
}

func (g *Game) Update() error {
	// allow user to edit board when game is paused
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && g.gamePaused {
		x, y := ebiten.CursorPosition()
		row := y/gridSize
		col := x/gridSize
		if g.board[row][col] {
			g.board[row][col] = false
		} else if g.board[row][col] == false {
			g.board[row][col] = true
		}
	}

	// user to switch between interesting presets
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) && g.gamePaused {
		if g.presetId+1 == len(getPresets()) {
			g.presetId = 0
			g.liveCells = getPresets()[g.presetId]
			g.board = getBoard(g.liveCells)
		} else {
			g.presetId ++
			g.liveCells = getPresets()[g.presetId]
			g.board = getBoard(g.liveCells)
		}
	}
	// user to control game params
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if g.evolutionSpeed > 0 { // cannot reduce speed to below 0
			g.evolutionSpeed --
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if g.evolutionSpeed < len(getSpeed())-1 { // cannot increase above top speed
			g.evolutionSpeed ++
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) { // pause game
		if g.gamePaused {
			g.gamePaused = false
		} else {
			g.gamePaused = true
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) { // reset game
		g.reset()
	}

	// update board state
	if len(g.liveCells) > 0 && !g.gamePaused{
		g.generations ++
		g.board = UpdateBoard(g.board)
		time.Sleep(time.Duration(int64(1000000000/getSpeed()[g.evolutionSpeed]))) // convert speed into wait time in nanoseconds
	}
	g.liveCells = g.getLiveCells()
	return nil
}

func (g *Game) reset() {
	// clear all live cells
	var b board = InitializeBoard()
	g.board = b
	g.generations = 0
	g.evolutionSpeed = 3
	g.gamePaused = true
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draws the world of cells
	for _, c := range g.liveCells {
		ebitenutil.DrawRect(screen, float64(c[1]*gridSize)+1, float64(c[0]*gridSize)+1, gridSize-2, gridSize-2, color.RGBA{R: 0x80, G: 0xa0, B: 0xc0, A: 0xff})
	}
	// game information
	msg := fmt.Sprintf("Number of cells alive: %v", len(g.liveCells))
	msg += fmt.Sprintf("\nNumber of generations: %v", g.generations)
	msg += fmt.Sprintf("\nCurrent speed: %v", getSpeed()[g.evolutionSpeed])
	msg += fmt.Sprintf("\nGame is currently paused: %v", g.gamePaused)
	msg += fmt.Sprintf("\nBoard preset ID: %v", g.presetId)
	mx, my := ebiten.CursorPosition()
	msg += fmt.Sprintf("\n(%d, %d)", mx, my)
	ebitenutil.DebugPrint(screen, msg)
}

func NewGame() *Game {
	var b board = InitializeBoard()
	g := &Game{
		generations:    0,
		evolutionSpeed: 3,
		gamePaused:     true,
		board: 			b,
		presetId: 		0,
	}
	return g
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	/*
	Start of game:
		1. start with black screen with prompt for user to add live cells
		2. press space once finish adding cells to start game
	Once game is running:
		1. update screen for each generation
		2. up and down arrow to adjust speed of evolution
		3. press space to pause/unpause
	Game is terminated if:
		1. all cells are dead
		OR
		2. user presses Escape
	Game terminated:
		1. show number of generations ran
		2. show current cells alive
		3. prompt to start new round
	 */
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Game of Life")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

