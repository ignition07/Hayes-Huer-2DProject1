package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Player struct {
	image *ebiten.Image
	xloc  int
	yloc  int
	dX    int
	dY    int
}

func MovePlayer(game *Game) {

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		game.player.dX = -playerSpeed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		game.player.dX = playerSpeed
	} else if (inpututil.IsKeyJustReleased(ebiten.KeyA) && game.player.dX < 0) ||
		(inpututil.IsKeyJustReleased(ebiten.KeyD) && game.player.dX > 0) {
		game.player.dX = 0
	}
	game.player.xloc += game.player.dX
	if game.player.xloc <= 0 {
		game.player.dX = 0
		game.player.xloc = 0
	} else if game.player.xloc+game.player.image.Bounds().Size().X > screenWidth {
		game.player.dX = 0
		game.player.xloc = screenWidth - game.player.image.Bounds().Size().X
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		game.player.dY = -playerSpeed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		game.player.dY = playerSpeed
	} else if (inpututil.IsKeyJustReleased(ebiten.KeyW) && game.player.dY < 0) ||
		(inpututil.IsKeyJustReleased(ebiten.KeyS) && game.player.dY > 0) {
		game.player.dY = 0
	}
	game.player.yloc += game.player.dY
	if game.player.yloc <= 0 {
		game.player.dY = 0
		game.player.yloc = 0
	} else if game.player.yloc+game.player.image.Bounds().Size().Y > screenHeight {
		game.player.dY = 0
		game.player.yloc = screenHeight - game.player.image.Bounds().Size().Y
	}
}

func ShootLaser(game *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		x := game.player.xloc + 30
		y := game.player.yloc
		game.lasers = append(game.lasers, &Sprite{image: game.laser.image, xloc: x, yloc: y, dY: laserSpeed})

		game.laserShot.Rewind()
		game.laserShot.Play()
	}
	for i, shot := range game.lasers {
		shot.yloc += shot.dY
		if shot.yloc == 0 {
			game.lasers = game.RemoveSprite(game.lasers, i)
		}
		for j, enemy := range game.enemies {
			if game.CheckCollision(shot, enemy) {
				game.lasers = game.RemoveSprite(game.lasers, i)
				game.enemies = game.RemoveSprite(game.enemies, j)

				game.explosionSound.Rewind()
				game.explosionSound.Play()

				game.score += 2
				game.numEnemies--

			}
		}
	}
}
