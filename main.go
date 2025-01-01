package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
)

const SCREEN_WIDTH = 100
const SCREEN_HEIGHT = 50

const GRAVITY = 0.2

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

func (p Particle) Draw() {
	posXInt := int(p.posX)
	posYInt := int(p.posY)
	if posYInt < SCREEN_HEIGHT && posXInt < SCREEN_WIDTH && posXInt >= 0 && posYInt >= 0 {
		screen[posYInt][posXInt] = p.pixel
	}
}

func (p *Particle) ChildrenUpdate() {
	// p.posX += p.velX
	// p.posY += p.velY
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
		fmt.Print(len(childrenParticles))
	}
}

func DrawScreen() {
	for i := 0; i < SCREEN_HEIGHT; i++ {
		for j := 0; j < SCREEN_WIDTH; j++ {

			char := "#"
			size := screen[i][j].size
			if size <= 0. {
				char = " "
			} else if size <= 1. {
				char = sizeChars[int(math.Round(size*10.))]
			}

			fmt.Printf("\x1b[38;2;%d;%d;%dm", screen[i][j].color.r, screen[i][j].color.g, screen[i][j].color.b)
			fmt.Printf("%s", char)
		}
		fmt.Print("\n")
	}
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

func main() {

	ClearScreen()

	particle := Particle{
		posX: 10, posY: 10,
		velX: 1., velY: 1.,
		decayRate: 0.001, childrenDecayRate: 0.04,
		colorDecayRate: Color{10, 0, 0},
		pixel:          Pixel{size: .9, color: Color{255, 255, 0}},
	}

	particles = append(particles, particle)

	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()

	//makes cursor invisible
	fmt.Print("\033[?25l")

	running := true

	for running {
		ClearScreen()
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
		fmt.Println(particle.posX)
	}

	//makes cursor visible
	fmt.Print("\033[?25h")

	//changes terminal colors to default
	fmt.Print("\x1b[0m\n")
}
