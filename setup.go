package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strconv"
)

func DrawImage(screen *ebiten.Image, image *ebiten.Image, x int, y int) {
	drawOps := ebiten.DrawImageOptions{}
	drawOps.GeoM.Reset()
	drawOps.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(image, &drawOps)
}

func (game *Game) DrawBackground(screen *ebiten.Image) {
	drawOps := ebiten.DrawImageOptions{}
	const repeat = 3
	backgroundHeight := game.background.Bounds().Dy()
	for count := 0; count < repeat; count += 1 {
		drawOps.GeoM.Reset()
		drawOps.GeoM.Translate(0, -float64(backgroundHeight*count))
		drawOps.GeoM.Translate(0, float64(game.backgroundYView))
		screen.DrawImage(game.background, &drawOps)
	}
}

func LoadImage(imageFile string) *ebiten.Image {
	image, _, err := ebitenutil.NewImageFromFile(imageFile)
	if err != nil {
		fmt.Println("Error loading image:", err)
	}
	return image
}

// LoadSound code from https://github.com/jsantore/SimpleEbitenSound (LoadWav)
func LoadSound(name string, context *audio.Context) *audio.Player {
	soundFile, err := os.Open(name)
	if err != nil {
		fmt.Println("Error Loading sound: ", err)
	}
	s, err := wav.DecodeWithoutResampling(soundFile)
	if err != nil {
		fmt.Println("Error interpreting sound file: ", err)
	}
	soundPlayer, err := context.NewPlayer(s)
	if err != nil {
		fmt.Println("Couldn't create sound player: ", err)
	}
	return soundPlayer
}

func LoadMusic(name string, context *audio.Context) *audio.Player {
	soundFile, err := os.Open(name)
	if err != nil {
		fmt.Println("Error Loading sound: ", err)
	}
	s, err := wav.DecodeWithoutResampling(soundFile)
	if err != nil {
		fmt.Println("Error interpreting sound file: ", err)
	}
	loop := audio.NewInfiniteLoop(s, musicLoopLength)
	musicPlayer, err := context.NewPlayer(loop)
	if err != nil {
		fmt.Println("Couldn't create sound player: ", err)
	}
	return musicPlayer
}

func SpawnEnemy(game *Game) {

	if game.numEnemies == 0 {
		for i := 0; i < maxEnemies; i++ {
			x := rand.Intn(screenWidth-(game.enemy.image.Bounds().Size().X*2)) + game.enemy.image.Bounds().Size().X
			y := rand.Intn(screenHeight)*3 + game.enemy.image.Bounds().Size().Y
			game.enemies = append(game.enemies, &Sprite{image: game.enemy.image, xloc: x, yloc: -y, dY: enemySpeed})
		}
		game.numEnemies = maxEnemies
		game.enemyAlert.Rewind()
		game.enemyAlert.Play()
	}

	for i, enemy := range game.enemies {
		enemy.yloc += enemy.dY
		if enemy.yloc > screenHeight+game.enemy.image.Bounds().Size().Y {
			game.score--
			game.enemies = game.RemoveSprite(game.enemies, i)
			game.numEnemies--
		}
	}
}

func (game *Game) Draw(screen *ebiten.Image) {
	game.DrawBackground(screen)
	DrawImage(screen, game.player.image, game.player.xloc, game.player.yloc)

	for _, shot := range game.lasers {
		DrawImage(screen, shot.image, shot.xloc, shot.yloc)
	}
	for _, enemy := range game.enemies {
		DrawImage(screen, enemy.image, enemy.xloc, enemy.yloc)
	}

	game.backgroundMusic.Play()

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	mplusNormalFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	text.Draw(screen, "score: "+strconv.Itoa(game.score), mplusNormalFont, 20, 30, color.White)
}
