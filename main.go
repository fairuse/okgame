package main

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/color"
	"log"
	"math"
)

const screenWidth = 768
const screenHeight = 512

var emptyImage = ebiten.NewImage(3, 3)

func init() {
	emptyImage.Fill(color.White)
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
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.targetx = float64(tx)
		g.targety = float64(ty)

		// snapping
		if math.Abs(g.targetx - g.ballx) < 10 {
			g.targetx = g.ballx
		}
		if math.Abs(g.targety - g.bally) < 10 {
			g.targety = g.bally
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		l := vecLength(g.targetx-g.ballx, g.targety-g.bally)
		g.speedx = (g.targetx - g.ballx) / l * 5
		g.speedy = (g.targety - g.bally) / l * 5
		//		g.ballx = float64(tx)
		//		g.bally = float64(ty)
		g.positioning = false
	}

	if !g.positioning {
		// fmt.Println(g.ballx, g.bally, g.speedx, g.speedy)
		g.ballx += g.speedx
		g.bally += g.speedy

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
	screen.Fill( color.RGBA{0x30,0x30,0x50,0xff})
	src := emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

	if g.positioning {
		v, i := dotline(float64(g.ballx), float64(g.bally), float64(g.targetx), float64(g.targety), color.RGBA{0x80, 0x40, 0xa0, 0xff})
		screen.DrawTriangles(v, i, src, nil)

		polygon := makeNGon( Point{x:screenWidth/2.0, y:screenHeight/2.0 }, 50.0, 5 )

		v, i = render(polygon, color.RGBA{0xff,0xff,0xff,0xff})
		screen.DrawTriangles(v, i, src, nil)

		interp := intersectPolygonNorm(Point{x:g.ballx, y:g.bally}, Point{x:g.targetx,y:g.targety},
				      polygon )  // Point{x:screenWidth/2.0, y:0}, Point{x:screenWidth/2.0,y:screenHeight} )
		if interp != nil {
			v, i := ball(interp.src.x, interp.src.y, color.RGBA{0xff,0x00,0xff,0xff})
			screen.DrawTriangles(v, i, src, nil)

			v, i = dotline(float64(interp.src.x), float64(interp.src.y), float64(interp.src.x+15.0*interp.dst.x), float64(interp.src.y+15.0*interp.dst.y), color.RGBA{0x0, 0x40, 0xa0, 0xff})
			screen.DrawTriangles(v, i, src, nil)

			newdir := interp.src.sub(Point{g.ballx,g.bally}).reflect(interp.dst)
			v, i = dotline(float64(interp.src.x), float64(interp.src.y), float64(interp.src.x+newdir.x), float64(interp.src.y+newdir.y), color.RGBA{0x0, 0x80, 0xa0, 0xff})
			screen.DrawTriangles(v, i, src, nil)
		}

		v, i = ball(float64(g.targetx), float64(g.targety), color.RGBA{0xa0, 0xa0, 0xa0, 0xff})
		screen.DrawTriangles(v, i, src, nil)
	}

	v, i := ball(float64(g.ballx), float64(g.bally), color.RGBA{0xff, 0xff, 0, 0xff})
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
