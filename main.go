package main

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"log"
	"os"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	x, y      float64
	dx, dy    float64
	img       *ebiten.Image
	audioData []byte
	ctx       *oto.Context
	player    *oto.Player
	state     string // "menu" or "playing"
	lives     int
	score     int
	paddleX   float64
}

func (g *Game) Update() error {
	if g.state == "menu" {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = "playing"
			g.lives = 3
			g.score = 0
			g.x, g.y = 320, 240
			g.dx, g.dy = 2, 2
			g.paddleX = 320 - 40 // largeur raquette 80, centrer
		}
		return nil
	}

	// État de jeu
	// Déplacer la raquette
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.paddleX > 0 {
		g.paddleX -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && g.paddleX < 640-80 {
		g.paddleX += 5
	}

	// Mettre à jour la balle
	g.x += g.dx
	g.y += g.dy
	bounced := false
	if g.x < 0 || g.x > 640-16 {
		g.dx = -g.dx
		bounced = true
	}
	if g.y < 0 {
		g.dy = -g.dy
		bounced = true
	}
	// Vérifier collision avec raquette
	paddleY := float64(480 - 20)
	if g.y > paddleY-16 && g.y < paddleY && g.x > g.paddleX && g.x < g.paddleX+80 {
		g.dy = -g.dy
		bounced = true
		g.score += 10
	}
	// Vérifier si balle sort en bas
	if g.y > 480 {
		g.lives--
		if g.lives <= 0 {
			g.state = "menu"
		} else {
			// Réinitialiser la balle
			g.x, g.y = 320, 240
			g.dx, g.dy = 2, 2
		}
	}
	if bounced {
		g.playSound()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state == "menu" {
		ebitenutil.DebugPrint(screen, "Jeu de Balle Rebondissante\n\nAppuyez sur ESPACE pour commencer")
		return
	}

	// État de jeu
	// Dessiner la balle
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.x, g.y)
	screen.DrawImage(g.img, op)

	// Dessiner la raquette
	vector.FillRect(screen, float32(g.paddleX), 460, 80, 20, color.White, false)

	// Afficher vies et score
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Vies: %d", g.lives), 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.score), 10, 30)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func (g *Game) playSound() {
	if g.player != nil {
		g.player.Seek(0, io.SeekStart)
		g.player.Play()
	}
}

func loadImage() *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile("ball.png")
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func loadAudio() ([]byte, error) {
	file, err := os.Open("bounce.wav")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// Ignorer l'en-tête WAV (44 octets)
	header := make([]byte, 44)
	file.Read(header)
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func main() {
	img := loadImage()
	audioData, err := loadAudio()
	if err != nil {
		log.Fatal(err)
	}
	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	})
	if err != nil {
		log.Fatal(err)
	}
	<-ready
	player := ctx.NewPlayer(bytes.NewReader(audioData))
	game := &Game{
		x: 320, y: 240,
		dx: 2, dy: 2,
		img:       img,
		audioData: audioData,
		ctx:       ctx,
		player:    player,
		state:     "menu",
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Jeu de Balle Rebondissante")
	os.Setenv("EBITEN_GRAPHICS_LIBRARY", "opengl")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
