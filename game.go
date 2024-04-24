package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

func NewGrid() map[Hex]bool {
	g := make(map[Hex]bool)
	N := 10
	for q := -N; q <= N; q++ {
		r1 := Max(-N, -q-N)
		r2 := Min(N, -q+N)
		for r := r1; r <= r2; r++ {
			g[Hex{q, r, -q - r}] = false
		}
	}
	return g
}

type Game struct {
	layout     Layout
	grid       map[Hex]bool
	generation int
	running    bool
	lastTick   time.Time
	lifeCount  int
}

func NewGame() *Game {
	g := &Game{
		layout:     *NewLayout(&PointyTop, Point{25, 25}, Point{500, 500}),
		grid:       NewGrid(),
		generation: 0,
		running:    false,
		lastTick:   time.Now(),
		lifeCount:  0,
	}
	return g
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) handleMouseButton(button ebiten.MouseButton, grid *map[Hex]bool) {
	if ebiten.IsMouseButtonPressed(button) {
		x, y := ebiten.CursorPosition()
		mouse := Point{float32(x), float32(y)}
		fhex := g.layout.PixelToHex(mouse)
		hex := *fhex.Round()
		if _, ok := (*grid)[hex]; ok {
			(*grid)[hex] = (button == ebiten.MouseButtonLeft)
		}
	}
}

func (g *Game) countLife() {
	count := 0
	for _, value := range g.grid {
		if value {
			count++
		}
	}
	g.lifeCount = count
}

func (g *Game) NextGen() {
	g.generation++
	for hex, value := range g.grid {
		neighbors := 0
		for _, n := range hex.Neighbors() {
			if value, ok := g.grid[n]; ok {
				if value {
					neighbors++
				}
			}
		}
		if value {
			if neighbors != 3 && neighbors != 5 {
				g.grid[hex] = false
			}
		} else {
			if neighbors == 2 {
				g.grid[hex] = true
			}
		}
	}
}

func (g *Game) Stop() {
	g.running = false
	g.generation = 0
	for hex := range g.grid {
		g.grid[hex] = false
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.Stop()
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.running = true
		g.lastTick = time.Now()
	}
	if !g.running {
		g.handleMouseButton(ebiten.MouseButtonLeft, &g.grid)
		g.handleMouseButton(ebiten.MouseButtonRight, &g.grid)
	} else {
		if time.Since(g.lastTick).Seconds() <= 1.0 {
			return nil
		}
		g.lastTick = time.Now()
		g.NextGen()
	}
	g.countLife()
	if g.running && g.lifeCount == 0 {
		g.Stop()
	}
	return nil
}

func (g *Game) DrawInfo(screen *ebiten.Image) {
	var info string
	if !g.running {
		info = fmt.Sprintf("Life count: %d\nLeft click to place cell\nRight click to remove\nS to start\n", g.lifeCount)
	} else {
		info = fmt.Sprintf("Life count: %d\nGeneration: %d\nQ to stop\n", g.lifeCount, g.generation)
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(10, 20)
	op.LineSpacing = mplusNormalFace.Size * 1.5
	text.Draw(screen, info, mplusNormalFace, op)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawInfo(screen)
	for hex := range g.grid {
		corners := g.layout.GetCorners(&hex)
		for i, c := range corners {
			var prev Point
			if i == 0 {
				prev = corners[len(corners)-1]
			} else {
				prev = corners[i-1]
			}
			vector.StrokeLine(screen, c.x, c.y, prev.x, prev.y, 1.0, color.White, false)
		}
		if g.grid[hex] {
			center := g.layout.HexToPixel(&hex)
			vector.DrawFilledCircle(screen, center.x, center.y, 10, colornames.BlueGrey500, false)
		}
	}
}
