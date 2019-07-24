package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type GimulType int

const (
	GimulTypeNone = -1 + iota
	GimulTypeGreenWang
	GimulTypeGreenJa
	GimulTypeGreenJang
	GimulTypeGreenSang
	GimulTypeGreenQueen
	GimulTypeRedWang
	GimulTypeRedJa
	GimulTypeRedJang
	GimulTypeRedSang
	GimulTypeRedQueen
	GimulTypeMax
)

const (
	ScreenWidth  = 480
	ScreenHeight = 412 //362
	S            = 20
	GridWidth    = 116
	T            = 23
	GridHeight   = 116

	BoardWidth  = 4
	BoardHeight = 3
)

type TeamType int

const (
	TeamNone TeamType = iota
	TeamGreen
	TeamRed
)

var (
	board       [BoardWidth][BoardHeight]GimulType
	bgimg       *ebiten.Image
	gimulImgs   [GimulTypeMax]*ebiten.Image
	selected    bool
	selectedX   int
	selectedY   int
	selectedImg *ebiten.Image
	currentTeam TeamType = TeamGreen
	gameover    bool
	flag        bool
)

func GetTeamType(gimulType GimulType) TeamType {
	if gimulType == GimulTypeGreenJa ||
		gimulType == GimulTypeGreenJang ||
		gimulType == GimulTypeGreenWang ||
		gimulType == GimulTypeGreenSang {
		return TeamGreen
	}
	if gimulType == GimulTypeRedJa ||
		gimulType == GimulTypeRedJang ||
		gimulType == GimulTypeRedWang ||
		gimulType == GimulTypeRedSang {
		return TeamRed
	}
	return TeamNone
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func OnDie(gimulType GimulType) {
	if gimulType == GimulTypeGreenWang ||
		gimulType == GimulTypeRedWang {
		gameover = true
	}
}

func move(prevX, prevY, tarX, tarY int) {
	if isMoveble(prevX, prevY, tarX, tarY) {
		OnDie(board[tarX][tarY])
		board[prevX][prevY], board[tarX][tarY] = GimulTypeNone, board[prevX][prevY]
		selected = false
		if currentTeam == TeamGreen {
			currentTeam = TeamRed
		} else {
			currentTeam = TeamGreen
		}
		if tarX == 3 && board[tarX][tarY] == GimulTypeGreenJa {
			flag = true
		}
	}
}
func isMoveble(prevX, prevY, tarX, tarY int) bool {
	if tarX < 0 || tarY < 0 {
		return false
	}
	if tarX >= BoardWidth || tarY >= BoardHeight {
		return false
	}

	if GetTeamType(board[prevX][prevY]) == GetTeamType(board[tarX][tarY]) {
		return false
	}
	switch board[prevX][prevY] {
	case GimulTypeGreenJa:
		return prevX+1 == tarX && prevY == tarY
	case GimulTypeRedJa:
		return prevX-1 == tarX && prevY == tarY
	case GimulTypeGreenSang, GimulTypeRedSang:
		return abs(prevX-tarX) == 1 && abs(prevY-tarY) == 1
	case GimulTypeGreenJang, GimulTypeRedJang:
		return abs(prevX-tarX)+abs(prevY-tarY) == 1
	case GimulTypeRedWang, GimulTypeGreenWang:
		return abs(prevX-tarX) == 1 || abs(prevY-tarY) == 1
	}
	return false
}
func update(screen *ebiten.Image) error {

	screen.DrawImage(bgimg, nil)
	if gameover {
		return nil
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		i, j := x/GridWidth, y/GridHeight
		if i >= 0 && i < GridWidth && j >= 0 && j < GridHeight {
			if selected {
				if i == selectedX && j == selectedY {
					selected = false
				} else {
					move(selectedX, selectedY, i, j)
				}
			} else {
				if board[i][j] != GimulTypeNone && currentTeam == GetTeamType(board[i][j]) {
					selected = true
					selectedX, selectedY = i, j
				}
			}
		}
	}

	for i := 0; i < BoardWidth; i++ {
		for j := 0; j < BoardHeight; j++ {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(S+GridWidth*i), float64(T+GridHeight*j))
			switch board[i][j] {
			case GimulTypeGreenWang:
				screen.DrawImage(gimulImgs[GimulTypeGreenWang], opts)
			case GimulTypeGreenJa:
				screen.DrawImage(gimulImgs[GimulTypeGreenJa], opts)
				if flag == true {
					screen.DrawImage(gimulImgs[GimulTypeGreenQueen], opts)
				}
			case GimulTypeGreenSang:
				screen.DrawImage(gimulImgs[GimulTypeGreenSang], opts)
			case GimulTypeGreenJang:
				screen.DrawImage(gimulImgs[GimulTypeGreenJang], opts)
			case GimulTypeGreenQueen:
				screen.DrawImage(gimulImgs[GimulTypeGreenQueen], opts)
				//
			case GimulTypeRedWang:
				screen.DrawImage(gimulImgs[GimulTypeRedWang], opts)
			case GimulTypeRedJa:
				screen.DrawImage(gimulImgs[GimulTypeRedJa], opts)
			case GimulTypeRedSang:
				screen.DrawImage(gimulImgs[GimulTypeRedSang], opts)
			case GimulTypeRedJang:
				screen.DrawImage(gimulImgs[GimulTypeRedJang], opts)
			case GimulTypeRedQueen:
				screen.DrawImage(gimulImgs[GimulTypeRedQueen], opts)
			}
		}
	}
	if selected {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(S+GridWidth*selectedX), float64(T+GridHeight*selectedY))
		screen.DrawImage(selectedImg, opts)
	}
	return nil
}
func main() {
	var err error
	bgimg, _, err = ebitenutil.NewImageFromFile("./bgimg.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	gimulImgs[GimulTypeGreenWang], _, err = ebitenutil.NewImageFromFile("./green_wang.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	gimulImgs[GimulTypeGreenJa], _, err = ebitenutil.NewImageFromFile("./green_ja.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	gimulImgs[GimulTypeGreenJang], _, err = ebitenutil.NewImageFromFile("./green_jang.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	gimulImgs[GimulTypeGreenSang], _, err = ebitenutil.NewImageFromFile("./green_sang.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	gimulImgs[GimulTypeGreenQueen], _, err = ebitenutil.NewImageFromFile("./green_queen.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	//
	gimulImgs[GimulTypeRedWang], _, err = ebitenutil.NewImageFromFile("./red_wang.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	gimulImgs[GimulTypeRedJa], _, err = ebitenutil.NewImageFromFile("./red_ja.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	gimulImgs[GimulTypeRedJang], _, err = ebitenutil.NewImageFromFile("./red_jang.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	gimulImgs[GimulTypeRedSang], _, err = ebitenutil.NewImageFromFile("./red_sang.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	gimulImgs[GimulTypeRedQueen], _, err = ebitenutil.NewImageFromFile("./red_queen.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	selectedImg, _, err = ebitenutil.NewImageFromFile("./selected.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	//intialize board
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			board[i][j] = GimulTypeNone
		}
	}
	board[0][0] = GimulTypeGreenSang
	board[0][1] = GimulTypeGreenWang
	board[0][2] = GimulTypeGreenJang
	board[1][1] = GimulTypeGreenJa

	board[3][0] = GimulTypeRedSang
	board[3][1] = GimulTypeRedWang
	board[3][2] = GimulTypeRedJang
	board[2][1] = GimulTypeRedJa

	ebiten.Run(update, ScreenWidth, ScreenHeight, 1.0, "12 ")
	if err != nil {
		log.Fatal("ebiten run error : %v", err)
	}
}
