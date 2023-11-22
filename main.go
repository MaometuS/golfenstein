package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"math"
)

type Walls struct {
	LimX  int
	LimY  int
	Size  int
	Walls [][]uint8
}

func NewWalls() *Walls {
	return &Walls{
		LimX: 8,
		LimY: 8,
		Size: 64,
		Walls: [][]uint8{
			{1, 1, 1, 1, 1, 1, 1, 1},
			{1, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 2, 0, 3, 3, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 1},
			{1, 0, 0, 0, 0, 4, 0, 1},
			{1, 0, 5, 0, 0, 4, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 1},
			{1, 1, 1, 1, 1, 1, 1, 1},
		},
	}
}

func (w *Walls) Draw(screen *ebiten.Image) {
	for y := 0; y < w.LimY; y++ {
		for x := 0; x < w.LimX; x++ {
			if w.Walls[y][x] == 1 {
				vector.DrawFilledRect(screen, float32(x*w.Size)+1, float32(y*w.Size)+1, float32(w.Size)-1, float32(w.Size)-1, color.RGBA{
					R: 255,
					G: 255,
					B: 255,
					A: 255,
				}, false)
			} else {
				vector.DrawFilledRect(screen, float32(x*w.Size)+1, float32(y*w.Size)+1, float32(w.Size)-1, float32(w.Size)-1, color.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 255,
				}, false)
			}
		}
	}
}

type Player struct {
	X      float32
	Y      float32
	DeltaX float32
	DeltaY float32
	MouseX int
	Angle  float32
}

func NewPlayer() *Player {
	return &Player{
		X:      89,
		Y:      87,
		DeltaX: float32(math.Cos(0)) * 5,
		DeltaY: float32(math.Sin(0)) * 5,
		Angle:  0,
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, p.X, p.Y, 10, color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	}, false)

	vector.StrokeLine(screen, p.X, p.Y, p.X+p.DeltaX*5, p.Y+p.DeltaY*5, 2, color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	}, false)
}

func checkCollision(rx, ry float32, walls [][]uint8) bool {
	mx := int(rx) >> 6
	my := int(ry) >> 6

	return mx >= 0 && my >= 0 && mx < 8 && my < 8 && walls[my][mx] > 0
}

func (p *Player) HandleMovement(walls [][]uint8) {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		panic("stop")
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		if !checkCollision(p.X+p.DeltaX*2, p.Y+p.DeltaY*2, walls) {
			p.X += p.DeltaX
			p.Y += p.DeltaY
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		if !checkCollision(p.X-p.DeltaX*2, p.Y-p.DeltaY*2, walls) {
			p.X -= p.DeltaX
			p.Y -= p.DeltaY
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		if !checkCollision(p.X+p.DeltaY*2, p.Y-p.DeltaX*2, walls) {
			p.X += p.DeltaY
			p.Y -= p.DeltaX
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		if !checkCollision(p.X-p.DeltaY*2, p.Y+p.DeltaX*2, walls) {
			p.X -= p.DeltaY
			p.Y += p.DeltaX
		}
	}

	ebiten.SetCursorMode(ebiten.CursorModeCaptured)

	curX, _ := ebiten.CursorPosition()
	curDX := p.MouseX - curX
	p.MouseX = curX

	da := float32(curDX) / float32(1024) * 4

	p.Angle -= da
	if p.Angle < 0 {
		p.Angle += math.Pi * 2
	}

	if p.Angle > math.Pi*2 {
		p.Angle -= math.Pi * 2
	}

	p.DeltaX = float32(math.Cos(float64(p.Angle))) * 3
	p.DeltaY = float32(math.Sin(float64(p.Angle))) * 3
}

type Game struct {
	player *Player
	walls  *Walls
}

func (g *Game) Update() error {
	g.player.HandleMovement(g.walls.Walls)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{
		R: 50,
		G: 50,
		B: 50,
		A: 255,
	})

	g.castRays(screen)

	ebitenutil.DebugPrintAt(screen, "Press 'Esc' to exit", 0, 0)
	ebitenutil.DebugPrintAt(screen, "Use 'WASD' for movement", 0, 15)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 512
}

func dist(ax, ay, bx, by float64) float64 {
	return math.Sqrt(math.Pow(ax-bx, 2) + math.Pow(ay-by, 2))
}

func checkHorizontal(angle float64, px float64, py float64, walls [][]uint8) (uint8, float64) {
	aTan := -1 / math.Tan(angle)

	var rx, ry, xo, yo, disH, hx, hy float64
	var mx, my, dof int
	var col uint8

	disH = math.MaxFloat64

	if angle > math.Pi {
		ry = float64((int(py)>>6)<<6) - 0.0001
		rx = (py-ry)*aTan + px
		yo = float64(-64)
		xo = -yo * aTan
	}

	if angle < math.Pi {
		ry = float64((int(py)>>6)<<6) + 64
		rx = (py-ry)*aTan + px
		yo = float64(64)
		xo = -yo * aTan
	}

	if angle == 0 || angle == math.Pi {
		rx = px
		ry = py
		dof = 8
	}

	for dof < 8 {
		mx = int(rx) >> 6
		my = int(ry) >> 6

		if mx >= 0 && my >= 0 && mx < 8 && my < 8 && walls[my][mx] > 0 {
			hx = rx
			hy = ry
			col = walls[my][mx]
			disH = dist(hx, hy, px, py)
			dof = 8
		} else {
			rx += xo
			ry += yo

			dof++
		}
	}

	return col, disH
}

func checkVertical(angle float64, px float64, py float64, walls [][]uint8) (uint8, float64) {
	var rx, ry, xo, yo, disV, vx, vy float64
	var mx, my, dof int
	var col uint8

	disV = math.MaxFloat64

	dof = 0
	nTan := -math.Tan(angle)

	pi2 := math.Pi / 2
	pi3 := 3 * math.Pi / 2

	if angle > pi2 && angle < pi3 {
		rx = float64((int(px)>>6)<<6) - 0.0001
		ry = (px-rx)*nTan + py
		xo = float64(-64)
		yo = -xo * nTan
	}

	if angle < pi2 || angle > pi3 {
		rx = float64((int(px)>>6)<<6) + 64
		ry = (px-rx)*nTan + py
		xo = float64(64)
		yo = -xo * nTan
	}

	if angle == pi2 || angle == pi3 {
		rx = px
		ry = py
		dof = 8
	}

	for dof < 8 {
		mx = int(rx) >> 6
		my = int(ry) >> 6

		if mx >= 0 && my >= 0 && mx < 8 && my < 8 && walls[my][mx] > 0 {
			vx = rx
			vy = ry
			col = walls[my][mx]
			disV = dist(vx, vy, px, py)
			dof = 8
		} else {
			rx += xo
			ry += yo

			dof++
		}
	}

	return col, disV
}

func (p *Game) castRays(screen *ebiten.Image) {
	for ray := 0; ray < 1024; ray++ {
		delta := math.Pi / 4 / 1024
		rayAngle := float64(p.player.Angle) - 512*delta + float64(ray)*delta
		if rayAngle < 0 {
			rayAngle += 2 * math.Pi
		}
		if rayAngle > 2*math.Pi {
			rayAngle -= 2 * math.Pi
		}

		colH, disH := checkHorizontal(rayAngle, float64(p.player.X), float64(p.player.Y), p.walls.Walls)
		colV, disV := checkVertical(rayAngle, float64(p.player.X), float64(p.player.Y), p.walls.Walls)

		var disT float64
		var cl color.Color

		if disH > disV {
			disT = disV

			switch colV {
			case 1:
				cl = color.RGBA{
					R: 200,
					G: 200,
					B: 0,
					A: 255,
				}
			case 2:
				cl = color.RGBA{
					R: 200,
					G: 0,
					B: 0,
					A: 255,
				}
			case 3:
				cl = color.RGBA{
					R: 0,
					G: 200,
					B: 0,
					A: 255,
				}
			case 4:
				cl = color.RGBA{
					R: 0,
					G: 0,
					B: 200,
					A: 255,
				}
			case 5:
				cl = color.RGBA{
					R: 200,
					G: 0,
					B: 200,
					A: 255,
				}
			}
		} else {
			disT = disH

			switch colH {
			case 1:
				cl = color.RGBA{
					R: 150,
					G: 150,
					B: 0,
					A: 255,
				}
			case 2:
				cl = color.RGBA{
					R: 150,
					G: 0,
					B: 0,
					A: 255,
				}
			case 3:
				cl = color.RGBA{
					R: 0,
					G: 150,
					B: 0,
					A: 255,
				}
			case 4:
				cl = color.RGBA{
					R: 0,
					G: 0,
					B: 150,
					A: 255,
				}
			case 5:
				cl = color.RGBA{
					R: 150,
					G: 0,
					B: 150,
					A: 255,
				}
			}
		}

		ca := rayAngle - float64(p.player.Angle)
		if ca < 0 {
			ca += 2 * math.Pi
		}

		if ca > 2*math.Pi {
			ca -= 2 * math.Pi
		}

		disT = disT * math.Cos(ca)

		lineH := float64(p.walls.Size*512) / disT
		if lineH > 512 {
			lineH = 512
		}
		lineO := 256 - lineH/2
		width := float32(1)
		vector.StrokeLine(screen, float32(ray)*width, float32(lineO), float32(ray)*width, float32(lineH+lineO), width, cl, false)
	}
}

func main() {
	ebiten.SetWindowSize(1024, 512)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{
		player: NewPlayer(),
		walls:  NewWalls(),
	}); err != nil {
		log.Fatal(err)
	}
}
