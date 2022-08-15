package main

import (
	"log"
	"errors"
        "github.com/hajimehoshi/ebiten/v2"
        "github.com/hajimehoshi/ebiten/v2/ebitenutil"
//        "github.com/hajimehoshi/ebiten/v2/examples/resources/images"
        "github.com/hajimehoshi/ebiten/v2/inpututil"
)

const screenWidth = 256
const screenHeight = 256

func (g *Game) handleMovement() {
        tx, ty := ebiten.CursorPosition()
        if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && tx != 0 && ty != 0 {
                g.px = tx
                g.py = ty
        }
}

type Game struct {
        showRays bool
        px, py   int
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
        ebitenutil.DebugPrintAt(screen, "Click and drag to launch ball", 0, 0)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
        return screenWidth, screenHeight
}


func main() {
        g := &Game{
                px: screenWidth / 2,
                py: screenHeight / 2,
        }


        ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
        ebiten.SetWindowTitle("Ray casting and shadows (Ebiten demo)")
        if err := ebiten.RunGame(g); err != nil {
                log.Fatal(err)
        }
}
 
