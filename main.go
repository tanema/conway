package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/tanema/gluey/term"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	world := newWorld()
	defer world.buf.Done()
	for {
		world.update()
		world.draw()
		time.Sleep(time.Millisecond)
	}
}

type world struct {
	size int
	Data [][]int
	buf  *term.ScreenBuf
}

func newWorld() *world {
	w, h := term.Width(), term.Height()
	size := w
	if h < w {
		size = h
	}

	a := make([][]int, size)
	for i := range a {
		a[i] = make([]int, size)
		for j := range a[i] {
			if rand.Float32() < 0.5 {
				a[i][j] = 1
			}
		}
	}

	return &world{
		size: size,
		Data: a,
		buf:  term.NewScreenBuf(os.Stdout),
	}
}

func (world *world) update() {
	newData := make([][]int, world.size)
	for i := range newData {
		newData[i] = make([]int, world.size)
	}

	for i := range world.Data {
		for j := range world.Data[i] {
			population := world.pop(i, j)
			if population == 3 || population == 2 && world.Data[i][j] == 1 {
				newData[i][j] = 1
			} else {
				newData[i][j] = 0
			}
		}
	}
	world.Data = newData
}

func (world *world) pop(i, j int) int {
	return world.neighbor(i-1, j-1) +
		world.neighbor(i-1, j) +
		world.neighbor(i-1, j+1) +
		world.neighbor(i, j-1) +
		world.neighbor(i, j+1) +
		world.neighbor(i+1, j-1) +
		world.neighbor(i+1, j) +
		world.neighbor(i+1, j+1)
}

func (world *world) neighbor(i, j int) int {
	i = (i + world.size) % world.size
	j = (j + world.size) % world.size
	return world.Data[i][j]
}

func (world *world) draw() {
	world.buf.WriteTmpl(`
{{ range $n, $i := .Data }}{{ range $n, $j := $i }} {{ if eq $j 0 }} {{ else }}o{{end}} {{ end }}
{{ end }}`,
		world,
	)
}
