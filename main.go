package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
)

const SCREEN_WIDTH = 100
const SCREEN_HEIGHT = 50

const GRAVITY = 0.01

var sizeChars = [11]string{" ", ".", "-", "=", "+", "*", "o", "O", "0", "@", "#"}

type Color struct {
	r uint8
	g uint8
	b uint8
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
	childrenDecayRate    float64
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

func (f Firework) Spawn() {
	for i := 0; i < f.particleCount; i++ {
		angle := float64(i) * (2. * math.Pi / float64(f.particleCount))

		newParticle := Particle{
			posX: f.posX, posY: f.posY,
			velX:      math.Cos(angle+(rand.Float64()*f.particleVelocitySpread)) * f.particleVelocity,
			velY:      math.Sin(angle+(rand.Float64()*f.particleVelocitySpread)) * f.particleVelocity,
			decayRate: 0.02, childrenDecayRate: 0.04,
			colorDecayRate:       f.particleDecayRate,
			pixel:                Pixel{size: .9, color: f.particleColor},
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

func (p *Particle) ChildrenUpdate() {
	p.pixel.size -= p.decayRate
	p.pixel.color.r -= p.colorDecayRate.r
	p.pixel.color.g -= p.colorDecayRate.g
	p.pixel.color.b -= p.colorDecayRate.b
}

func (p *Particle) Update() {
	oldPosX := int(p.posX)
	oldPosY := int(p.posY)
	p.posX += p.velX
	p.posY += p.velY
	p.velY += GRAVITY
	p.pixel.size -= p.decayRate
	if oldPosX != int(p.posX) || oldPosY != int(p.posY) {
		newParticle := *p
		newParticle.decayRate = p.childrenDecayRate
		childrenParticles = append(childrenParticles, newParticle)
	}
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
			// if i == 0 || i == SCREEN_HEIGHT-1 || j == 0 || j == SCREEN_WIDTH-1 {
			// 	char = "#"
			// }

			toPrint += fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s", screen[i][j].color.r, screen[i][j].color.g, screen[i][j].color.b, char)
			// toPrint += "\x1b[38;2;" +
			// 	string(screen[i][j].color.r) + ";" +
			// 	string(screen[i][j].color.g) + ";" +
			// 	string(screen[i][j].color.b) + "m" + char
			// toPrint += char
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

var particles = []Particle{}
var childrenParticles = []Particle{}

var possiblePositionsX = []int{75, 50, 25}
var possiblePositionsY = []int{25, 40, 10}

var possibleColors = []Color{{255, 255, 0}, {255, 0, 255}, {0, 255, 255}}
var possibleDecayRates = []Color{{10, 0, 0}, {0, 0, 10}, {0, 10, 0}}

func main() {
	ClearScreen()

	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	//makes cursor invisible
	fmt.Print("\033[?25l")

	running := true
	frame := 0

	posXRandIndex := rand.Int() % 3
	posYRandIndex := rand.Int() % 3
	colorRandIndex1 := rand.Int() % 3
	colorRandIndex2 := rand.Int() % 3

	posXRandIndex = 1
	posYRandIndex = 0

	randPosX := rand.Int() % 5
	randPosY := rand.Int() % 5

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

	for running {
		frame += 1
		fmt.Print(frame)
		ClearScreen()

		if frame%25 == 0 && frame > 50 {
			posXRandIndex := rand.Int() % 3
			posYRandIndex := rand.Int() % 3
			colorRandIndex1 := rand.Int() % 3
			colorRandIndex2 := rand.Int() % 3

			randPosX := rand.Int() % 5
			randPosY := rand.Int() % 5

			firework := Firework{
				posX:                   float64(possiblePositionsX[posXRandIndex] + randPosX),
				posY:                   float64(possiblePositionsY[posYRandIndex] + randPosY),
				particleCount:          10,
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
				particleCount:          10,
				particleVelocity:       .5,
				particleVelocitySpread: .5,
				particleColor:          possibleColors[colorRandIndex2],
				particleDecayRate:      possibleDecayRates[colorRandIndex2],
				hasExplodingChildren:   false,
			}
			firework2.Spawn()
		}

		//set cursor pos to 0, 0

		for i := 0; i < len(particles); i++ {
			particles[i].Update()
		}
		for i := 0; i < len(childrenParticles); i++ {
			childrenParticles[i].ChildrenUpdate()
		}

		for i := 0; i < len(particles); i++ {
			particles[i].Draw()
		}
		for i := 0; i < len(childrenParticles); i++ {
			childrenParticles[i].Draw()
		}
		DrawScreen()
	}

	//makes cursor visible
	fmt.Print("\033[?25h")

	//changes terminal colors to default
	fmt.Print("\x1b[0m\n")
}
