package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"log"
)

const (
	screenWidth  = 340
	screenHeight = 600
	paddleWidth  = 100
	paddleHeight = 20
	ballSize     = 15
	brickWidth   = 60
	brickHeight  = 20
	numBricksRow = 10
	numBricksCol = 5
	scoreHeight  = 20
)

type Game struct {
	paddleX    float32
	ballX      float32
	ballY      float32
	ballVX     float32
	ballVY     float32
	bricks     [][]bool
	score      int
	isGameOver bool
}

func (g *Game) Init() {
	g.bricks = make([][]bool, numBricksCol)
	for i := range g.bricks {
		g.bricks[i] = make([]bool, numBricksRow)
		for j := range g.bricks[i] {
			g.bricks[i][j] = true
		}
	}
}

func (g *Game) Update() error {
	if g.isGameOver {
		return nil
	}

	// Убедитесь, что все кирпичи уничтожены
	if g.allBricksDestroyed() {
		g.isGameOver = true
		return nil
	}

	//отрисовка игрового состояния
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.paddleX -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.paddleX += 5
	}

	//ограничение платформы внутри игрового поля
	if g.paddleX < 0 {
		g.paddleX = 0
	}
	if g.paddleX+paddleWidth > screenWidth {
		g.paddleX = screenWidth - paddleWidth
	}

	//обновление позиции мяча
	g.ballX += g.ballVX
	g.ballY += g.ballVY

	//ограничение мяча внутри игрового поля
	if g.ballX < 0 || g.ballX+ballSize > screenWidth {
		g.ballVX *= -1
	}
	if g.ballY < scoreHeight {
		g.ballVY *= -1
	}
	//вылет снизу
	if g.ballY+ballSize > screenHeight {
		g.ballX, g.ballY = (screenWidth+ballSize)/2, (screenHeight+ballSize)/2
		g.ballVX, g.ballVY = 3, 3
		g.score -= 10
	}

	//проверка столкновения мяча и платформы
	if g.ballY+ballSize > screenHeight-paddleHeight && g.ballX+ballSize > g.paddleX && g.ballX < g.paddleX+paddleWidth {
		g.ballVY *= -1
		g.ballY = screenHeight - paddleHeight - ballSize //корректировка позиции мяча
	}

	//проверка столкновения мяча и кирпича
	for i := 0; i < numBricksCol; i++ {
		for j := 0; j < numBricksRow; j++ {
			if g.bricks[i][j] {
				brickX := float32(i * (brickWidth + 10))
				brickY := float32(j*(brickHeight+10) + scoreHeight)

				if g.ballX+ballSize > brickX && g.ballX < brickX+brickWidth && g.ballY > brickY && g.ballY < brickY+brickHeight {
					g.bricks[i][j] = false
					g.score++
					g.ballVY *= -1
				}
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.isGameOver {
		// Отрисовка сообщения об окончании игры
		text.Draw(screen, fmt.Sprintf("GAME OVER! Score: %d", g.score), basicfont.Face7x13, screenWidth/2-60, screenHeight/2, color.White)
		return
	}

	//отрисовка платформы
	vector.DrawFilledRect(screen, g.paddleX, screenHeight-paddleHeight, paddleWidth, paddleHeight, color.RGBA{R: 255, G: 255, B: 255, A: 255}, true)

	//отрисовка мяча
	vector.DrawFilledRect(screen, g.ballX, g.ballY, ballSize, ballSize, color.RGBA{R: 255, G: 255, B: 255, A: 255}, true)

	//отрисовка кирпичей
	for i := 0; i < numBricksCol; i++ {
		for j := 0; j < numBricksRow; j++ {
			if g.bricks[i][j] {
				brickX := float32(i * (brickWidth + 10))
				brickY := float32(j*(brickHeight+10) + scoreHeight)
				vector.DrawFilledRect(screen, brickX, brickY, brickWidth, brickHeight, color.RGBA{R: 61, G: 61, B: 61, A: 255}, true)
			}
		}
	}

	//отрисуем счет
	text.Draw(screen, fmt.Sprintf("Score: %d", g.score), basicfont.Face7x13, screenWidth/2-40, 10, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screeWidth, screenHeight int) {
	//установка размеров экрана
	return outsideWidth, outsideHeight
}

func (g *Game) allBricksDestroyed() bool {
	for i := 0; i < numBricksCol; i++ {
		for j := 0; j < numBricksRow; j++ {
			if g.bricks[i][j] {
				return false
			}
		}
	}
	return true
}

func main() {
	game := &Game{
		paddleX: (screenWidth - paddleWidth) / 2,
		ballX:   (screenWidth - ballSize) / 2,
		ballY:   (screenHeight + scoreHeight - ballSize) / 2,
		ballVX:  3,
		ballVY:  3,
		score:   0,
	}
	game.Init()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Арканоид")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
