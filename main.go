package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var bgImg *ebiten.Image
var blackImg *ebiten.Image

var ROWS = 20
var COLUMNS = 20

const SQUARE_SIZE = 32

var playing bool

var board [][]bool

func init() {
	playing = false
	initBoard()
	initImages()
}

func initBoard() {
	board = make([][]bool, ROWS)
	for i := range board {
		board[i] = make([]bool, COLUMNS)
	}
}

func initImages() {
	// background image
	bgImg = ebiten.NewImage(COLUMNS*SQUARE_SIZE, ROWS*SQUARE_SIZE)
	bgImg.Fill(color.White)
	for i := 1; i < COLUMNS; i++ {
		// vertical lines
		vector.StrokeLine(bgImg, float32(i*SQUARE_SIZE), 0, float32(i*SQUARE_SIZE), float32(ROWS*SQUARE_SIZE), 1.0, color.Black, false)
	}
	for i := 1; i < ROWS; i++ {
		// horizontal lines
		vector.StrokeLine(bgImg, 0, float32(i*SQUARE_SIZE), float32(COLUMNS*SQUARE_SIZE), float32(i*SQUARE_SIZE), 1.0, color.Black, false)
	}

	// black image for live cells
	blackImg = ebiten.NewImage(SQUARE_SIZE, SQUARE_SIZE)
	blackImg.Fill(color.Black)
}

type Game struct{}

func (g *Game) Update() error {
	if !playing {
		playing = drawingPhase()
	} else {
		playing = playingPhase()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(bgImg, nil)
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLUMNS; j++ {
			if board[i][j] {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(j*SQUARE_SIZE), float64(i*SQUARE_SIZE))
				screen.DrawImage(blackImg, op)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return COLUMNS * SQUARE_SIZE, ROWS * SQUARE_SIZE
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func drawingPhase() bool {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		board[y/SQUARE_SIZE][x/SQUARE_SIZE] = true
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		board[y/SQUARE_SIZE][x/SQUARE_SIZE] = false
		//println("x:", x, "y:", y, "column:", x/SQUARE_SIZE, "row:", y/SQUARE_SIZE)
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return true
	}
	return false
}

func playingPhase() bool {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		return false
	}
	return true
}
