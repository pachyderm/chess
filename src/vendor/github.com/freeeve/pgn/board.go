package pgn

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrAmbiguousMove       = errors.New("pgn: ambiguous algebraic move")
	ErrUnknownMove         = errors.New("pgn: unknown move")
	ErrAttackerNotFound    = errors.New("pgn: attacker not found")
	ErrMoveFromEmptySquare = errors.New("pgn: move from empty square")
	ErrMoveWrongColor      = errors.New("pgn: move from wrong color")
	ErrMoveThroughPiece    = errors.New("pgn: move through piece")
	ErrMoveThroughCheck    = errors.New("pgn: move through check")
	ErrMoveIntoCheck       = errors.New("pgn: move into check")
	ErrMoveInvalidCastle   = errors.New("pgn: move invalid castle")
)

type Board struct {
	wPawns        uint64
	bPawns        uint64
	wRooks        uint64
	bRooks        uint64
	wKnights      uint64
	bKnights      uint64
	wBishops      uint64
	bBishops      uint64
	wQueens       uint64
	bQueens       uint64
	wKings        uint64
	bKings        uint64
	lastMove      Move
	wCastle       CastleStatus
	bCastle       CastleStatus
	toMove        Color
	fullmove      int
	halfmoveClock int
}

type Piece byte

const (
	NoPiece     Piece = ' '
	BlackPawn   Piece = 'p'
	BlackKnight Piece = 'n'
	BlackBishop Piece = 'b'
	BlackRook   Piece = 'r'
	BlackQueen  Piece = 'q'
	BlackKing   Piece = 'k'
	WhitePawn   Piece = 'P'
	WhiteKnight Piece = 'N'
	WhiteBishop Piece = 'B'
	WhiteRook   Piece = 'R'
	WhiteQueen  Piece = 'Q'
	WhiteKing   Piece = 'K'
)

func (p Piece) Color() Color {
	if 'a' <= p && p <= 'z' {
		return Black
	}
	if 'A' <= p && p <= 'Z' {
		return White
	}
	return NoColor
}

func (p *Piece) Normalize() {
	*p = Piece(bytes.ToLower([]byte{byte(*p)})[0])
}

type Color int8

const (
	NoColor Color = iota
	Black
	White
)

func (c Color) String() string {
	if c == White {
		return "w"
	} else if c == Black {

		return "b"
	}
	return " "
}

type CastleStatus int8

const (
	Both CastleStatus = iota
	None
	Kingside
	Queenside
)

func (cs CastleStatus) String(c Color) string {
	type p struct {
		CastleStatus
		Color
	}
	switch (p{cs, c}) {
	case p{Both, Black}:
		return "kq"
	case p{Both, White}:
		return "KQ"
	case p{Kingside, Black}:
		return "k"
	case p{Kingside, White}:
		return "K"
	case p{Queenside, Black}:
		return "q"
	case p{Queenside, White}:
		return "Q"
	}
	if cs == None {
		return "-"
	}
	return ""
}

func (b *Board) MakeCoordMove(str string) error {
	move, err := MoveFromCoord(str)
	if err != nil {
		return err
	}
	err = b.MakeMove(move)
	if err != nil {
		return err
	}
	return nil
}

func MoveFromCoord(str string) (Move, error) {
	length := len(str)
	if length < 4 {
		return NilMove, ErrUnknownMove
	}
	// handle promotion
	promote := NoPiece
	if length == 5 {
		promote = Piece(str[length-1])
		promote.Normalize()
		str = str[:length-1]
	}
	fromPos, err := ParsePosition(str[:2])
	if err != nil {
		return NilMove, ErrUnknownMove
	}
	toPos, err := ParsePosition(str[2:])
	if err != nil {
		return NilMove, ErrUnknownMove
	}
	return Move{fromPos, toPos, promote}, nil
}

func (b *Board) MakeAlgebraicMove(str string, color Color) error {
	move, err := b.MoveFromAlgebraic(str, color)
	if err != nil {
		return err
	}
	err = b.MakeMove(move)
	if err != nil {
		return err
	}
	return nil
}

func (b *Board) MoveFromAlgebraic(str string, color Color) (Move, error) {
	str = strings.Trim(str, "+!?#")
	//fmt.Println("move from alg:", str, "..", color)
	if b.toMove != color {
		return NilMove, ErrMoveWrongColor
	}
	// handle promotion
	promote := NoPiece
	if strings.Contains(str, "=") {
		promote = Piece(str[len(str)-1])
		promote.Normalize()
		str = str[:len(str)-2]
	}
	pos, err := ParsePosition(str)
	testPos := pos
	if err == nil {
		for testPos != NoPosition {
			r := testPos.GetRank()
			f := testPos.GetFile()
			// if it's a raw position, it's a pawn move
			if color == White {
				testPos = PositionFromFileRank(f, r-1)
				if b.GetPiece(testPos) == WhitePawn {
					return Move{testPos, pos, promote}, nil
				}
				if pos == NoPosition {
					return NilMove, fmt.Errorf("Position out of bounds")
				}
			} else {
				testPos = PositionFromFileRank(f, r+1)
				if b.GetPiece(testPos) == BlackPawn {
					return Move{testPos, pos, promote}, nil
				}
				if pos == NoPosition {
					return NilMove, fmt.Errorf("Position out of bounds")
				}
			}
		}
	} else {
		// otherwise it's a non-pawn move (or pawn take)
		switch str[0] {
		case 'O':
			if str == "O-O" {
				return b.getKingsideCastle(color)
			} else if str == "O-O-O" {
				return b.getQueensideCastle(color)
			} else {
				return NilMove, ErrUnknownMove
			}
		case 'N':
			pos, err := ParsePosition(str[len(str)-2 : len(str)])
			if err != nil {
				return NilMove, err
			}
			fromPos, err := b.findAttackingKnight(pos, color, true)
			if err == ErrAmbiguousMove {
				if str[1] >= 'a' && str[1] <= 'h' {
					fromPos, err = b.findAttackingKnightFromFile(pos, color, File(str[1]))
					if err == nil {
						return Move{fromPos, pos, NoPiece}, nil
					}
				} else if str[1] >= '1' && str[1] <= '8' {
					fromPos, err = b.findAttackingKnightFromRank(pos, color, Rank(str[1]))
					if err == nil {
						return Move{fromPos, pos, NoPiece}, nil
					}
				}
				return NilMove, ErrAmbiguousMove
			}
			if err != nil {
				return NilMove, err
			}
			return Move{fromPos, pos, NoPiece}, nil
		case 'B':
			pos, err := ParsePosition(str[len(str)-2 : len(str)])
			if err != nil {
				return NilMove, err
			}
			fromPos, err := b.findAttackingBishop(pos, color, true)
			// TODO handle rare ambiguous move
			if err != nil {
				return NilMove, err
			}
			return Move{fromPos, pos, NoPiece}, nil
		case 'R':
			pos, err := ParsePosition(str[len(str)-2 : len(str)])
			if err != nil {
				return NilMove, err
			}
			fromPos, err := b.findAttackingRook(pos, color, true)
			if err == ErrAmbiguousMove {
				if str[1] >= 'a' && str[1] <= 'h' {
					fromPos, err = b.findAttackingRookFromFile(pos, color, File(str[1]))
					if err == nil {
						return Move{fromPos, pos, NoPiece}, nil
					}
				} else if str[1] >= '1' && str[1] <= '8' {
					fromPos, err = b.findAttackingRookFromRank(pos, color, Rank(str[1]))
					if err == nil {
						return Move{fromPos, pos, NoPiece}, nil
					}
				}
			}
			if err != nil {
				return NilMove, err
			}
			return Move{fromPos, pos, NoPiece}, nil
		case 'Q':
			pos, err := ParsePosition(str[len(str)-2 : len(str)])
			if err != nil {
				return NilMove, err
			}

			fromPos, err := b.findAttackingQueen(pos, color, true)

			if err == ErrAmbiguousMove {
				if str[1] >= 'a' && str[1] <= 'h' {
					fromPos, err = b.findAttackingQueenFromFile(pos, color, File(str[1]))
					if err == nil {
						return Move{fromPos, pos, NoPiece}, nil
					}
				} else if str[1] >= '1' && str[1] <= '8' {
					fromPos, err = b.findAttackingQueenFromRank(pos, color, Rank(str[1]))
					if err == nil {
						return Move{fromPos, pos, NoPiece}, nil
					}
				}
			}

			if err != nil {
				return NilMove, err
			}
			return Move{fromPos, pos, NoPiece}, nil
		case 'K':
			pos, err := ParsePosition(str[len(str)-2 : len(str)])
			if err != nil {
				return NilMove, err
			}
			fromPos, err := b.findAttackingKing(pos, color)
			if err != nil {
				return NilMove, err
			}
			return Move{fromPos, pos, NoPiece}, nil
		}
		// pawn taking move
		if str[0] >= 'a' && str[0] <= 'h' && str[1] == 'x' {
			pos, err := ParsePosition(str[len(str)-2 : len(str)])
			if err != nil {
				return NilMove, err
			}
			fromPos, err := b.findAttackingPawn(pos, color, true)
			if err == ErrAmbiguousMove {
				fromPos, err = b.findAttackingPawnFromFile(pos, color, File(str[0]))
				if err == nil {
					return Move{fromPos, pos, promote}, nil
				}
			}
			if err != nil {
				return NilMove, err
			}
			return Move{fromPos, pos, promote}, nil
		}
	}
	return NilMove, ErrUnknownMove
}

func NewBoard() *Board {
	b, _ := NewBoardFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	return b
}

func NewBoardFEN(fen string) (*Board, error) {
	f, err := ParseFEN(fen)
	if err != nil {
		return nil, err
	}
	b := Board{}
	b.toMove = f.ToMove
	b.wCastle = f.WhiteCastleStatus
	b.bCastle = f.BlackCastleStatus
	b.fullmove = f.Fullmove
	b.halfmoveClock = f.HalfmoveClock

	x := byte('a')
	y := byte('8')
	for i := 0; i < len(f.FOR); i++ {
		// if we're at the end of the row
		if f.FOR[i] == '/' {
			x = 'a'
			y--
			continue
		} else if f.FOR[i] >= '1' && f.FOR[i] <= '8' {
			// if we have blank squares
			j, err := strconv.Atoi(string(f.FOR[i]))
			if err != nil {
				fmt.Println(err)
			}
			x += byte(j)
			continue
		} else {
			// if we have a piece
			pos, err := ParsePosition(fmt.Sprintf("%c%c", x, y))
			if err != nil {
				fmt.Println(err)
			}
			b.SetPiece(pos, Piece(f.FOR[i]))
			x++
		}
	}
	return &b, nil
}

func (b *Board) String() string {
	return FENFromBoard(b).String()
}

func (b *Board) MakeMove(m Move) error {
	p := b.GetPiece(m.From)
	if p == NoPiece {
		return ErrMoveFromEmptySquare
	}

	// sanity check for wrong color
	if p.Color() != b.toMove {
		return ErrMoveWrongColor
	}

	// see if we're taking a piece
	take := b.GetPiece(m.To)
	if take != NoPiece {
		b.RemovePiece(m.To, take)
	}

	// the moving
	b.SetPiece(m.To, p)
	b.RemovePiece(m.From, p)

	// handle special cases
	switch p {
	case BlackRook:
		switch m.From {
		case A8:
			switch b.bCastle {
			case Both:
				b.bCastle = Kingside
			case Queenside:
				b.bCastle = None
			}
		case H8:
			switch b.bCastle {
			case Both:
				b.bCastle = Queenside
			case Kingside:
				b.bCastle = None
			}
		}
	case BlackPawn:
		if m.From.GetFile() != m.To.GetFile() &&
			take == NoPiece {
			b.RemovePiece(PositionFromFileRank(m.To.GetFile(), m.To.GetRank()+1), WhitePawn)
		}
		if m.Promote != NoPiece {
			b.SetPiece(m.To, m.Promote)
			b.RemovePiece(m.To, BlackPawn)
		}
	case BlackKing:
		// handle castles
		if m.From == E8 && m.To == G8 {
			if b.bCastle != Kingside && b.bCastle != Both {
				return ErrMoveInvalidCastle
			}
			rook := b.GetPiece(H8)
			b.RemovePiece(H8, rook)
			b.SetPiece(F8, rook)
		} else if m.From == E8 && m.To == C8 {
			if b.bCastle != Queenside && b.bCastle != Both {
				return ErrMoveInvalidCastle
			}
			rook := b.GetPiece(A8)
			b.RemovePiece(A8, rook)
			b.SetPiece(D8, rook)
		}
		b.bCastle = None
	case WhiteRook:
		switch m.From {
		case A1:
			switch b.wCastle {
			case Both:
				b.wCastle = Kingside
			case Queenside:
				b.wCastle = None
			}
		case H1:
			switch b.wCastle {
			case Both:
				b.wCastle = Queenside
			case Kingside:
				b.wCastle = None
			}
		}
	case WhitePawn:
		if m.From.GetFile() != m.To.GetFile() &&
			take == NoPiece {
			b.RemovePiece(PositionFromFileRank(m.To.GetFile(), m.To.GetRank()-1), BlackPawn)
		}
		if m.Promote != NoPiece {
			// TODO refactor this. semi hacky
			b.RemovePiece(m.To, WhitePawn)
			switch m.Promote {
			case BlackQueen:
				b.SetPiece(m.To, WhiteQueen)
			case BlackRook:
				b.SetPiece(m.To, WhiteRook)
			case BlackBishop:
				b.SetPiece(m.To, WhiteBishop)
			case BlackKnight:
				b.SetPiece(m.To, WhiteKnight)
			}
		}
	case WhiteKing:
		// handle castles
		if m.From == E1 && m.To == G1 {
			if b.wCastle != Kingside && b.wCastle != Both {
				return ErrMoveInvalidCastle
			}
			rook := b.GetPiece(H1)
			b.RemovePiece(H1, rook)
			b.SetPiece(F1, rook)
		} else if m.From == E1 && m.To == C1 {
			if b.wCastle != Queenside && b.wCastle != Both {
				return ErrMoveInvalidCastle
			}
			rook := b.GetPiece(A1)
			b.RemovePiece(A1, rook)
			b.SetPiece(D1, rook)
		}
		b.wCastle = None
	}

	// swap next color
	switch b.toMove {
	case White:
		b.toMove = Black
	case Black:
		b.toMove = White
	}

	// handle move number increment
	if b.toMove == White {
		b.fullmove++
	}

	// handle halfmove clock
	if take != NoPiece || p == WhitePawn || p == BlackPawn {
		b.halfmoveClock = 0
	} else {
		b.halfmoveClock++
	}

	// save lastMove
	b.lastMove = m

	return nil
}

// refPiece returns a pointer to the Board field corresponding to p
// fallback just simplifies code to avoid nil checking
// (a dummy value will be passed in if p is invalid)
func (b *Board) refPiece(p Piece, fallback *uint64) *uint64 {
	switch p {
	case WhitePawn:
		return &b.wPawns
	case BlackPawn:
		return &b.bPawns
	case WhiteKnight:
		return &b.wKnights
	case BlackKnight:
		return &b.bKnights
	case WhiteBishop:
		return &b.wBishops
	case BlackBishop:
		return &b.bBishops
	case WhiteRook:
		return &b.wRooks
	case BlackRook:
		return &b.bRooks
	case WhiteQueen:
		return &b.wQueens
	case BlackQueen:
		return &b.bQueens
	case WhiteKing:
		return &b.wKings
	case BlackKing:
		return &b.bKings
	}
	return fallback
}

func (b *Board) RemovePiece(pos Position, p Piece) {
	npos := ^uint64(pos) // negation of pos
	pp := b.refPiece(p, &npos)
	*pp &= npos
}

func (b *Board) SetPiece(pos Position, p Piece) {
	cpos := uint64(pos) // (nominative) conversion of pos
	pp := b.refPiece(p, &cpos)
	*pp |= cpos
}

func (b Board) GetPiece(p Position) Piece {
	q := uint64(p)
	switch {
	case q&b.bPawns != 0:
		return BlackPawn
	case q&b.wPawns != 0:
		return WhitePawn
	case q&b.bKnights != 0:
		return BlackKnight
	case q&b.wKnights != 0:
		return WhiteKnight
	case q&b.bBishops != 0:
		return BlackBishop
	case q&b.wBishops != 0:
		return WhiteBishop
	case q&b.bRooks != 0:
		return BlackRook
	case q&b.wRooks != 0:
		return WhiteRook
	case q&b.bQueens != 0:
		return BlackQueen
	case q&b.wQueens != 0:
		return WhiteQueen
	case q&b.bKings != 0:
		return BlackKing
	case q&b.wKings != 0:
		return WhiteKing
	}
	return NoPiece
}

func (b Board) findAttackingPawn(pos Position, color Color, check bool) (Position, error) {
	retPos := NoPosition
	count := 0
	switch color {
	case White:
		// special en-passant case
		if b.lastMove.To.GetFile() == pos.GetFile() &&
			b.lastMove.To.GetRank() == pos.GetRank()-1 &&
			b.lastMove.From.GetRank() == pos.GetRank()+1 &&
			b.GetPiece(PositionFromFileRank(pos.GetFile(), pos.GetRank()-1)) == BlackPawn {
			testPos := PositionFromFileRank(pos.GetFile()+1, pos.GetRank()-1)
			if b.GetPiece(testPos) == WhitePawn &&
				(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
				retPos = testPos
				count++
				// TODO remove these breaks
				break
			}
			testPos = PositionFromFileRank(pos.GetFile()-1, pos.GetRank()-1)
			if b.GetPiece(testPos) == WhitePawn &&
				(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
				retPos = testPos
				count++
				break
			}
		}
		testPos := PositionFromFileRank(pos.GetFile()+1, pos.GetRank()-1)
		if b.GetPiece(testPos) == WhitePawn &&
			(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		}
		testPos = PositionFromFileRank(pos.GetFile()-1, pos.GetRank()-1)
		if b.GetPiece(testPos) == WhitePawn &&
			(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		}
	case Black:
		// special en-passant case
		if b.lastMove.To.GetFile() == pos.GetFile() &&
			b.lastMove.To.GetRank() == pos.GetRank()+1 &&
			b.lastMove.From.GetRank() == pos.GetRank()-1 &&
			b.GetPiece(PositionFromFileRank(pos.GetFile(), pos.GetRank()+1)) == WhitePawn {
			testPos := PositionFromFileRank(pos.GetFile()+1, pos.GetRank()+1)
			if b.GetPiece(testPos) == BlackPawn &&
				(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
				retPos = testPos
				count++
				break
			}
			testPos = PositionFromFileRank(pos.GetFile()-1, pos.GetRank()+1)
			if b.GetPiece(testPos) == BlackPawn &&
				(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
				retPos = testPos
				count++
				break
			}
		}
		testPos := PositionFromFileRank(pos.GetFile()+1, pos.GetRank()+1)
		if b.GetPiece(testPos) == BlackPawn &&
			(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		}
		testPos = PositionFromFileRank(pos.GetFile()-1, pos.GetRank()+1)
		if b.GetPiece(testPos) == BlackPawn &&
			(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		}
	}

	if count > 1 {
		return NoPosition, ErrAmbiguousMove
	}
	if retPos == NoPosition {
		return retPos, ErrAttackerNotFound
	}
	return retPos, nil
}

func (b Board) findAttackingPawnFromFile(pos Position, color Color, file File) (Position, error) {
	retPos := NoPosition
	count := 0
	if color == White {
		// special en-passant case
		if b.lastMove.To.GetFile() == pos.GetFile() &&
			b.lastMove.To.GetRank() == pos.GetRank()-1 &&
			b.lastMove.From.GetRank() == pos.GetRank()+1 &&
			b.GetPiece(PositionFromFileRank(pos.GetFile(), pos.GetRank()-1)) == BlackPawn {
			if b.GetPiece(PositionFromFileRank(file, pos.GetRank()-1)) == WhitePawn {
				retPos = PositionFromFileRank(file, pos.GetRank()-1)
				count++
			}
		}
		if b.GetPiece(PositionFromFileRank(file, pos.GetRank()-1)) == WhitePawn {
			retPos = PositionFromFileRank(file, pos.GetRank()-1)
			count++
		}
	} else {
		// special en-passant case
		if b.lastMove.To.GetFile() == pos.GetFile() &&
			b.lastMove.To.GetRank() == pos.GetRank()+1 &&
			b.lastMove.From.GetRank() == pos.GetRank()-1 &&
			b.GetPiece(PositionFromFileRank(pos.GetFile(), pos.GetRank()+1)) == WhitePawn {
			if b.GetPiece(PositionFromFileRank(file, pos.GetRank()+1)) == BlackPawn {
				retPos = PositionFromFileRank(file, pos.GetRank()+1)
				count++
			}
		}
		if b.GetPiece(PositionFromFileRank(file, pos.GetRank()+1)) == BlackPawn {
			retPos = PositionFromFileRank(file, pos.GetRank()+1)
			count++
		}
	}
	if count > 1 {
		return NoPosition, ErrAmbiguousMove
	}
	if retPos == NoPosition {
		return retPos, ErrAttackerNotFound
	}
	return retPos, nil
}

func (b Board) findAttackingBishop(pos Position, color Color, check bool) (Position, error) {
	count := 0
	retPos := NoPosition

	r := pos.GetRank()
	f := pos.GetFile()
	for {
		f--
		r--
		testPos := PositionFromFileRank(f, r)
		if b.checkBishopColor(testPos, color) &&
			(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		f--
		r++
		testPos := PositionFromFileRank(f, r)
		if b.checkBishopColor(testPos, color) &&
			(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		f++
		r++
		testPos := PositionFromFileRank(f, r)
		if b.checkBishopColor(testPos, color) &&
			(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	r = pos.GetRank()
	f = pos.GetFile()
	for {
		f++
		r--
		testPos := PositionFromFileRank(f, r)
		if b.checkBishopColor(testPos, color) &&
			(!check || !b.moveIntoCheck(Move{testPos, pos, NoPiece}, color)) {
			retPos = testPos
			count++
		} else if testPos == NoPosition || b.containsPieceAt(testPos) {
			break
		}
	}

	if count > 1 {
		return NoPosition, ErrAmbiguousMove
	}
	if count == 0 {
		return NoPosition, ErrAttackerNotFound
	}
	return retPos, nil
}

func (b Board) checkBishopColor(pos Position, color Color) bool {
	return (b.GetPiece(pos) == WhiteBishop && color == White) ||
		(b.GetPiece(pos) == BlackBishop && color == Black)
}

func (b Board) findAttackingKing(pos Position, color Color) (Position, error) {
	count := 0
	retPos := NoPosition

	// straight
	r := pos.GetRank()
	f := pos.GetFile()
	f--
	testPos := PositionFromFileRank(f, r)
	if b.checkKingColor(testPos, color) {
		retPos = testPos
		count++
	}

	r = pos.GetRank()
	f = pos.GetFile()
	f++
	testPos = PositionFromFileRank(f, r)
	if b.checkKingColor(testPos, color) {
		retPos = testPos
		count++
	}

	r = pos.GetRank()
	f = pos.GetFile()
	r++
	testPos = PositionFromFileRank(f, r)
	if b.checkKingColor(testPos, color) {
		retPos = testPos
		count++
	}

	r = pos.GetRank()
	f = pos.GetFile()
	r--
	testPos = PositionFromFileRank(f, r)
	if b.checkKingColor(testPos, color) {
		retPos = testPos
		count++
	}

	// diagonals
	r = pos.GetRank()
	f = pos.GetFile()
	f--
	r--
	testPos = PositionFromFileRank(f, r)
	if b.checkKingColor(testPos, color) {
		retPos = testPos
		count++
	}

	r = pos.GetRank()
	f = pos.GetFile()
	f--
	r++
	testPos = PositionFromFileRank(f, r)
	if b.checkKingColor(testPos, color) {
		retPos = testPos
		count++
	}

	r = pos.GetRank()
	f = pos.GetFile()
	f++
	r++
	testPos = PositionFromFileRank(f, r)
	if b.checkKingColor(testPos, color) {
		retPos = testPos
		count++
	}

	r = pos.GetRank()
	f = pos.GetFile()
	f++
	r--
	testPos = PositionFromFileRank(f, r)
	if b.checkKingColor(testPos, color) {
		retPos = testPos
		count++
	}

	if count == 0 {
		return NoPosition, ErrAttackerNotFound
	}
	return retPos, nil
}

func (b Board) checkKingColor(pos Position, color Color) bool {
	return (b.GetPiece(pos) == WhiteKing && color == White) ||
		(b.GetPiece(pos) == BlackKing && color == Black)
}

func (b Board) containsPieceAt(pos Position) bool {
	return (uint64(b.wPawns)|uint64(b.bPawns)|uint64(b.wKnights)|uint64(b.bKnights)|
		uint64(b.wBishops)|uint64(b.bBishops)|uint64(b.wRooks)|uint64(b.bRooks)|
		uint64(b.wQueens)|uint64(b.bQueens)|uint64(b.wKings)|uint64(b.bKings))&uint64(pos) > 0
}

func (b Board) getKingsideCastle(color Color) (Move, error) {
	if color == White {
		if b.wCastle != Both && b.wCastle != Kingside {
			return NilMove, ErrMoveInvalidCastle
		}
		if b.containsPieceAt(F1) || b.containsPieceAt(G1) {
			return NilMove, ErrMoveThroughPiece
		}
		if b.positionAttackedBy(F1, Black) || b.positionAttackedBy(G1, Black) {
			return NilMove, ErrMoveThroughCheck
		}
		return Move{E1, G1, NoPiece}, nil
	} else {
		if b.bCastle != Both && b.bCastle != Kingside {
			return NilMove, ErrMoveInvalidCastle
		}
		if b.containsPieceAt(F8) || b.containsPieceAt(G8) {
			return NilMove, ErrMoveThroughPiece
		}
		if b.positionAttackedBy(F8, White) || b.positionAttackedBy(G8, White) {
			return NilMove, ErrMoveThroughCheck
		}
		return Move{E8, G8, NoPiece}, nil
	}
}

func (b Board) getQueensideCastle(color Color) (Move, error) {
	if color == White {
		if b.wCastle != Both && b.wCastle != Queenside {
			return NilMove, ErrMoveInvalidCastle
		}
		if b.containsPieceAt(B1) || b.containsPieceAt(C1) || b.containsPieceAt(D1) {
			return NilMove, ErrMoveThroughPiece
		}
		if b.positionAttackedBy(C1, Black) || b.positionAttackedBy(D1, Black) {
			return NilMove, ErrMoveThroughCheck
		}
		return Move{E1, C1, NoPiece}, nil
	} else {
		if b.bCastle != Both && b.bCastle != Queenside {
			return NilMove, ErrMoveInvalidCastle
		}
		if b.containsPieceAt(B8) || b.containsPieceAt(C8) || b.containsPieceAt(D8) {
			return NilMove, ErrMoveThroughPiece
		}
		if b.positionAttackedBy(C8, White) || b.positionAttackedBy(D8, White) {
			return NilMove, ErrMoveThroughCheck
		}
		return Move{E8, C8, NoPiece}, nil
	}
}

func (b Board) positionAttackedBy(pos Position, color Color) bool {
	p, err := b.findAttackingPawn(pos, color, false)
	if p != NoPosition || err == ErrAmbiguousMove {
		return true
	}
	p, err = b.findAttackingKnight(pos, color, false)
	if p != NoPosition || err == ErrAmbiguousMove {
		return true
	}
	p, err = b.findAttackingBishop(pos, color, false)
	if p != NoPosition || err == ErrAmbiguousMove {
		return true
	}
	p, err = b.findAttackingRook(pos, color, false)
	if p != NoPosition || err == ErrAmbiguousMove {
		return true
	}
	p, err = b.findAttackingQueen(pos, color, false)
	if p != NoPosition || err == ErrAmbiguousMove {
		return true
	}
	p, err = b.findAttackingKing(pos, color)
	if p != NoPosition || err == ErrAmbiguousMove {
		return true
	}
	return false
}

func (b Board) moveIntoCheck(move Move, color Color) bool {
	tempb := b
	tempb.MakeMove(move)
	if color == White {
		p := tempb.FindKing(White)
		return tempb.positionAttackedBy(p, Black)
	} else if color == Black {
		p := tempb.FindKing(Black)
		return tempb.positionAttackedBy(p, White)
	}
	return false
}

func (b Board) FindKing(color Color) Position {
	for r := Rank8; r >= Rank1; r-- {
		for f := FileA; f <= FileH; f++ {
			pos := PositionFromFileRank(f, r)
			p := b.GetPiece(pos)
			if (p == BlackKing || p == WhiteKing) && p.Color() == color {
				return pos
			}
		}
	}
	return NoPosition
}
