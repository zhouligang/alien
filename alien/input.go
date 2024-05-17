package alien

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"time"
)

type Input struct {
	Msg            string
	lastBulletTime time.Time
}

func (i *Input) Update(game *Game) {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		//fmt.Println("←←←←←←←←←←←←←←←←←←←←←←←")
		i.Msg = "left pressed"
		game.ship.x -= game.cfg.ShipSpeedFactor
		if game.ship.x < -float64(game.ship.width) {
			game.ship.x = -float64(game.ship.width)
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		//fmt.Println("→→→→→→→→→→→→→→→→→→→→→→→")
		i.Msg = "right pressed"
		game.ship.x += game.cfg.ShipSpeedFactor
		if game.ship.x > float64(game.cfg.ScreenWidth)-float64(game.ship.width) {
			game.ship.x = float64(game.cfg.ScreenWidth) - float64(game.ship.width)
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		//fmt.Println("-----------------------")
		i.Msg = "space pressed"
		if len(game.bullets) >= game.cfg.MaxBulletNum || time.Now().Sub(i.lastBulletTime).Milliseconds() < game.cfg.BulletInterval {
			return
		}
		bullet := NewBullet(game.cfg, game.ship)
		game.addBullet(bullet)
		i.lastBulletTime = time.Now()
	}
}

func (i *Input) IsKeyPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
}
