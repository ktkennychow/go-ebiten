package main

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type Meteor struct {
	position      Vector
	rotation      float64
	rotationSpeed float64
	movement      Vector
	sprite        *ebiten.Image
}

func NewMeteor() *Meteor {
	target := Vector{
		// random within center 50% ofscreen
		X: ScreenWidth/2 + rand.Float64()*ScreenWidth/4,
		Y: ScreenHeight/2 + rand.Float64()*ScreenHeight/4,
	}

	// Pick a random edge: 0=top, 1=right, 2=bottom, 3=left
	edge := rand.IntN(4)

	var pos Vector
	switch edge {
	case 0: // Top edge
		pos = Vector{
			X: rand.Float64()*(ScreenWidth+200) - 100,
			Y: -100,
		}
	case 1: // Right edge
		pos = Vector{
			X: ScreenWidth + 100,
			Y: rand.Float64()*(ScreenHeight+200) - 100,
		}
	case 2: // Bottom edge
		pos = Vector{
			X: rand.Float64()*(ScreenWidth+200) - 100,
			Y: ScreenHeight + 100,
		}
	case 3: // Left edge
		pos = Vector{
			X: -100,
			Y: rand.Float64()*(ScreenHeight+200) - 100,
		}
	}

	velocity := 0.25 + rand.Float64()*2

	direction := Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}

	normalizedDirection := direction.Normalize()

	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	sprite := MeteorSprites[rand.IntN(len(MeteorSprites))]

	rotationSpeed := -0.02 + rand.Float64()*0.06

	return &Meteor{
		position:      pos,
		sprite:        sprite,
		movement:      movement,
		rotationSpeed: rotationSpeed,
	}
}

func (m *Meteor) Collider() Rect {
	bounds := m.sprite.Bounds()

	return NewRect(
		m.position.X,
		m.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}

func (m *Meteor) Update() {
	m.position.X += m.movement.X
	m.position.Y += m.movement.Y
	m.rotation += m.rotationSpeed
}

func (m *Meteor) Draw(screen *ebiten.Image) {
	bounds := m.sprite.Bounds()

	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(m.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(m.position.X, m.position.Y)

	screen.DrawImage(m.sprite, op)
}
