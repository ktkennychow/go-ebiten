package main

import (
	"math"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	position   Vector
	rotation   float64
	shootTimer *Timer
	sprite     *ebiten.Image
	IsShooting bool
}

func NewPlayer() *Player {
	sprite := PlayerSprites[rand.IntN(len(PlayerSprites))]
	bounds := sprite.Bounds()

	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos := Vector{
		X: ScreenWidth/2 - halfW,
		Y: ScreenHeight/2 - halfH,
	}

	return &Player{
		position:   pos,
		sprite:     sprite,
		shootTimer: NewTimer(1 * time.Second),
	}
}

func (p *Player) Collider() Rect {
	bounds := p.sprite.Bounds()

	return NewRect(
		p.position.X,
		p.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}

func (p *Player) Update() {
	speed := float64(300 / ebiten.TPS())
	rotationSpeed := math.Pi / float64(ebiten.TPS())

	var delta Vector

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= rotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += rotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		// calculate depends on current position and rotation
		delta.X = math.Sin(p.rotation) * speed
		delta.Y = math.Cos(p.rotation) * -speed
	}

	p.position.X += delta.X
	p.position.Y += delta.Y

	p.shootTimer.Update()
	if ebiten.IsKeyPressed(ebiten.KeySpace) && p.shootTimer.IsReady() {
		p.shootTimer.Reset()
		p.IsShooting = true
	} else {
		p.IsShooting = false
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	bounds := p.sprite.Bounds()

	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(p.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(p.position.X, p.position.Y)

	screen.DrawImage(p.sprite, op)
}
