package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type Board struct {
	data   []bool
	width  int
	height int
}

func NewBoard(width, height int, maxLiveCells int) *Board {
	b := &Board{
		data:   make([]bool, width*height),
		width:  width,
		height: height,
	}

	b.init(maxLiveCells)
	return b
}

func (b *Board) init(maxLiveCells int) {
	for i := 0; i < maxLiveCells; i++ {
		x := rand.Intn(b.width)
		y := rand.Intn(b.height)

		b.data[x+y*b.width] = true
	}
}

func (b *Board) neighbors(x, y int) int {
	// Run clockwise around the given cell, if a cell is alive, increase count
	cnt := 0
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			// If we're at the center point, ignore
			if i == 0 && j == 0 {
				continue
			}

			x2 := x + i
			y2 := y + j

			if x2 < 0 || y2 < 0 || b.width <= x2 || b.height <= y2 {
				continue
			}

			if b.data[y2*b.width+x2] {
				cnt++
			}
		}
	}

	return cnt
}

func (b *Board) update() {
	// For each position on the board,
	// Determine moore neighbors
	// Conway Rules determine the value of the cell
	next := make([]bool, b.width*b.height)

	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			nCnt := b.neighbors(x, y)
			//pos := &(b.data[x*y+b.width])

			switch {
			// Any live cell with 2 or 3 live neighbors survives
			case (b.data[y*b.width+x] && nCnt == 2 || nCnt == 3):
				next[y*b.width+x] = true
			// Any dead cell with three live neighbors becomes alive
			case (b.data[x*y+b.width] == false && nCnt == 3):
				next[y*b.width+x] = true
			// All other live cells die. All other dead cells stay dead
			default:
				next[y*b.width+x] = false
			}
		}
	}

	b.data = next
}

func (b *Board) draw(pixels []byte) {
	for i, p := range b.data {
		if p {
			pixels[4*i] = 0xff
			pixels[4*i+1] = 0xff
			pixels[4*i+2] = 0xff
			pixels[4*i+3] = 0xff
		} else {
			pixels[4*i] = 0
			pixels[4*i+1] = 0
			pixels[4*i+2] = 0
			pixels[4*i+3] = 0
		}
	}
}

type Game struct {
	world  *Board
	pixels []byte
}

func (g *Game) Update(screen *ebiten.Image) error {
	// Update the logical state
	g.world.update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Render screen
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}

	g.world.draw(g.pixels)
	screen.ReplacePixels(g.pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Return game screen size
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(42)

	game := &Game{
		world: NewBoard(screenWidth, screenHeight, 5000),
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Game title")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
