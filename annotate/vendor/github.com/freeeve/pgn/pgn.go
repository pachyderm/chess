package pgn

import (
	"fmt"
	"io"
	"strings"
	"text/scanner"
)

type PGNScanner struct {
	s scanner.Scanner
}

type Game struct {
	Moves []Move
	Tags  map[string]string
}

type Move struct {
	From    Position
	To      Position
	Promote Piece
}

func (m Move) String() string {
	if m.Promote == NoPiece {
		return fmt.Sprintf("%v%v", m.From, m.To)
	}
	return fmt.Sprintf("%v%v%v", m.From, m.To, m.Promote)
}

var (
	NilMove Move = Move{From: NoPosition, To: NoPosition}
)

func ParseGame(s *scanner.Scanner) (*Game, error) {
	g := Game{Tags: map[string]string{}, Moves: []Move{}}
	err := ParseTags(s, &g)
	if err != nil {
		return nil, err
	}
	err = ParseMoves(s, &g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func ParseTags(s *scanner.Scanner, g *Game) error {
	//fmt.Println("starting tags parse")
	run := s.Peek()
	for run != scanner.EOF {
		switch run {
		case '[', ']', '\n', '\r':
			run = s.Next()
		case '1':
			return nil
		default:
			s.Scan()
			tag := s.TokenText()
			s.Scan()
			val := s.TokenText()
			//fmt.Println("tag:", tag, "; val:", val)
			g.Tags[tag] = strings.Trim(val, "\"")
		}
		run = s.Peek()
	}
	return nil
}

func isEnd(str string) bool {
	if str == "1/2-1/2" {
		return true
	}
	if str == "0-1" {
		return true
	}
	if str == "1-0" {
		return true
	}
	return false
}

func ParseMoves(s *scanner.Scanner, g *Game) error {
	//fmt.Println("starting moves parse")
	s.Mode = scanner.ScanIdents | scanner.ScanChars | scanner.ScanInts | scanner.ScanStrings
	run := s.Peek()
	board := NewBoard()
	var err error
	if len(g.Tags["FEN"]) > 0 {
		board, err = NewBoardFEN(g.Tags["FEN"])
		if err != nil {
			return err
		}
	}
	num := ""
	white := ""
	black := ""
	for run != scanner.EOF {
		switch run {
		case '(':
			for run != ')' && run != scanner.EOF {
				run = s.Next()
			}
		case '{':
			for run != '}' && run != scanner.EOF {
				run = s.Next()
			}
		case '#', '.', '+', '!', '?', '\n', '\r':
			run = s.Next()
			run = s.Peek()
		default:
			s.Scan()
			if s.TokenText() == "{" {
				run = '{'
				continue
			}
			if num == "" {
				num = s.TokenText()
				for s.Peek() == '-' {
					s.Scan()
					num += s.TokenText()
					s.Scan()
					num += s.TokenText()
				}
				for s.Peek() == '/' {
					s.Scan()
					num += s.TokenText()
					s.Scan()
					num += s.TokenText()
					s.Scan()
					num += s.TokenText()
					s.Scan()
					num += s.TokenText()
					s.Scan()
					num += s.TokenText()
					s.Scan()
					num += s.TokenText()
				}
				if isEnd(num) {
					return nil
				}
			} else if white == "" {
				white = s.TokenText()
				for s.Peek() == '-' {
					s.Scan()
					white += s.TokenText()
					s.Scan()
					white += s.TokenText()
				}
				for s.Peek() == '/' {
					s.Scan()
					white += s.TokenText()
					s.Scan()
					white += s.TokenText()
					s.Scan()
					white += s.TokenText()
					s.Scan()
					white += s.TokenText()
					s.Scan()
					white += s.TokenText()
					s.Scan()
					white += s.TokenText()
				}
				if isEnd(white) {
					return nil
				}
				if s.Peek() == '=' {
					s.Scan()
					white += s.TokenText()
					s.Scan()
					white += s.TokenText()
				}
				move, err := board.MoveFromAlgebraic(white, White)
				if err != nil {
					fmt.Println(board)
					return err
				}
				g.Moves = append(g.Moves, move)
				board.MakeMove(move)
			} else if black == "" {
				black = s.TokenText()
				for s.Peek() == '-' {
					s.Scan()
					black += s.TokenText()
					s.Scan()
					black += s.TokenText()
				}
				for s.Peek() == '/' {
					s.Scan()
					black += s.TokenText()
					s.Scan()
					black += s.TokenText()
					s.Scan()
					black += s.TokenText()
					s.Scan()
					black += s.TokenText()
					s.Scan()
					black += s.TokenText()
					s.Scan()
					black += s.TokenText()
				}
				if isEnd(black) {
					return nil
				}
				if s.Peek() == '=' {
					s.Scan()
					black += s.TokenText()
					s.Scan()
					black += s.TokenText()
				}
				move, err := board.MoveFromAlgebraic(black, Black)
				if err != nil {
					fmt.Println(board)
					return err
				}
				g.Moves = append(g.Moves, move)
				board.MakeMove(move)
				num = ""
				white = ""
				black = ""
			}
			run = s.Peek()
		}
	}
	return nil
}

func NewPGNScanner(r io.Reader) *PGNScanner {
	s := scanner.Scanner{}
	s.Init(r)
	return &PGNScanner{s: s}
}

func (ps *PGNScanner) Next() bool {
	if ps.s.Peek() == scanner.EOF {
		return false
	}
	return true
}

func (ps *PGNScanner) Scan() (*Game, error) {
	game, err := ParseGame(&ps.s)
	//fmt.Println(game)
	return game, err
}
