package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bullet struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
}

func NewBullet(origin Vector, rotation float64) *Bullet {
	sprite := mustLoadImage("assets/PNG/Lasers/laserBlue01.png")
	return &Bullet{
		position: origin,
		rotation: rotation,
		sprite:   sprite,
	}
}

func (b *Bullet) Collider() Rect {
	bounds := b.sprite.Bounds()

	return NewRect(
		b.position.X,
		b.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}

func (b *Bullet) Update() {
	speed := float64(1000 / ebiten.TPS())
	// use delta for constant speed for all 8 directions
	delta := Vector{
		X: math.Sin(b.rotation) * speed,
		Y: math.Cos(b.rotation) * -speed,
	}

	b.position.X += delta.X
	b.position.Y += delta.Y
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	bounds := b.sprite.Bounds()

	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(b.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(b.position.X, b.position.Y)

	screen.DrawImage(b.sprite, op)
}
