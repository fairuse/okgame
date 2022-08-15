package main

import (
	"log"
	"errors"
        "image/color"
	"image"
        "github.com/hajimehoshi/ebiten/v2"
        "github.com/hajimehoshi/ebiten/v2/ebitenutil"
        "github.com/hajimehoshi/ebiten/v2/inpututil"
)

const screenWidth = 512
const screenHeight = 512

var  emptyImage = ebiten.NewImage(3, 3)

func init() {
	emptyImage.Fill( color.White )
}

func (g *Game) handleMovement() {
        tx, ty := ebiten.CursorPosition()
        if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
                g.ballx = tx
                g.bally = ty
        }
        if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
                g.targetx = tx
                g.targety = ty
        }
        if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
                g.ballx = tx
                g.bally = ty
        }
}

type Game struct {
        showRays bool
        ballx, bally   int
        targetx, targety   int
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
        src := emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

        ebitenutil.DebugPrintAt(screen, "Click and drag to launch ball", 0, 0)
	ebitenutil.DebugPrintAt(screen, "Ball", g.ballx, g.bally)
	ebitenutil.DebugPrintAt(screen, "Target", g.targetx, g.targety)

        v,i := ball(float64(g.ballx),float64(g.bally),color.RGBA{0xff,0xff,0,0xff})
        screen.DrawTriangles(v, i, src, nil)

        v,i = ball(float64(g.targetx),float64(g.targety),color.RGBA{0xa0,0xa0,0xa0,0xff})
        screen.DrawTriangles(v, i, src, nil)
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
        ebiten.SetWindowTitle("Ray casting and shadows (Ebiten demo)")
        if err := ebiten.RunGame(g); err != nil {
                log.Fatal(err)
        }
}
 
