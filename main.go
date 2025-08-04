package main

import (
	"embed"
	"fmt"
	"io/fs"

	"image"
	_ "image/png"
	"log"
	"math"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

type Vector struct {
	X float64
	Y float64
}

type Timer struct {
	currentTick int
	targetTick  int
}

type Player struct {
	position   Vector
	rotation   float64
	shootTimer *Timer
	sprite     *ebiten.Image
	IsShooting bool
}

type Bullet struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
}

type Meteor struct {
	position      Vector
	rotation      float64
	rotationSpeed float64
	movement      Vector
	sprite        *ebiten.Image
}

type Game struct {
	player           *Player
	meteorSpawnTimer *Timer
	meteors          []*Meteor
	bullets          []*Bullet
}

//go:embed assets/*
var assets embed.FS
var PlayerSprites = mustLoadImages("assets/PNG/PlayerShips/*.png")
var MeteorSprites = mustLoadImages("assets/PNG/Meteors/*.png")

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func mustLoadImages(pattern string) []*ebiten.Image {
	files, err := fs.Glob(assets, pattern)
	if err != nil {
		panic(err)
	}
	images := make([]*ebiten.Image, len(files))
	for i, file := range files {
		images[i] = mustLoadImage(file)
	}
	return images
}

// get just the direction without the length
func (v *Vector) Normalize() Vector {
	length := math.Sqrt(v.X*v.X + v.Y*v.Y)
	return Vector{
		X: v.X / length,
		Y: v.Y / length,
	}
}

func NewTimer(d time.Duration) *Timer {
	return &Timer{
		currentTick: 0,
		targetTick:  int(d.Milliseconds()) / 1000 * ebiten.TPS(),
	}
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

func NewBullet(origin Vector, rotation float64) *Bullet {
	sprite := mustLoadImage("assets/PNG/Lasers/laserBlue01.png")
	return &Bullet{
		position: origin,
		rotation: rotation,
		sprite:   sprite,
	}
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

	cm := colorm.ColorM{}
	cm.Translate(1.0, 1.0, 1.0, 0.0)
	cmOp := &colorm.DrawImageOptions{}
	cmOp.GeoM.Translate(-halfW*0.1, -halfH*0.1)
	cmOp.GeoM.Rotate(p.rotation)
	cmOp.GeoM.Translate(halfW*0.1, halfH*0.1)
	cmOp.GeoM.Scale(0.1, 0.1)

	fmt.Println(p.rotation)
	// -halfW * 0.1 to halfW * 0.1 depending on rotation, rotation = 0 ->  halfW * 0.1 but 180 is -halfW * 0.1
	delta := Vector{
		X: math.Cos(p.rotation) * halfW * 0.1,
		Y: math.Sin(p.rotation) * halfH * 0.1,
	}
	fmt.Println(delta)
	cmOp.GeoM.Translate(p.position.X+halfW+delta.X, p.position.Y+halfH+delta.Y)

	screen.DrawImage(p.sprite, op)
	colorm.DrawImage(screen, p.sprite, cm, cmOp)
}

func (b *Bullet) Update() {
	speed := float64(1000 / ebiten.TPS())
	// use delta
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

func (t *Timer) Update() {
	if t.currentTick < t.targetTick {
		t.currentTick++
	}
}

func (t *Timer) IsReady() bool {
	return t.currentTick >= t.targetTick
}

func (t *Timer) Reset() {
	t.currentTick = 0
}

func (g *Game) Update() error {
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

		bulletSpawnOffset := float64(playerBounds.Size().Y / 2)

		playerCenter := Vector{
			X: g.player.position.X + halfW + math.Sin(g.player.rotation)*bulletSpawnOffset,
			Y: g.player.position.Y + halfH - math.Cos(g.player.rotation)*bulletSpawnOffset,
		}

		bullet := NewBullet(playerCenter, g.player.rotation)
		g.bullets = append(g.bullets, bullet)
	}

	for _, meteor := range g.meteors {
		meteor.Update()
	}

	for _, bullet := range g.bullets {
		bullet.Update()
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	p := NewPlayer()
	g := &Game{
		player:           p,
		meteorSpawnTimer: NewTimer(1 * time.Second),
		meteors:          []*Meteor{},
		bullets:          []*Bullet{},
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
