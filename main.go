package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const SCREEN_WIDTH = 100
const SCREEN_HEIGHT = 50
const FPS = 30
const GRAVITY = 0.02

var sizeChars = [11]string{" ", ".", "-", "=", "+", "*", "o", "O", "0", "@", "#"}

type Color struct {
	r float64
	g float64
	b float64
}

type Pixel struct {
	size  float64
	color Color
}

var screen [SCREEN_HEIGHT][SCREEN_WIDTH]Pixel

type Particle struct {
	posX                 float64
	posY                 float64
	velX                 float64
	velY                 float64
	decayRate            float64
	colorDecayRate       Color
	pixel                Pixel
	hasExplodingChildren bool
	didExplode           bool
}

type Firework struct {
	posX                   float64
	posY                   float64
	particleCount          int
	particleVelocity       float64
	particleVelocitySpread float64
	particleColor          Color
	particleDecayRate      Color
	hasExplodingChildren   bool
}

func randomFloat(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func (f Firework) Spawn() {
	for i := 0; i < f.particleCount; i++ {
		angle := float64(i) * (2. * math.Pi / float64(f.particleCount))

		particleClr := f.particleColor
		if particleClr.r == 0 {
			particleClr.r += randomFloat(0, 100)

		} else {
			particleClr.r += randomFloat(-100, 0)
		}

		if particleClr.g == 0 {
			particleClr.g += randomFloat(0, 100)

		} else {
			particleClr.g += randomFloat(-100, 0)
		}

		if particleClr.b == 0 {
			particleClr.b += randomFloat(0, 100)
		} else {
			particleClr.b += randomFloat(-100, 0)
		}

		newParticle := Particle{
			posX: f.posX, posY: f.posY,
			velX:                 math.Cos(angle+(rand.Float64()*f.particleVelocitySpread)) * f.particleVelocity,
			velY:                 math.Sin(angle+(rand.Float64()*f.particleVelocitySpread)) * f.particleVelocity,
			decayRate:            0.02,
			colorDecayRate:       f.particleDecayRate,
			pixel:                Pixel{size: .9, color: particleClr},
			hasExplodingChildren: f.hasExplodingChildren,
			didExplode:           false,
		}
		particles = append(particles, newParticle)
	}
}

func (p Particle) Draw() {
	posXInt := int(p.posX)
	posYInt := int(p.posY)
	if posYInt < SCREEN_HEIGHT && posXInt < SCREEN_WIDTH && posXInt >= 0 && posYInt >= 0 {
		screen[posYInt][posXInt] = p.pixel
	}
}

func (p *Particle) Update() {
	p.posX += p.velX
	p.posY += p.velY
	p.velY += GRAVITY
	p.pixel.size -= p.decayRate
	if p.hasExplodingChildren && int(math.Round(p.pixel.size*10.)) == 2 && !p.didExplode {
		p.didExplode = true
		colorRandIndex1 := rand.Int() % 3
		firework := Firework{
			posX:                   p.posX,
			posY:                   p.posY,
			particleCount:          10,
			particleVelocity:       .4,
			particleVelocitySpread: .5,
			particleColor:          possibleColors[colorRandIndex1],
			particleDecayRate:      possibleDecayRates[colorRandIndex1],
			hasExplodingChildren:   false,
		}
		firework.Spawn()
	}
}

func DrawScreen() {
	// sets cursor to 0, 0
	toPrint := "\033[0;0H"
	for i := 0; i < SCREEN_HEIGHT; i++ {
		for j := 0; j < SCREEN_WIDTH; j++ {

			char := "#"
			size := screen[i][j].size
			if size <= 0. {
				char = " "
			} else if size <= 1. {
				char = sizeChars[int(math.Round(size*10.))]
			}

			toPrint += fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s", uint8(screen[i][j].color.r), uint8(screen[i][j].color.g), uint8(screen[i][j].color.b), char)
		}
		toPrint += "\n"
	}
	fmt.Print(toPrint)
}

func ClearScreen() {
	emptyPixel := Pixel{size: .0, color: Color{0, 0, 0}}
	for i := 0; i < SCREEN_HEIGHT; i++ {
		for j := 0; j < SCREEN_WIDTH; j++ {

			screen[i][j] = emptyPixel
		}
	}
}

func FadeScreen() {
	for i := 0; i < SCREEN_HEIGHT; i++ {
		for j := 0; j < SCREEN_WIDTH; j++ {
			screen[i][j].color.r *= 0.85
			screen[i][j].color.g *= 0.85
			screen[i][j].color.b *= 0.85

			if screen[i][j].color.r < 5 && screen[i][j].color.g < 5 {
				screen[i][j].size = 0
			}
		}
	}
}

var particles = []Particle{}

var possiblePositionsX = []int{75, 50, 25}
var possiblePositionsY = []int{25, 40, 10}

var possibleColors = []Color{{255, 255, 0}, {255, 0, 255}, {0, 255, 255}, {255, 0, 0}, {0, 255, 0}, {0, 0, 255}}
var possibleDecayRates = []Color{{10, 0, 0}, {0, 0, 10}, {0, 10, 0}, {10, 0, 0}, {0, 10, 0}, {0, 0, 10}}

func main() {
	ClearScreen()

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	//makes cursor invisible
	fmt.Print("\033[?25l")

	frame := 0

	posXRandIndex := rand.Int() % 3
	posYRandIndex := rand.Int() % 3
	colorRandIndex1 := rand.Int() % 6
	colorRandIndex2 := rand.Int() % 6

	posXRandIndex = 1
	posYRandIndex = 0

	randPosX := rand.Int() % 5
	randPosY := rand.Int() % 5

	ticker := time.NewTicker(time.Second / FPS)
	defer ticker.Stop()

	firework := Firework{
		posX:                   float64(possiblePositionsX[posXRandIndex] + randPosX),
		posY:                   float64(possiblePositionsY[posYRandIndex] + randPosY),
		particleCount:          10,
		particleVelocity:       .75,
		particleVelocitySpread: .5,
		particleColor:          possibleColors[colorRandIndex1],
		particleDecayRate:      possibleDecayRates[colorRandIndex1],
		hasExplodingChildren:   true,
	}
	firework.Spawn()

	firework2 := Firework{
		posX:                   float64(possiblePositionsX[posXRandIndex] + randPosX),
		posY:                   float64(possiblePositionsY[posYRandIndex] + randPosY),
		particleCount:          10,
		particleVelocity:       .5,
		particleVelocitySpread: .5,
		particleColor:          possibleColors[colorRandIndex2],
		particleDecayRate:      possibleDecayRates[colorRandIndex2],
		hasExplodingChildren:   false,
	}
	firework2.Spawn()

	for range ticker.C {
		frame += 1
		FadeScreen()
		fmt.Print(frame)

		if frame%25 == 0 && frame > 50 {
			posXRandIndex := rand.Int() % 3
			posYRandIndex := rand.Int() % 3
			colorRandIndex1 := rand.Int() % 6
			colorRandIndex2 := rand.Int() % 6

			randPosX := rand.Int() % 5
			randPosY := rand.Int() % 5

			firework := Firework{
				posX:                   float64(possiblePositionsX[posXRandIndex] + randPosX),
				posY:                   float64(possiblePositionsY[posYRandIndex] + randPosY),
				particleCount:          20,
				particleVelocity:       .75,
				particleVelocitySpread: .5,
				particleColor:          possibleColors[colorRandIndex1],
				particleDecayRate:      possibleDecayRates[colorRandIndex1],
				hasExplodingChildren:   false,
			}
			firework.Spawn()

			firework2 := Firework{
				posX:                   float64(possiblePositionsX[posXRandIndex] + randPosX),
				posY:                   float64(possiblePositionsY[posYRandIndex] + randPosY),
				particleCount:          20,
				particleVelocity:       .5,
				particleVelocitySpread: .5,
				particleColor:          possibleColors[colorRandIndex2],
				particleDecayRate:      possibleDecayRates[colorRandIndex2],
				hasExplodingChildren:   false,
			}
			firework2.Spawn()
		}

		for i := 0; i < len(particles); i++ {
			particles[i].Update()
		}

		for i := 0; i < len(particles); i++ {
			particles[i].Draw()
		}

		DrawScreen()

		n := 0
		for _, p := range particles {
			if p.pixel.size > 0 {
				particles[n] = p
				n++
			}
		}

		particles = particles[:n]

	}

	//makes cursor visible
	fmt.Print("\033[?25h")

	//changes terminal colors to default
	fmt.Print("\x1b[0m\n")
}
