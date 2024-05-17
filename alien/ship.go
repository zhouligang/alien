package alien

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "golang.org/x/image/bmp"
	"log"
)

type Ship struct {
	image *ebiten.Image
	GameObject
}

func (ship *Ship) Draw(screen *ebiten.Image, cfg *Config) {
	op := &ebiten.DrawImageOptions{}
	//op.GeoM.Translate(float64(cfg.ScreenWidth/2-ship.width)/2, float64(cfg.ScreenHeight/2-ship.height))
	op.GeoM.Translate(ship.x, ship.y)
	screen.DrawImage(ship.image, op)
}

func NewShip(screenWidth, screenHeight int) *Ship {
	img, _, err := ebitenutil.NewImageFromFile("./images/ship.bmp")
	if err != nil {
		log.Fatal(err)
	}
	width, height := img.Size()
	ship := &Ship{
		image: img,
		GameObject: GameObject{
			width:  width,
			height: height,
			x:      float64(screenWidth-width) / 2,
			y:      float64(screenHeight - height),
		},
	}

	return ship
}
