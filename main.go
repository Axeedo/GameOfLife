package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var bgImg *ebiten.Image
var blackImg *ebiten.Image

var ROWS int = 300
var COLUMNS int = 300

const SQUARE_SIZE int = 32

var prevPlaying bool
var start time.Time
var simulationSpeed int = 300
var iterCount int
var prevIterTime time.Time

var playing bool
var showFPS bool = false

var board [][]bool

func main() {
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetVsyncEnabled(true)

	//ebiten.SetWindowSize(1280, 720)
	//ebiten.SetWindowSize(1920, 1080)
	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Game of Life")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func init() {
	iterCount = 0
	start = time.Now()
	prevIterTime = time.Now()
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
	/*for i := 1; i < COLUMNS; i++ {
		// vertical lines
		vector.StrokeLine(bgImg, float32(i*SQUARE_SIZE), 0, float32(i*SQUARE_SIZE), float32(ROWS*SQUARE_SIZE), 1.0, color.Black, false)
	}
	for i := 1; i < ROWS; i++ {
		// horizontal lines
		vector.StrokeLine(bgImg, 0, float32(i*SQUARE_SIZE), float32(COLUMNS*SQUARE_SIZE), float32(i*SQUARE_SIZE), 1.0, color.Black, false)
	}*/

	// black image for live cells
	blackImg = ebiten.NewImage(SQUARE_SIZE, SQUARE_SIZE)
	blackImg.Fill(color.Black)
}

type Game struct{}

func (g *Game) Update() error {
	prevPlaying = playing

	if !playing {
		playing = drawingPhase()
	} else {
		playing = playingPhase()

	}
	if prevPlaying != playing {
		if !playing {
			println("Drawing phase")
		} else {
			println("Playing phase")
		}
	}
	if time.Since(start) >= time.Second && showFPS {
		start = time.Now()
		fmt.Printf("Fps: %.2f\t", ebiten.ActualTPS())
		fmt.Printf("Tps: %.2f\n", ebiten.ActualTPS())
		println("Iteration count per second:", iterCount)
		iterCount = 0
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
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		println("Reseting board...")
		initBoard() //reset board
	}
	return false
}

func playingPhase() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return false
	}
	if time.Since(prevIterTime) >= time.Duration(1000/simulationSpeed)*time.Millisecond {
		iterCount++
		board = updateBoard(board)
		prevIterTime = time.Now()
	}
	return true
}

func liveNeighbors(board [][]bool, row int, column int) int {
	count := 0
	// first row
	if row > 0 {
		if column > 0 && board[row-1][column-1] {
			count++
		}
		if board[row-1][column] {
			count++
		}
		if column < COLUMNS-1 && board[row-1][column+1] {
			count++
		}
	}
	// second row
	if column > 0 && board[row][column-1] {
		count++
	}
	if column < COLUMNS-1 && board[row][column+1] {
		count++
	}
	//third row
	if row < ROWS-1 {
		if column > 0 && board[row+1][column-1] {
			count++
		}
		if board[row+1][column] {
			count++
		}
		if column < COLUMNS-1 && board[row+1][column+1] {
			count++
		}
	}
	return count
}

func updateBoard(board [][]bool) [][]bool {
	newboard := make([][]bool, ROWS)
	for i := range newboard {
		newboard[i] = make([]bool, COLUMNS)
	}
	for row := range board {
		for column := range board[row] {
			nb := liveNeighbors(board, row, column)
			if (board[row][column] && (nb == 2 || nb == 3)) || !board[row][column] && nb == 3 {
				// remain alive
				newboard[row][column] = true
			} else {
				newboard[row][column] = false
			}
		}
	}

	return newboard
}
