package main

import (
	"fmt"
	"github.com/co0p/tankism/lib/collision"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	_ "image/png"
)

// CREDITS:

// background music: "RetroRide" by tubelesshalo (49 seconds)
// https://tubelesshalo.itch.io/free-music-loop-pack

// laser and alarm sound by Pixabay
// https://pixabay.com/sound-effects/search/game/

// art: Nico Hayes-Huer

const (
	screenWidth     = 1000
	screenHeight    = 1200
	soundSampleRate = 48000
	playerSpeed     = 8
	enemySpeed      = 6
	laserSpeed      = -10
	musicLoopLength = 8740380
	maxEnemies      = 10
	maxLasers       = 100
)

type Sprite struct {
	image *ebiten.Image
	xloc  int
	yloc  int
	dY    int
}

type Game struct {
	player          Player
	enemy           Sprite
	enemies         []*Sprite
	laser           Sprite
	lasers          []*Sprite
	explosionImage  Sprite
	background      *ebiten.Image
	backgroundXView int
	backgroundYView int
	audioContext    *audio.Context
	explosionSound  *audio.Player
	laserShot       *audio.Player
	enemyAlert      *audio.Player
	backgroundMusic *audio.Player
	score           int
	numEnemies      int
}

func (game *Game) Update() error {
	MovePlayer(game)
	ShootLaser(game)
	SpawnEnemy(game)

	backgroundHeight := game.background.Bounds().Dy()
	maxY := backgroundHeight * 2
	game.backgroundYView += 2
	game.backgroundYView %= maxY

	return nil
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hayes-Huer | 2D Project 1")

	// sound
	audioContext := audio.NewContext(soundSampleRate)
	explosionSoundFile := LoadSound("assets/Thunder1.wav", audioContext)
	laserSoundFile := LoadSound("assets/shoot.wav", audioContext)
	alertSoundFile := LoadSound("assets/alarm.wav", audioContext)
	backgroundMusicFile := LoadMusic("assets/RetroRide.wav", audioContext)

	// images
	background := LoadImage("assets/background.png")
	playerShip := LoadImage("assets/player_ship.png")
	enemyShip := LoadImage("assets/enemy_ship.png")
	laserImage := LoadImage("assets/laser.png")
	boom := LoadImage("assets/boom1.png")

	lasers := make([]*Sprite, 0, maxLasers)
	enemies := make([]*Sprite, 0, maxEnemies)

	game := Game{
		audioContext:    audioContext,
		explosionSound:  explosionSoundFile,
		laserShot:       laserSoundFile,
		enemyAlert:      alertSoundFile,
		backgroundMusic: backgroundMusicFile,
		background:      background,
		lasers:          lasers,
		enemies:         enemies,
		numEnemies:      0,
		score:           0,
	}
	game.player = Player{image: playerShip, xloc: (screenWidth - playerShip.Bounds().Size().X) / 2,
		yloc: 1000, dX: 0, dY: 0}
	game.enemy = Sprite{image: enemyShip, xloc: 0, yloc: 0, dY: 0}
	game.laser = Sprite{image: laserImage, xloc: 0, yloc: 0, dY: 0}
	game.explosionImage = Sprite{image: boom, xloc: 0, yloc: 0, dY: 0}

	err := ebiten.RunGame(&game)
	if err != nil {
		fmt.Println("Failed to run game", err)
	}
}

func (game *Game) RemoveSprite(spriteList []*Sprite, i int) []*Sprite {
	s := spriteList
	if i >= 0 && i < len(s) {
		s[i] = s[len(s)-1]
	}
	return s[:len(s)-1]
}

// CheckCollision from class slides
func (game *Game) CheckCollision(sprite1 *Sprite, sprite2 *Sprite) bool {
	sprite1Bounds := collision.BoundingBox{
		X:      float64(sprite1.xloc),
		Y:      float64(sprite1.yloc),
		Width:  float64(sprite1.image.Bounds().Dx()),
		Height: float64(sprite1.image.Bounds().Dy()),
	}
	sprite2Bounds := collision.BoundingBox{
		X:      float64(sprite2.xloc),
		Y:      float64(sprite2.yloc),
		Width:  float64(sprite2.image.Bounds().Dx()),
		Height: float64(sprite2.image.Bounds().Dy()),
	}
	return collision.AABBCollision(sprite1Bounds, sprite2Bounds)
}
