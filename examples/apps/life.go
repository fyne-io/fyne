package apps

import "image/color"
import "time"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/theme"

func aliveColor() color.RGBA {
	return theme.TextColor()
}

func deadColor() color.RGBA {
	return theme.BackgroundColor()
}

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

func (b *board) renderer(x, y, w, h int) color.RGBA {
	xpos := int(float64(b.width) * (float64(x) / float64(w)))
	ypos := int(float64(b.height) * (float64(y) / float64(h)))

	if xpos >= b.width || ypos >= b.height {
		return deadColor()
	}
	if b.cells[ypos][xpos] {
		return aliveColor()
	}

	return deadColor()
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

type boardLayout struct {
	board  *board
	render *canvas.Image
}

// Layout sets the Life render raster to be full size
func (l *boardLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	for _, obj := range objs {
		obj.Resize(size)
	}
}

// MinSize returns a suitable size to fit every cell at 10x10
func (l *boardLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(l.board.width*10, l.board.height*10)
}

func (l *boardLayout) animate() {
	go func() {
		tick := time.NewTicker(time.Second / 6)
		canvas := fyne.GetCanvas(l.render)

		for {
			select {
			case <-tick.C:
				state := l.board.nextGen()
				l.board.renderState(state)
				canvas.Refresh(l.render)
			}
		}
	}()
}

func newBoardLayout(b *board) *boardLayout {
	render := canvas.NewRaster(b.renderer)
	return &boardLayout{b, render}
}

// Life starts a new game of life
func Life(app fyne.App) {
	board := newBoard()
	board.load()

	layout := newBoardLayout(board)

	window := app.NewWindow("Life")
	window.SetContent(fyne.NewContainerWithLayout(layout,
		layout.render))

	// start the board animation before we show the window - it will block
	layout.animate()

	window.Show()
}
