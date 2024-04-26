package goh

import (
	"fmt"
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

const (
	ScreenWidth  = 1024
	ScreenHeight = 1024
)

const (
	ActionNone = iota
	ActionDraw
	ActionErase
)

var (
	Image = ebiten.NewImage(3, 3)

	// subImage is an internal sub image of Image.
	// Use bImage at DrawTriangles instead of mage in order to avoid bleeding edges.
	subImage = Image.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	Image.Fill(colornames.AmberA700)
}

type Cell struct {
	neighbors int
	alive     bool
}

func (c *Cell) AliveNextGen() bool {
	if c.neighbors == 3 {
		return true
	} else if c.neighbors != 2 {
		return false
	}
	return c.alive
}

func NewBoard() map[Hex]*Cell {
	b := make(map[Hex]*Cell)
	N := 10
	for q := -N; q <= N; q++ {
		r1 := Max(-N, -q-N)
		r2 := Min(N, -q+N)
		for r := r1; r <= r2; r++ {
			b[Hex{q, r, -q - r}] = &Cell{0, false}
		}
	}
	return b
}

type Game struct {
	layout     Layout
	board      map[Hex]*Cell
	generation int
	running    bool
	lastTick   time.Time
	lifeCount  int
	action     int
}

func NewGame() *Game {
	g := &Game{
		layout:     *NewLayout(&PointyTop, Point{25, 25}, Point{500, 500}),
		board:      NewBoard(),
		generation: 0,
		running:    false,
		lastTick:   time.Now(),
		lifeCount:  0,
		action:     ActionNone,
	}
	return g
}

func (g *Game) GetCell(hex Hex) *Cell {
	if cell, ok := g.board[hex]; ok {
		return cell
	}
	return nil
}

func (g *Game) GetNeighbors(hex Hex) []*Cell {
	var neighbors []*Cell
	for _, dir := range hex_directions {
		neighbor := g.GetCell(AddHexes(hex, dir))
		if neighbor == nil {
			continue
		}
		neighbors = append(neighbors, neighbor)

	}
	return neighbors
}

func (g *Game) UpdateBoardState() {
	g.lifeCount = 0
	for hex, cell := range g.board {
		if cell.alive {
			g.lifeCount++
		}

		count := 0
		for _, cell := range g.GetNeighbors(hex) {
			if cell.alive {
				count++
			}
		}
		cell.neighbors = count
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) NextGen() {
	changed := false

	for _, cell := range g.board {
		old_state := cell.alive
		if cell.neighbors == 3 {
			cell.alive = true
		} else if cell.neighbors != 2 {
			cell.alive = false
		}
		if cell.alive != old_state {
			changed = true
		}
	}

	g.UpdateBoardState()
	g.generation++

	if !changed {
		g.Stop()
	}
}

func (g *Game) Stop() {
	g.running = false
	g.generation = 0
	for hex := range g.board {
		g.board[hex].alive = false
		g.board[hex].neighbors = 0
	}
}

func (g *Game) Start() {
	if !g.running {
		g.generation = 0
		g.running = true
		g.lastTick = time.Now()
	}
}

func (g *Game) GetHexUnderCursor() Hex {
	x, y := ebiten.CursorPosition()
	mouse := Point{float32(x), float32(y)}
	return g.layout.PixelToHex(mouse)
}

func (g *Game) GetCellUnderCursor() *Cell {
	hex := g.GetHexUnderCursor()
	return g.GetCell(hex)
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.Stop()
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Start()
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.action = ActionNone
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cell := g.GetCellUnderCursor()
		if cell != nil && cell.alive {
			g.action = ActionErase
		} else {
			g.action = ActionDraw
		}
	}

	if !g.running {
		if g.action == ActionDraw || g.action == ActionErase {
			cell := g.GetCellUnderCursor()
			if cell != nil {
				if g.action == ActionDraw {
					cell.alive = true
				} else {
					cell.alive = false
				}
			}
		}
		g.UpdateBoardState()

	} else {
		if time.Since(g.lastTick).Seconds() <= 0.5 {
			return nil
		}
		g.lastTick = time.Now()
		g.NextGen()

		if g.lifeCount == 0 {
			g.Stop()
		}
	}

	return nil
}

func (g *Game) GetCellColor(cell *Cell) float32 {
	if cell.alive {
		return 1.0
	}
	return 0.1
}

func (g *Game) DrawInfo(screen *ebiten.Image) {
	var info string
	if !g.running {
		info = fmt.Sprintf("Life count: %d\nLeft click to place or remove cells\nS to start\n", g.lifeCount)
	} else {
		info = fmt.Sprintf("Life count: %d\nGeneration: %d\nQ to stop\n", g.lifeCount, g.generation)
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(10, 20)
	op.LineSpacing = mplusNormalFace.Size * 1.5
	text.Draw(screen, info, mplusNormalFace, op)
}

func (g *Game) DrawHex(screen *ebiten.Image, hex Hex, cell *Cell) {
	corners := g.layout.GetCorners(&hex)
	center := g.layout.HexToPixel(&hex)

	var alpha float32
	alpha = 0.1
	if cell.alive {
		alpha = 1.0
	}

	for i, c := range corners {
		path := vector.Path{}
		path.MoveTo(float32(center.x), float32(center.y))
		path.LineTo(float32(c.x), float32(c.y))
		path.LineTo(float32(corners[(i+1)%6].x), float32(corners[(i+1)%6].y))
		path.Close()

		var vs []ebiten.Vertex
		var is []uint16
		vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)

		op := &ebiten.DrawTrianglesOptions{}
		op.FillRule = ebiten.EvenOdd

		for i := range vs {
			vs[i].ColorA = alpha
		}
		screen.DrawTriangles(vs, is, subImage, op)
	}

	for i, c := range corners {
		var prev Point
		if i == 0 {
			prev = corners[len(corners)-1]
		} else {
			prev = corners[i-1]
		}
		vector.StrokeLine(screen, c.x, c.y, prev.x, prev.y, 2.0, colornames.BlueGrey500, false)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawInfo(screen)
	for hex, cell := range g.board {
		g.DrawHex(screen, hex, cell)
	}
}
