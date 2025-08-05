package main

import (
	"embed"
	"io/fs"

	"image"
	_ "image/png"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

type Rect struct {
	X float64
	Y float64
	W float64
	H float64
}

type Vector struct {
	X float64
	Y float64
}

type Timer struct {
	currentTick int
	targetTick  int
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

func main() {
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowSize(ScreenWidth*2, ScreenHeight*2)
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
