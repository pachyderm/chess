package pgn

import "fmt"

type Rank byte

const (
	NoRank Rank = '0' + iota
	Rank1
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)

type File byte

const (
	FileA File = 'a' + iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
	NoFile File = ' '
)

type Position uint64

const (
	A1 Position = 1 << iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1
	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8

	NoPosition Position = 0
)

func ParsePosition(pstr string) (Position, error) {
	p, ok := parsePosition(pstr)
	if !ok {
		return 0, fmt.Errorf("pgn: invalid position string: %s", pstr)
	}
	return p, nil
}

func parsePosition(pstr string) (Position, bool) {
	if len(pstr) != 2 {
		return 0, false
	}

	file := File(pstr[0])
	switch file {
	case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h':
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H':
		file += 'a' - 'A' // lowercase
	default:
		return 0, false
	}

	rank := Rank(pstr[1])
	switch rank {
	case '1', '2', '3', '4', '5', '6', '7', '8':
	default:
		return 0, false
	}

	p := PositionFromFileRank(file, rank)

	return p, true
}

func (p Position) String() string {
	if NoPosition == p {
		return "-"
	}
	f := byte(p.GetFile())
	r := byte(p.GetRank())
	return string([]byte{f, r})
}

func (p Position) GetRank() Rank {
	switch {
	case p == NoPosition:
		return NoRank
	case p <= H1:
		return Rank1
	case p <= H2:
		return Rank2
	case p <= H3:
		return Rank3
	case p <= H4:
		return Rank4
	case p <= H5:
		return Rank5
	case p <= H6:
		return Rank6
	case p <= H7:
		return Rank7
	case p <= H8:
		return Rank8
	}
	panic("unreachable")
}

func (p Position) GetFile() File {
	switch p % 255 {
	case 1 << 7:
		return FileH
	case 1 << 6:
		return FileG
	case 1 << 5:
		return FileF
	case 1 << 4:
		return FileE
	case 1 << 3:
		return FileD
	case 1 << 2:
		return FileC
	case 1 << 1:
		return FileB
	case 1 << 0:
		return FileA
	}
	return NoFile
}

func PositionFromFileRank(f File, r Rank) Position {
	// shift ['a'..'h'] and ['1'..'8'] to [0..7]
	f -= FileA
	r -= Rank1
	if f > 7 || r > 7 {
		return NoPosition
	}
	return Position(1) << (uint(r*8) + uint(f))
}
