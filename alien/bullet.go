package alien

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Bullet struct {
	image       *ebiten.Image
	speedFactor float64
	GameObject
}

func (bullet *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(bullet.x, bullet.y)
	screen.DrawImage(bullet.image, op)
}

func NewBullet(cfg *Config, ship *Ship) *Bullet {
	rect := image.Rect(0, 0, cfg.BulletWidth, cfg.BulletHeight)
	img := ebiten.NewImageWithOptions(rect, nil)
	img.Fill(cfg.BulletColor)

	return &Bullet{
		image: img,
		GameObject: GameObject{
			width:  cfg.BulletWidth,
			height: cfg.BulletHeight,
			x:      ship.x + float64(ship.width-cfg.BulletWidth)/2,
			y:      float64(cfg.ScreenHeight - ship.height - cfg.BulletHeight),
		},
		speedFactor: cfg.BulletSpeedFactor,
	}
}

func (bullet *Bullet) outOfScreen() bool {
	return bullet.y < -float64(bullet.height)
}
