package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ScreenHeight = 800
	ScreenWidth  = 800
	GridSize     = 40
	GridWidth    = ScreenWidth / GridSize
	GridHeight   = ScreenHeight / GridSize
	StepDelay    = 1
)

type Cell struct {
	visited bool
	walls   [4]bool
}

type Game struct {
	grid    [][]*Cell
	stack   [][2]int
	current [2]int
	frame   int
}

func NewGame() (*Game, error) {
	grid := make([][]*Cell, GridHeight)
	for y := range grid {
		grid[y] = make([]*Cell, GridWidth)
		for x := range grid[y] {
			grid[y][x] = &Cell{walls: [4]bool{true, true, true, true}}
		}
	}

	game := &Game{
		grid:    grid,
		stack:   [][2]int{{0, 0}},
		current: [2]int{0, 0},
	}
	game.grid[0][0].visited = true
	return game, nil
}

func (g *Game) Update() error {
	g.frame++
	if g.frame >= StepDelay {
		g.frame = 0
		g.step()
	}
	return nil
}

func (g *Game) step() {
	if len(g.stack) > 0 {
		x, y := g.current[0], g.current[1]
		neighbors := g.getUnvisitedNeighbors(x, y)

		if len(neighbors) > 0 {
			next := neighbors[rand.Intn(len(neighbors))]
			g.removeWall(x, y, next[0], next[1])
			g.grid[next[1]][next[0]].visited = true
			g.stack = append(g.stack, next)
			g.current = next
		} else {
			g.stack = g.stack[:len(g.stack)-1]
			if len(g.stack) > 0 {
				g.current = g.stack[len(g.stack)-1]
			}
		}
	}
}

func (g *Game) getUnvisitedNeighbors(x, y int) [][2]int {
	directions := [][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	var neighbors [][2]int

	for _, d := range directions {
		nx, ny := x+d[0], y+d[1]
		if nx >= 0 && nx < GridWidth && ny >= 0 && ny < GridHeight && !g.grid[ny][nx].visited {
			neighbors = append(neighbors, [2]int{nx, ny})
		}
	}

	return neighbors
}

func (g *Game) removeWall(x1, y1, x2, y2 int) {
	if x1-x2 == 1 {
		g.grid[y1][x1].walls[3] = false
		g.grid[y2][x2].walls[1] = false
	} else if x1-x2 == -1 {
		g.grid[y1][x1].walls[1] = false
		g.grid[y2][x2].walls[3] = false
	} else if y1-y2 == 1 {
		g.grid[y1][x1].walls[0] = false
		g.grid[y2][x2].walls[2] = false
	} else if y1-y2 == -1 {
		g.grid[y1][x1].walls[2] = false
		g.grid[y2][x2].walls[0] = false
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y := 0; y < GridHeight; y++ {
		for x := 0; x < GridWidth; x++ {
			cell := g.grid[y][x]
			px, py := float32(x*GridSize), float32(y*GridSize)

			bgColor := color.RGBA{255, 255, 255, 255}
			if cell.visited {
				bgColor = color.RGBA{200, 200, 255, 255}
			}
			vector.DrawFilledRect(screen, px, py, GridSize, GridSize, bgColor, false)

			wallColor := color.Black
			if cell.walls[0] {
				vector.StrokeLine(screen, px, py, px+GridSize, py, 1, wallColor, false)
			}
			if cell.walls[1] {
				vector.StrokeLine(screen, px+GridSize, py, px+GridSize, py+GridSize, 1, wallColor, false)
			}
			if cell.walls[2] {
				vector.StrokeLine(screen, px, py+GridSize, px+GridSize, py+GridSize, 1, wallColor, false)
			}
			if cell.walls[3] {
				vector.StrokeLine(screen, px, py, px, py+GridSize, 1, wallColor, false)
			}
		}
	}

	cx, cy := float32(g.current[0]*GridSize), float32(g.current[1]*GridSize)
	vector.DrawFilledRect(screen, cx, cy, GridSize, GridSize, color.RGBA{255, 100, 100, 255}, false)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	game, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Maze Generator")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
