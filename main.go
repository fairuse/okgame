package main

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/color"
	"log"
	//	"fmt"
	"math"
)

const screenWidth = 1366
const screenHeight = 768

var emptyImage = ebiten.NewImage(3, 3)
var level Level

func init() {
	emptyImage.Fill(color.White)
	level = Level{}
	// polygon := makeNGon(Point{x: screenWidth / 2.0, y: screenHeight / 2.0}, 100.0, 100)
	// level.add(Obstacle{geom: polygon, enabled: true, color: color.RGBA{255, 255, 0, 255}})

	for nr := 0; nr < 10; nr++ {
		polygon := makeNGon(Point{x: screenWidth / 2.0, y: screenHeight / 2.0}, 100-5.0*float64(nr), 100)
		level.add(Obstacle{geom: polygon, enabled: true, color: color.RGBA{R: uint8(255 - nr*15), G: uint8(255 - nr*20), A: 255}})
	}

	for i := 0.0; i < 2*3.1415; i += 0.4 {
		polygon := makeNGon(Point{x: screenWidth*0.5 + screenHeight*0.4*math.Sin(i), y: screenHeight * (0.5 + 0.4*math.Cos(i))}, 15.0, 4)
		level.add(Obstacle{geom: polygon, enabled: true, color: color.RGBA{R: 255, G: 255, A: 255}})
	}
	//fmt.Println("LEVEL",level)
	//v,i := level.draw()
	//fmt.Println(v)
	//fmt.Println(i)
}

func vecLength(x float64, y float64) float64 {
	return math.Sqrt(x*x + y*y)
}

func (g *Game) handleMovement() {
	tx, ty := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.ballx = float64(tx)
		g.bally = float64(ty)
		g.positioning = true
		level.enableAll()
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.targetx = float64(tx)
		g.targety = float64(ty)

		// snapping
		if math.Abs(g.targetx-g.ballx) < 10 {
			g.targetx = g.ballx
		}
		if math.Abs(g.targety-g.bally) < 10 {
			g.targety = g.bally
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		l := vecLength(g.targetx-g.ballx, g.targety-g.bally)
		g.speedx = (g.targetx - g.ballx) / l * 5
		g.speedy = (g.targety - g.bally) / l * 5
		if l == 0 {
			g.speedx = 0
			g.speedy = 0
		}
		//		g.ballx = float64(tx)
		//		g.bally = float64(ty)
		g.positioning = false
	}

	if !g.positioning {
		// fmt.Println(g.ballx, g.bally, g.speedx, g.speedy)
		oldPos := Point{g.ballx, g.bally}

		g.ballx += g.speedx
		g.bally += g.speedy

		newPos := Point{g.ballx, g.bally}

		interp, obj := level.findHit(oldPos, newPos)

		//interp := intersectPolygonNorm(Point{x: g.ballx, y: g.bally}, Point{x: g.targetx, y: g.targety},
		//	polygon) // Point{x:screenWidth/2.0, y:0}, Point{x:screenWidth/2.0,y:screenHeight} )
		if interp != nil {
			obj.enabled = false
			// traveled := interp.src.sub(oldPos).length() // distance traveled before hit
			newdir := interp.src.sub(oldPos).reflect(interp.dst).normalized().mul(5.0)

			// newtarget := interp.src.add(newdir.normalized().mul(screenWidth+screenHeight))
			newball := interp.src.add(newdir.mul(epsilon))

			g.ballx = newball.x
			g.bally = newball.y
			g.speedx = newdir.x
			g.speedy = newdir.y
		}

		if g.ballx > screenWidth && g.speedx > 0 {
			g.speedx *= -1
		}
		if g.bally > screenHeight && g.speedy > 0 {
			g.speedy *= -1
		}

		if g.ballx < 0 && g.speedx < 0 {
			g.speedx *= -1
		}
		if g.bally < 0 && g.speedy < 0 {
			g.speedy *= -1
		}
	}
}

type Game struct {
	showRays         bool
	ballx, bally     float64
	targetx, targety float64
	speedx, speedy   float64
	positioning      bool
	trail            []Point
	//        objects  []object
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("game ended by player")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.showRays = !g.showRays
	}

	g.handleMovement()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0x30, G: 0x30, B: 0x50, A: 0xff})
	src := emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

	if true {
		v, i := level.draw()
		screen.DrawTriangles(v, i, src, nil)
	}

	g.trail = append(g.trail, Point{g.ballx, g.bally})
	const maxtrail = 100
	if len(g.trail) > maxtrail {
		g.trail = g.trail[len(g.trail)-maxtrail:]
	}
	for nr, trailpos := range g.trail {
		polygon := makeNGon(trailpos, 5.0*float64(nr)/float64(maxtrail+1), 4)
		v, i := render(polygon, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: uint8(0xa0 * nr / maxtrail)})
		screen.DrawTriangles(v, i, src, nil)
	}

	if g.positioning {
		v, i := dotline(g.ballx, g.bally, g.targetx, g.targety, color.RGBA{R: 0x80, G: 0x40, B: 0xa0, A: 0xff})
		screen.DrawTriangles(v, i, src, nil)

		// polygon := makeNGon(Point{x: screenWidth / 2.0, y: screenHeight / 2.0}, 50.0, 5)

		// v, i = render(polygon, color.RGBA{0xff, 0xff, 0xff, 0xff})
		// screen.DrawTriangles(v, i, src, nil)

		interp, _ := level.findHit(Point{x: g.ballx, y: g.bally}, Point{x: g.targetx, y: g.targety})

		//interp := intersectPolygonNorm(Point{x: g.ballx, y: g.bally}, Point{x: g.targetx, y: g.targety},
		//	polygon) // Point{x:screenWidth/2.0, y:0}, Point{x:screenWidth/2.0,y:screenHeight} )
		if interp != nil {
			v, i := mkBallGeom(interp.src.x, interp.src.y, color.RGBA{R: 0xff, B: 0xff, A: 0xff})
			screen.DrawTriangles(v, i, src, nil)

			v, i = dotline(interp.src.x, interp.src.y, interp.src.x+15.0*interp.dst.x, interp.src.y+15.0*interp.dst.y, color.RGBA{G: 0x40, B: 0xa0, A: 0xff})
			screen.DrawTriangles(v, i, src, nil)

			newdir := interp.src.sub(Point{g.ballx, g.bally}).reflect(interp.dst)
			v, i = dotline(interp.src.x, interp.src.y, interp.src.x+newdir.x, interp.src.y+newdir.y, color.RGBA{G: 0x80, B: 0xa0, A: 0xff})
			screen.DrawTriangles(v, i, src, nil)

			newtarget := interp.src.add(newdir.normalized().mul(screenWidth + screenHeight))
			newball := interp.src.add(newdir.mul(epsilon))

			interp, _ := level.findHit(newball, newtarget)
			if interp != nil {
				newdir := interp.src.sub(newball).reflect(interp.dst)
				v, i = dotline(interp.src.x, interp.src.y, interp.src.x+newdir.x, interp.src.y+newdir.y, color.RGBA{R: 0xff, A: 0xff})
				screen.DrawTriangles(v, i, src, nil)
			}
		}

		v, i = mkBallGeom(g.targetx, g.targety, color.RGBA{R: 0xa0, G: 0xa0, B: 0xa0, A: 0xff})
		screen.DrawTriangles(v, i, src, nil)

	}

	v, i := mkBallGeom(g.ballx, g.bally, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
	screen.DrawTriangles(v, i, src, nil)

	ebitenutil.DebugPrintAt(screen, "Click and drag to launch ball", 0, 0)
	//ebitenutil.DebugPrintAt(screen, "Ball", g.ballx, g.bally)
	//ebitenutil.DebugPrintAt(screen, "Target", g.targetx, g.targety)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		ballx: screenWidth / 2,
		bally: screenHeight / 2,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Ray casting and shadows (Ebiten demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
