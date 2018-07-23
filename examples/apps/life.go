package apps

import "image/color"
import "time"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/theme"

type board struct {
	cells  [][]bool
	width  int
	height int
}

func (b *board) countNeighbours(x, y int) int {
	sum := 0
	if y > 0 {
		if x > 0 && b.cells[y-1][x-1] {
			sum++
		}
		if b.cells[y-1][x] {
			sum++
		}
		if x < b.width-1 && b.cells[y-1][x+1] {
			sum++
		}
	}

	if x > 0 && b.cells[y][x-1] {
		sum++
	}
	if x < b.width-1 && b.cells[y][x+1] {
		sum++
	}

	if y < b.height-1 {
		if x > 0 && b.cells[y+1][x-1] {
			sum++
		}
		if b.cells[y+1][x] {
			sum++
		}
		if x < b.width-1 && b.cells[y+1][x+1] {
			sum++
		}
	}

	return sum
}

func (b *board) nextGen() [][]bool {
	state := make([][]bool, b.height)

	for y := 0; y < b.height; y++ {
		state[y] = make([]bool, b.width)

		for x := 0; x < b.width; x++ {
			n := b.countNeighbours(x, y)

			if b.cells[y][x] {
				state[y][x] = n == 2 || n == 3
			} else {
				state[y][x] = n == 3
			}
		}
	}

	return state
}

func (b *board) renderState(state [][]bool) {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			b.cells[y][x] = state[y][x]
		}
	}
}

func (b *board) load() {
	b.cells[1][16] = true
	b.cells[2][17] = true
	b.cells[2][18] = true
	b.cells[3][16] = true
	b.cells[3][17] = true

	b.cells[15][22] = true
	b.cells[15][23] = true
	b.cells[15][24] = true

	b.cells[15][32] = true
	b.cells[15][33] = true
	b.cells[16][32] = true
	b.cells[16][33] = true
}

func newBoard() *board {
	b := &board{nil, 50, 40}
	b.cells = make([][]bool, b.height)

	for y := 0; y < b.height; y++ {
		b.cells[y] = make([]bool, b.width)
	}

	return b
}

type game struct {
	board  *board
	paused bool

	aliveColor color.RGBA
	deadColor  color.RGBA

	size     fyne.Size
	position fyne.Position
	render   *canvas.Image
	objects  []fyne.CanvasObject
}

func (g *game) CurrentSize() fyne.Size {
	return g.size
}

func (g *game) Resize(size fyne.Size) {
	g.size = size
	g.render.Resize(size)
}

func (g *game) CurrentPosition() fyne.Position {
	return g.position
}

func (g *game) Move(pos fyne.Position) {
	g.position = pos
	g.render.Move(pos)
}

func (g *game) MinSize() fyne.Size {
	return fyne.NewSize(g.board.width*10, g.board.height*10)
}

func (g *game) ApplyTheme() {
	g.aliveColor = theme.TextColor()
	g.deadColor = theme.BackgroundColor()
}

func (g *game) CanvasObjects() []fyne.CanvasObject {
	return g.objects
}

func (g *game) cellForCoord(x, y, w, h int) (int, int) {
	xpos := int(float64(g.board.width) * (float64(x) / float64(w)))
	ypos := int(float64(g.board.height) * (float64(y) / float64(h)))

	return xpos, ypos
}

func (g *game) renderer(x, y, w, h int) color.RGBA {
	xpos, ypos := g.cellForCoord(x, y, w, h)

	if xpos >= g.board.width || ypos >= g.board.height {
		return g.deadColor
	}
	if g.board.cells[ypos][xpos] {
		return g.aliveColor
	}

	return g.deadColor
}

func (g *game) run() {
	g.paused = false
}

func (g *game) stop() {
	g.paused = true
}

func (g *game) toggleRun() {
	g.paused = !g.paused
}

func (g *game) animate() {
	go func() {
		tick := time.NewTicker(time.Second / 6)
		canvas := fyne.GetCanvas(g.render)

		for {
			select {
			case <-tick.C:
				if g.paused {
					continue
				}

				state := g.board.nextGen()
				g.board.renderState(state)
				canvas.Refresh(g.render)
			}
		}
	}()
}

func (g *game) keyDown(ev *fyne.KeyEvent) {
	if ev.Name == "space" {
		g.toggleRun()
	}
}

func (g *game) OnMouseDown(ev *fyne.MouseEvent) {
	xpos, ypos := g.cellForCoord(ev.Position.X, ev.Position.Y, g.size.Width, g.size.Height)

	if xpos >= g.board.width || ypos >= g.board.height {
		return
	}

	g.board.cells[ypos][xpos] = !g.board.cells[ypos][xpos]

	canvas := fyne.GetCanvas(g.render)
	canvas.Refresh(g.render)
}

func newGame(b *board) *game {
	g := &game{board: b}
	render := canvas.NewRaster(g.renderer)
	g.ApplyTheme()

	g.render = render
	g.objects = []fyne.CanvasObject{render}
	return g
}

// Life starts a new game of life
func Life(app fyne.App) {
	board := newBoard()
	board.load()

	game := newGame(board)

	window := app.NewWindow("Life")
	window.SetContent(game)
	window.Canvas().SetOnKeyDown(game.keyDown)

	// start the board animation before we show the window - it will block
	game.animate()

	window.Show()
}
