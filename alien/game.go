package alien

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/jinzhu/copier"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeOver
)

var mutex sync.RWMutex

type Game struct {
	input     *Input
	cfg       *Config
	ship      *Ship
	bullets   map[*Bullet]struct{}
	aliens    map[*Alien]struct{}
	mode      Mode
	failCount int // 被外星人碰撞和移出屏幕的外星人数量之和
	overMsg   string
	score     int
}

func NewGame() *Game {
	config := LoadConfig()
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenWidth)
	ebiten.SetWindowTitle("yes it is!")
	g := &Game{
		input:   &Input{Msg: "Hello, World!"},
		cfg:     config,
		ship:    NewShip(config.ScreenWidth, config.ScreenHeight),
		bullets: make(map[*Bullet]struct{}),
		//aliens:  make(map[*Alien]struct{}),
		aliens: make(map[*Alien]struct{}),
	}
	g.init()
	return g
}

func (g *Game) Update() error {

	//var alienA Alien
	//alienA.x = 220
	//alienA.y = 344
	//alienA.width = 60
	//alienA.height = 58
	//var alienB Alien
	//alienB.x = 180
	//alienB.y = 344
	//alienB.width = 60
	//alienB.height = 58
	//log.Printf("result is %v", CheckCollision(&alienA, &alienB))
	//return errors.New("error happened")
	switch g.mode {
	case ModeTitle:
		if g.input.IsKeyPressed() {
			g.mode = ModeGame
		}
	case ModeGame:
		if g.input.IsKeyPressed() {
			g.mode = ModeTitle
		}
		for bullet := range g.bullets {
			// 向上移动
			bullet.y -= bullet.speedFactor
			// 出屏幕了则删除
			if bullet.outOfScreen() {
				delete(g.bullets, bullet)
			}
		}

		for alien := range g.aliens {
			// 向下移动
			alien.y += alien.speedFactor

			if alien.outOfScreen(g.cfg) {
				//delete(g.aliens, alien)
				// 回到初始位置
				alien.y = 0
				// 掉出屏幕了，算失败一次
				g.failCount++
			}
		}
		// 随机横向移动
		moveAlienX(g)

		if g.failCount > 0 {
			g.overMsg = "Game Over!"
		} else if len(g.aliens) == 0 {
			g.overMsg = "You Win!"
		}

		if len(g.overMsg) > 0 {
			g.mode = ModeOver
			// 清除剩余的子弹和外星人
			g.aliens = make(map[*Alien]struct{})
			g.bullets = make(map[*Bullet]struct{})
		}

		// 检查碰撞
		g.CheckCollision()

		g.input.Update(g)
	case ModeOver:
		if g.input.IsKeyPressed() {
			g.init()
			g.mode = ModeTitle
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.cfg.BgColor)
	// screen.Fill(color.RGBA{R: g.cfg.BgColor.R, G: g.cfg.BgColor.G, B: g.cfg.BgColor.B, A: g.cfg.BgColor.A})

	var titleTexts []string
	var texts []string

	switch g.mode {
	case ModeTitle:
		titleTexts = []string{"ALIEN INVASION"}
		texts = []string{"", "", "", "", "", "", "", "PRESS ENTER KEY", "", "OR RIGHT MOUSE"}
	case ModeGame:

		// ship
		g.ship.Draw(screen, g.cfg)

		// bullet
		for bullet := range g.bullets {
			bullet.Draw(screen)
		}

		// alien
		for alien := range g.aliens {
			alien.Draw(screen)
		}
		//x := (g.cfg.ScreenWidth - 1*g.cfg.FontSize) / 2
		score := strconv.Itoa(g.score)
		msg := "Your Score: " + score
		//var decodeBytes, _ = simplifiedchinese.GBK.NewDecoder().Bytes([]byte(msg))
		//str := string(decodeBytes)
		text.Draw(screen, msg, arcadeFont, g.cfg.ScreenWidth-len(msg)*g.cfg.FontSize, g.cfg.ScreenHeight-g.cfg.FontSize, g.cfg.ScoreColor)
	case ModeOver:
		texts = []string{g.overMsg}
	}

	for i, l := range titleTexts {
		x := (g.cfg.ScreenWidth - len(l)*g.cfg.TitleFontSize) / 2
		text.Draw(screen, l, titleArcadeFont, x, (i+4)*g.cfg.TitleFontSize, color.White)
	}
	for i, l := range texts {
		x := (g.cfg.ScreenWidth - len(l)*g.cfg.FontSize) / 2
		text.Draw(screen, l, arcadeFont, x, (i+4)*g.cfg.FontSize, g.cfg.ScoreColor)
	}

	// ebitenutil.DebugPrint(screen, g.input.Msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	//return g.cfg.ScreenWidth/2, g.cfg.ScreenHeight/2
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}

func (g *Game) addBullet(bullet *Bullet) {
	g.bullets[bullet] = struct{}{}
}

func (g *Game) addAlien(alien *Alien) {
	g.aliens[alien] = struct{}{}
}

func (g *Game) CreateAliens() {
	alien := NewAlien(g.cfg)
	availableSpaceX := g.cfg.ScreenWidth - 2*alien.width
	numAliens := availableSpaceX / (2 * alien.width)

	for row := 0; row < 2; row++ {
		for i := 0; i < numAliens; i++ {
			alien = NewAlien(g.cfg)
			alien.x = float64(alien.width + 2*alien.width*i)
			alien.y = float64(alien.height*row) * 1.5
			g.addAlien(alien)
		}
	}
}

func (g *Game) createRandomAlien() (alien *Alien) {
	alien = NewAlien(g.cfg)
	width := rand.Intn(g.cfg.ScreenWidth - alien.width)
	alien.x = float64(width)
	alien.y = float64(0)
	// 是否需要检查
	needCheck := true
	for {
		isCollision := false
		// 无限循环判断，直到不重叠
		if !needCheck {
			break
		}

		for alienExist := range g.aliens {
			// 如果和现有的重叠了就重新生成
			if CheckCollision(alien, alienExist) {
				isCollision = true
				width = rand.Intn(g.cfg.ScreenWidth - alien.width)
				alien.x = float64(width)
			}

		}
		if isCollision {
			//time.Sleep(50 * time.Millisecond)
			log.Printf("重叠了，重新生成中-------------")
		}
		// 有重叠，则进行下一轮检查
		needCheck = isCollision
	}
	return alien
}

// 判断是否命中目标
func (g *Game) CheckCollision() {

	for alien := range g.aliens {
		for bullet := range g.bullets {
			if CheckCollision(bullet, alien) {
				delete(g.aliens, alien)
				delete(g.bullets, bullet)
				g.score++
				newAlien := g.createRandomAlien()
				g.addAlien(newAlien)
			}
		}
	}
}

// 创建字体
var (
	titleArcadeFont font.Face
	arcadeFont      font.Face
	smallArcadeFont font.Face
)

func (g *Game) CreateFonts() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	titleArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(g.cfg.TitleFontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(g.cfg.FontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	smallArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(g.cfg.SmallFontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

}

func (g *Game) init() {
	g.CreateAliens()
	g.CreateFonts()
	g.score = 0
	g.failCount = 0
	g.overMsg = ""
}

var alienXStatus = make(map[int]int64)

func moveAlienX(g *Game) {
	// 同时只能有一个线程操作
	mutex.Lock()
	defer mutex.Unlock()

	lastUpdateTime := alienXStatus[0]

	// 一秒更新一次x轴位置
	if time.Now().UnixMilli()-lastUpdateTime < 1000 {
		return
	}

	// 备份移动前的数据
	newAliens := copyAliens(g.aliens)

	//log.Printf("old is %v,new is %v \n", g.aliens, newAliens)
	needCheck := true
	// 尝试次数，次数不宜太多，否则会引起卡顿
	var tryCount int
	for {
		if !needCheck || tryCount >= 3 {
			break
		}
		isCollision := false
		for alien := range newAliens {
			xMove := random(-1, 2)
			x := alien.x + float64(xMove)*alien.speedFactor*20
			// 必须要在屏幕范围内，并且无碰撞
			if x >= 0 && x <= float64(g.cfg.ScreenWidth-alien.width) {
				// 移动
				alien.x = x
			}
		}

		// 等全体移动完再检查一遍
		for alienA := range newAliens {
			for alienB := range newAliens {
				if alienA != alienB && CheckCollision(alienA, alienB) {
					isCollision = true
				}
			}
		}
		if isCollision {
			// todo 检查出重叠了那就重新来一遍。这个有风险，风险在于如果永远没有空位，会死循环，所有要加个尝试次数的判断
			// 在空位少的时候会等待很久，如果开启go异步，画面则会突然发生变动
			newAliens = copyAliens(g.aliens)
			//log.Printf("有碰撞，重置数据,newAliens is %v\n", newAliens)
			tryCount++
		}
		needCheck = isCollision
		//log.Printf("needCheck is %v", needCheck)
	}
	// 移动成功后替换game对象的数据，不成功就不移动
	if !needCheck {
		// 进行替换
		//log.Println("进行替换---------------")
		g.aliens = newAliens
	}

	alienXStatus[0] = time.Now().UnixMilli()

}

func copyAliens(aliensA map[*Alien]struct{}) map[*Alien]struct{} {
	newAliens := make(map[*Alien]struct{})
	for alien := range aliensA {
		var alienBak Alien
		copier.CopyWithOption(&alienBak, alien, copier.Option{
			IgnoreEmpty:   true,
			DeepCopy:      true,
			CaseSensitive: false,
		})
		newAliens[&alienBak] = struct{}{}
	}
	return newAliens
}
