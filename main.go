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

const GRAVITY = 0.02

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
	posX              float64
	posY              float64
	velX              float64
	velY              float64
	decayRate         float64
	childrenDecayRate float64
	colorDecayRate    Color
	pixel             Pixel
}

type Firework struct {
	posX                   float64
	posY                   float64
	particleCount          int
	particleVelocity       float64
	particleVelocitySpread float64
	particleColor          Color
	particleDecayRate      Color
}

func (f Firework) Spawn() {
	for i := 0; i < f.particleCount; i++ {
		angle := float64(i) * (2. * math.Pi / float64(f.particleCount))

		newParticle := Particle{
			posX: f.posX, posY: f.posY,
			velX:      math.Cos(angle+(rand.Float64()*f.particleVelocitySpread)) * f.particleVelocity,
			velY:      math.Sin(angle+(rand.Float64()*f.particleVelocitySpread)) * f.particleVelocity,
			decayRate: 0.02, childrenDecayRate: 0.04,
			colorDecayRate: f.particleDecayRate,
			pixel:          Pixel{size: .9, color: f.particleColor},
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
}

func DrawScreen() {
	toPrint := ""
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

var possiblePositionsX = []int{}
var possiblePositionsY = []int{}

var possibleColors = []Color{}
var possibleDecayRates = []Color{}

func main() {

	possiblePositionsX = append(possiblePositionsX, 75)
	possiblePositionsX = append(possiblePositionsX, 50)
	possiblePositionsX = append(possiblePositionsX, 25)

	possiblePositionsY = append(possiblePositionsY, 25)
	possiblePositionsY = append(possiblePositionsY, 40)
	possiblePositionsY = append(possiblePositionsY, 10)

	possibleColors = append(possibleColors, Color{255, 255, 0})
	possibleColors = append(possibleColors, Color{255, 0, 255})
	possibleColors = append(possibleColors, Color{0, 255, 255})

	possibleDecayRates = append(possibleDecayRates, Color{10, 0, 0})
	possibleDecayRates = append(possibleDecayRates, Color{0, 0, 10})
	possibleDecayRates = append(possibleDecayRates, Color{0, 10, 0})
	ClearScreen()

	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	//makes cursor invisible
	fmt.Print("\033[?25l")

	running := true
	frame := 0

	for running {
		frame += 1
		fmt.Print(frame)
		ClearScreen()

		if frame%15 == 0 {
			posXRandIndex := rand.Int() % 3
			posYRandIndex := rand.Int() % 3
			colorRandIndex := rand.Int() % 3

			firework := Firework{
				posX:                   float64(possiblePositionsX[posXRandIndex] + rand.Int()%5),
				posY:                   float64(possiblePositionsY[posYRandIndex] + rand.Int()%5),
				particleCount:          10,
				particleVelocity:       1.,
				particleVelocitySpread: .5,
				particleColor:          possibleColors[colorRandIndex],
				particleDecayRate:      possibleDecayRates[colorRandIndex],
			}
			firework.Spawn()
		}

		//set cursor pos to 0, 0
		fmt.Printf("\033[%d;%dH", 0, 0)

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
