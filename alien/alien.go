package alien

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

type Alien struct {
	image *ebiten.Image
	GameObject
	speedFactor float64
}

func NewAlien(cfg *Config) *Alien {
	img, _, err := ebitenutil.NewImageFromFile("./images/alien.png")
	if err != nil {
		log.Fatal(err)
	}

	width, height := img.Size()
	return &Alien{
		image: img,
		GameObject: GameObject{
			width:  width,
			height: height,
			x:      0,
			y:      0,
		},
		speedFactor: cfg.AlienSpeedFactor,
	}
}

func (alien *Alien) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(alien.x, alien.y)
	screen.DrawImage(alien.image, op)
}

func (alien *Alien) outOfScreen(config *Config) bool {
	return alien.y > float64(config.ScreenHeight)-float64(alien.height)
}
