package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Game struct {
	player           *Player
	meteorSpawnTimer *Timer
	meteors          []*Meteor
	bullets          []*Bullet
	score            int
}

func (g *Game) Reset() {
	g.player = NewPlayer()
	g.meteors = nil
	g.bullets = nil
	g.score = 0
}

func (g *Game) Update() error {
	bullet := NewBullet(Vector{}, g.player.rotation)

	g.player.Update()
	g.meteorSpawnTimer.Update()
	if g.meteorSpawnTimer.IsReady() {
		g.meteorSpawnTimer.Reset()

		g.meteors = append(g.meteors, NewMeteor())
	}

	if g.player.IsShooting {
		playerBounds := g.player.sprite.Bounds()
		halfW := float64(playerBounds.Dx()) / 2
		halfH := float64(playerBounds.Dy()) / 2

		playerCenter := Vector{
			X: g.player.position.X + halfW + math.Sin(g.player.rotation),
			Y: g.player.position.Y + halfH - math.Cos(g.player.rotation),
		}

		bulletBounds := bullet.sprite.Bounds()
		bulletHalfW := float64(bulletBounds.Dx()) / 2
		bulletHalfH := float64(bulletBounds.Dy()) / 2

		playerDimensionAsSquare := math.Min(float64(playerBounds.Dx()), float64(playerBounds.Dy()))

		// need a delta from playerbounds
		delta := Vector{
			X: math.Sin(g.player.rotation) * playerDimensionAsSquare / 2,
			Y: math.Cos(g.player.rotation) * playerDimensionAsSquare / 2,
		}

		bullet.position = Vector{
			X: playerCenter.X - bulletHalfW + delta.X,
			Y: playerCenter.Y - bulletHalfH - delta.Y,
		}

		g.bullets = append(g.bullets, bullet)
	}

	for _, meteor := range g.meteors {
		meteor.Update()
	}

	for _, bullet := range g.bullets {
		bullet.Update()
	}

	for i, meteor := range g.meteors {
		for j, bullet := range g.bullets {
			if bullet.Collider().Intersects(meteor.Collider()) {
				g.score++
				g.bullets = append(g.bullets[:j], g.bullets[j+1:]...)
				g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
			}
		}
	}

	for _, meteor := range g.meteors {
		if g.player.Collider().Intersects(meteor.Collider()) {
			g.Reset()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)

	for _, meteor := range g.meteors {
		meteor.Draw(screen)
	}

	for _, bullet := range g.bullets {
		bullet.Draw(screen)
	}

	scoreText := fmt.Sprintf("SCORE: %06d", g.score)
	scoreGeoM := ebiten.GeoM{}
	scoreGeoM.Translate(20, 30)

	text.Draw(screen, scoreText, ScoreFont, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: scoreGeoM,
		},
	})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
