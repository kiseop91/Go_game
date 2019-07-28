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
	GimulTypeGreenJa   //1
	GimulTypeGreenJang //2
	GimulTypeGreenSang //3
	GimulTypeGreenQueen
	GimulTypeRedWang
	GimulTypeRedJa   //6
	GimulTypeRedJang //7
	GimulTypeRedSang //8
	GimulTypeRedQueen
	GimulTypeMax
)

const (
	ScreenWidth  = 480
	ScreenHeight = 462 //362
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
	deadboard      [6]GimulType
	board          [BoardWidth][BoardHeight]GimulType
	bgimg          *ebiten.Image
	gimulImgs      [GimulTypeMax]*ebiten.Image
	mgimulImgs     [GimulTypeMax - 1]*ebiten.Image
	selected       bool
	selectedX      int
	selectedY      int
	selectedImg    *ebiten.Image
	currentTeam    TeamType = TeamGreen
	gameover       bool
	GreenJaToQueen bool
	RedJaToQueen   bool
	idx            int
	cnt            int
	deadsel        bool
	deadtar        GimulType
)

func initiate() {
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
}

func GetTeamType(gimulType GimulType) TeamType {
	if gimulType == GimulTypeGreenJa ||
		gimulType == GimulTypeGreenJang ||
		gimulType == GimulTypeGreenWang ||
		gimulType == GimulTypeGreenSang ||
		gimulType == GimulTypeGreenQueen {
		return TeamGreen
	}
	if gimulType == GimulTypeRedJa ||
		gimulType == GimulTypeRedJang ||
		gimulType == GimulTypeRedWang ||
		gimulType == GimulTypeRedSang ||
		gimulType == GimulTypeRedQueen {
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
		//initiate()
	}
}

func move(prevX, prevY, tarX, tarY int) {
	if isMoveble(prevX, prevY, tarX, tarY) {
		OnDie(board[tarX][tarY])
		if board[tarX][tarY] != GimulTypeNone {
			if board[tarX][tarY] > 4 {
				deadboard[cnt] = board[tarX][tarY] - 5
			} else {
				deadboard[cnt] = board[tarX][tarY] + 5
			}
			cnt += 1
			cnt %= 6
		}
		board[prevX][prevY], board[tarX][tarY] = GimulTypeNone, board[prevX][prevY]
		selected = false
		if currentTeam == TeamGreen {
			currentTeam = TeamRed
		} else {
			currentTeam = TeamGreen
		}
		if tarX == 3 && board[tarX][tarY] == GimulTypeGreenJa {
			board[tarX][tarY] = GimulTypeGreenQueen
		}
		if tarX == 0 && board[tarX][tarY] == GimulTypeRedJa {
			board[tarX][tarY] = GimulTypeRedQueen
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
	case GimulTypeGreenQueen:
		return !(prevY-1 == tarY && prevX-1 == tarX || prevY+1 == tarY && prevX-1 == tarX)
	case GimulTypeRedQueen:
		return !(prevY+1 == tarY && prevX+1 == tarX || prevY-1 == tarY && prevX+1 == tarX)
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
		if y > 400 {
			if x > 0 && x < 40 {
				idx = 0
				deadtar = deadboard[idx]
				deadsel = true
			}
			if x > 75 && x < 90 {
				idx = 1
				deadtar = deadboard[idx]
				deadsel = true
			}
			if x > 125 && x < 140 {
				idx = 2
				deadtar = deadboard[idx]
				deadsel = true
			}
			//
			if x > 330 && x < 350 {
				idx = 3
				deadtar = deadboard[idx]
				deadsel = true
			}
			if x > 385 && x < 395 {
				idx = 4
				deadtar = deadboard[idx]
				deadsel = true
			}
			if x > 435 && x < 450 {
				idx = 5
				deadtar = deadboard[idx]
				deadsel = true
			}
		} else if i >= 0 && i < GridWidth && j >= 0 && j < GridHeight {
			if deadsel {
				if GetTeamType(deadtar) == TeamGreen {
					currentTeam = TeamRed
				} else {
					currentTeam = TeamGreen
				}
				board[i][j] = deadtar
				deadsel = false
				deadboard[idx] = GimulTypeNone
			}
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

	for i := 0; i < 6; i++ {
		opts := &ebiten.DrawImageOptions{}
		trim := 0
		if i >= 3 {
			trim = 160
		}
		if deadboard[i] == GimulTypeGreenQueen {
			deadboard[i] = GimulTypeGreenJa
		}
		if deadboard[i] == GimulTypeRedQueen {
			deadboard[i] = GimulTypeRedJa
		}
		opts.GeoM.Translate(float64(10+50.0*i+trim), 400)
		switch deadboard[i] {
		case GimulTypeGreenJa:
			screen.DrawImage(mgimulImgs[GimulTypeGreenJa], opts)
		case GimulTypeGreenSang:
			screen.DrawImage(mgimulImgs[GimulTypeGreenSang], opts)
		case GimulTypeGreenJang:
			screen.DrawImage(mgimulImgs[GimulTypeGreenJang], opts)
			//
		case GimulTypeRedJa:
			screen.DrawImage(mgimulImgs[GimulTypeRedJa], opts)
		case GimulTypeRedSang:
			screen.DrawImage(mgimulImgs[GimulTypeRedSang], opts)
		case GimulTypeRedJang:
			screen.DrawImage(mgimulImgs[GimulTypeRedJang], opts)

		}
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

	///
	mgimulImgs[GimulTypeGreenJa], _, err = ebitenutil.NewImageFromFile("./m/mgreen_ja.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	mgimulImgs[GimulTypeGreenJang], _, err = ebitenutil.NewImageFromFile("./m/mgreen_jang.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	mgimulImgs[GimulTypeGreenSang], _, err = ebitenutil.NewImageFromFile("./m/mgreen_sang.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	//

	mgimulImgs[GimulTypeRedJa], _, err = ebitenutil.NewImageFromFile("./m/mred_ja.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	mgimulImgs[GimulTypeRedJang], _, err = ebitenutil.NewImageFromFile("./m/mred_jang.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	mgimulImgs[GimulTypeRedSang], _, err = ebitenutil.NewImageFromFile("./m/mred_sang.PNG", ebiten.FilterDefault)
	if err != nil {
		log.Fatalf("read file error : %v ", err)
	}

	//	initailize deadboard
	for i := 0; i < 6; i++ {
		deadboard[i] = GimulTypeNone
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
