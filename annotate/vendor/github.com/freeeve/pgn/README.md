pgn
===

A pgn parser for golang. This is release v1.

[![Build Status](https://travis-ci.org/freeeve/pgn.png?branch=master)](https://travis-ci.org/freeeve/pgn)
[![Coverage Status](https://coveralls.io/repos/freeeve/pgn/badge.svg?branch=v1&service=github)](https://coveralls.io/github/freeeve/pgn?branch=v1)

Normal go install... `go get gopkg.in/freeeve/pgn.v1`

## minimum viable snippet

```Go
package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/freeeve/pgn.v1"
)

func main() {
	f, err := os.Open("polgar.pgn")
	if err != nil {
		log.Fatal(err)
	}
	ps := pgn.NewPGNScanner(f)
	// while there's more to read in the file
	for ps.Next() {
		// scan the next game
		game, err := ps.Scan()
		if err != nil {
			log.Fatal(err)
		}

		// print out tags
		fmt.Println(game.Tags)

		// make a new board so we can get FEN positions
		b := pgn.NewBoard()
		for _, move := range game.Moves {
			// make the move on the board
			b.MakeMove(move)
			// print out FEN for each move in the game
			fmt.Println(b)
		}
	}
}
```

produces output like this for each game in the pgn file:

```
map[Event:Women's Chess Cup
    Site:Dresden GER 
    Round:7.1 
    Black:Polgar,Z 
    Result:1/2-1/2 
    Date:2006.07.08 
    White:Paehtz,E 
    WhiteElo:2438 
    BlackElo:2577 
    ECO:B35]
rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq d3 0 1
rnbqkbnr/pppp1ppp/4p3/8/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 2
rnbqkbnr/pppp1ppp/4p3/8/3P4/5N2/PPP1PPPP/RNBQKB1R b KQkq - 1 2
rnbqkbnr/pppp2pp/4p3/5p2/3P4/5N2/PPP1PPPP/RNBQKB1R w KQkq f6 0 3
rnbqkbnr/pppp2pp/4p3/3P1p2/8/5N2/PPP1PPPP/RNBQKB1R b KQkq - 0 3
rnbqkb1r/pppp2pp/4pn2/3P1p2/8/5N2/PPP1PPPP/RNBQKB1R w KQkq - 1 4
rnbqkb1r/pppp2pp/4Pn2/5p2/8/5N2/PPP1PPPP/RNBQKB1R b KQkq - 0 4
rnbqkb1r/ppp3pp/4pn2/5p2/8/5N2/PPP1PPPP/RNBQKB1R w KQkq - 0 5
rnbQkb1r/ppp3pp/4pn2/5p2/8/5N2/PPP1PPPP/RNB1KB1R b KQkq - 0 5
rnbk1b1r/ppp3pp/4pn2/5p2/8/5N2/PPP1PPPP/RNB1KB1R w KQ - 0 6
rnbk1b1r/ppp3pp/4pn2/5p2/8/2N2N2/PPP1PPPP/R1B1KB1R b KQ - 1 6
rnbk3r/ppp3pp/4pn2/5p2/1b6/2N2N2/PPP1PPPP/R1B1KB1R w KQ - 2 7
rnbk3r/ppp3pp/4pn2/5p2/1b6/2N2N2/PPPBPPPP/R3KB1R b KQ - 3 7
rnb4r/ppp1k1pp/4pn2/5p2/1b6/2N2N2/PPPBPPPP/R3KB1R w KQ - 4 8
rnb4r/ppp1k1pp/4pn2/5p2/1b6/P1N2N2/1PPBPPPP/R3KB1R b KQ - 0 8
rnb4r/ppp1k1pp/4pn2/2b2p2/8/P1N2N2/1PPBPPPP/R3KB1R w KQ - 1 9
rnb4r/ppp1k1pp/4pn2/2b2p2/1P6/P1N2N2/2PBPPPP/R3KB1R b KQ b3 0 9
rnb4r/ppp1k1pp/3bpn2/5p2/1P6/P1N2N2/2PBPPPP/R3KB1R w KQ - 1 10
rnb4r/ppp1k1pp/3bpn2/1N3p2/1P6/P4N2/2PBPPPP/R3KB1R b KQ - 2 10
rnb4r/ppp1k1pp/3bp3/1N3p2/1P2n3/P4N2/2PBPPPP/R3KB1R w KQ - 3 11
rnb4r/ppp1k1pp/3bp3/1N3p2/1P2n3/P3BN2/2P1PPPP/R3KB1R b KQ - 4 11
rnb4r/1pp1k1pp/3bp3/pN3p2/1P2n3/P3BN2/2P1PPPP/R3KB1R w KQ a6 0 12
rnb4r/1pp1k1pp/3Np3/p4p2/1P2n3/P3BN2/2P1PPPP/R3KB1R b KQ - 0 12
rnb4r/1p2k1pp/3pp3/p4p2/1P2n3/P3BN2/2P1PPPP/R3KB1R w KQ - 0 13
rnb4r/1p2k1pp/3pp3/P4p2/4n3/P3BN2/2P1PPPP/R3KB1R b KQ - 0 13
1nb4r/1p2k1pp/3pp3/r4p2/4n3/P3BN2/2P1PPPP/R3KB1R w KQ - 0 14
1nb4r/1p2k1pp/3pp3/r4p2/4n3/P3B3/2PNPPPP/R3KB1R b KQ - 1 14
1nb4r/1p2k1pp/3pp3/r4p2/8/P3B3/2PnPPPP/R3KB1R w KQ - 0 15
1nb4r/1p2k1pp/3pp3/r4p2/8/P7/2PBPPPP/R3KB1R b KQ - 0 15
1nb4r/1p2k1pp/3pp3/2r2p2/8/P7/2PBPPPP/R3KB1R w KQ - 1 16
1nb4r/1p2k1pp/3pp3/2r2p2/8/P7/2PBPPPP/2KR1B1R b - - 2 16
2b4r/1p2k1pp/2npp3/2r2p2/8/P7/2PBPPPP/2KR1B1R w - - 3 17
2b4r/1p2k1pp/2npp3/2r2p2/8/P3P3/2PB1PPP/2KR1B1R b - - 0 17
7r/1p1bk1pp/2npp3/2r2p2/8/P3P3/2PB1PPP/2KR1B1R w - - 1 18
7r/1p1bk1pp/2npp3/2r2p2/5P2/P3P3/2PB2PP/2KR1B1R b - f3 0 18
2r5/1p1bk1pp/2npp3/2r2p2/5P2/P3P3/2PB2PP/2KR1B1R w - - 1 19
2r5/1p1bk1pp/2npp3/2r2p2/5P2/P2BP3/2PB2PP/2KR3R b - - 2 19
2r5/1p1bk1pp/2np4/2r1pp2/5P2/P2BP3/2PB2PP/2KR3R w - - 0 20
2r5/1p1bk1pp/2np4/2r1pp2/4PP2/P2B4/2PB2PP/2KR3R b - - 0 20
2r5/1p1bk1pp/2np4/2r1p3/4pP2/P2B4/2PB2PP/2KR3R w - - 0 21
2r5/1p1bk1pp/2np4/2r1p3/4BP2/P7/2PB2PP/2KR3R b - - 0 21
2r5/1p1bk1pp/3p4/2r1p3/3nBP2/P7/2PB2PP/2KR3R w - - 1 22
2r5/1p1bk1pp/3p4/2r1P3/3nB3/P7/2PB2PP/2KR3R b - - 0 22
2r5/1p1bk1pp/3p4/4r3/3nB3/P7/2PB2PP/2KR3R w - - 0 23
2r5/1p1bk1pp/3p4/4r1B1/3nB3/P7/2P3PP/2KR3R b - - 1 23
2r5/1p1bk1pp/3p4/6r1/3nB3/P7/2P3PP/2KR3R w - - 0 24
2r5/1p1bk1pp/3p4/6r1/3RB3/P7/2P3PP/2K4R b - - 0 24
2r5/1p1bk2p/3p2p1/6r1/3RB3/P7/2P3PP/2K4R w - - 0 25
2r5/1B1bk2p/3p2p1/6r1/3R4/P7/2P3PP/2K4R b - - 0 25
8/1B1bk2p/3p2p1/6r1/3R4/P1r5/2P3PP/2K4R w - - 1 26
8/1B1bk2p/3p2p1/6r1/3R4/P1r5/2P3PP/2K1R3 b - - 2 26
8/1B1bk2p/3p2p1/4r3/3R4/P1r5/2P3PP/2K1R3 w - - 3 27
8/1B1bk2p/3p2p1/4r3/4R3/P1r5/2P3PP/2K1R3 b - - 4 27
8/1Brbk2p/3p2p1/4r3/4R3/P7/2P3PP/2K1R3 w - - 5 28
8/2rbk2p/B2p2p1/4r3/4R3/P7/2P3PP/2K1R3 b - - 6 28
8/2rb3p/B2p1kp1/4r3/4R3/P7/2P3PP/2K1R3 w - - 7 29
8/2rb3p/B2p1kp1/4r3/4R3/P7/2P3PP/2K2R2 b - - 8 29
8/2r4p/B2p1kp1/4rb2/4R3/P7/2P3PP/2K2R2 w - - 9 30
8/2r4p/3p1kp1/4rb2/4R3/P2B4/2P3PP/2K2R2 b - - 10 30
8/2r3kp/3p2p1/4rb2/4R3/P2B4/2P3PP/2K2R2 w - - 11 31
8/2r3kp/3p2p1/4rb2/3R4/P2B4/2P3PP/2K2R2 b - - 12 31
8/2r3kp/3p2p1/4r3/3R4/P2b4/2P3PP/2K2R2 w - - 0 32
8/2r3kp/3p2p1/4r3/8/P2R4/2P3PP/2K2R2 b - - 0 32
8/2r3kp/3p2p1/8/8/P2R4/2P1r1PP/2K2R2 w - - 1 33
8/2r3kp/3p2p1/8/8/P7/2PRr1PP/2K2R2 b - - 2 33
8/2r3kp/3p2p1/8/8/P3r3/2PR2PP/2K2R2 w - - 3 34
8/2r3kp/3R2p1/8/8/P3r3/2P3PP/2K2R2 b - - 0 34
8/2r3kp/3R2p1/8/8/r7/2P3PP/2K2R2 w - - 0 35
8/2r3kp/3R2p1/8/8/r7/2P3PP/2KR4 b - - 1 35
8/r1r3kp/3R2p1/8/8/8/2P3PP/2KR4 w - - 2 36
8/r1r3kp/6p1/8/8/3R4/2P3PP/2KR4 b - - 3 36
8/1rr3kp/6p1/8/8/3R4/2P3PP/2KR4 w - - 4 37
8/1rr3kp/6p1/8/8/2PR4/6PP/2KR4 b - - 0 37
8/r1r3kp/6p1/8/8/2PR4/6PP/2KR4 w - - 1 38
8/r1r3kp/6p1/8/8/2PR4/3R2PP/2K5 b - - 2 38
8/r1r4p/6pk/8/8/2PR4/3R2PP/2K5 w - - 3 39
8/r1r4p/6pk/8/8/2PR4/2KR2PP/8 b - - 4 39
8/2r4p/6pk/8/8/r1PR4/2KR2PP/8 w - - 5 40
8/2r4p/6pk/8/8/r1PR4/1K1R2PP/8 b - - 6 40
8/2r4p/6pk/r7/8/2PR4/1K1R2PP/8 w - - 7 41
8/2r4p/6pk/r7/3R4/2P5/1K1R2PP/8 b - - 8 41
8/1r5p/6pk/r7/3R4/2P5/1K1R2PP/8 w - - 9 42
8/1r5p/6pk/r7/1R6/2P5/1K1R2PP/8 b - - 10 42
8/r6p/6pk/r7/1R6/2P5/1K1R2PP/8 w - - 11 43
8/r6p/6pk/r7/7R/2P5/1K1R2PP/8 b - - 12 43
8/r5kp/6p1/r7/7R/2P5/1K1R2PP/8 w - - 13 44
8/r5kp/6p1/r7/3R4/2P5/1K1R2PP/8 b - - 14 44
8/r5kp/6p1/8/3R4/2P5/rK1R2PP/8 w - - 15 45
8/r5kp/6p1/8/3R4/1KP5/r2R2PP/8 b - - 16 45
8/r5kp/6p1/8/3R4/rKP5/3R2PP/8 w - - 17 46
8/r5kp/6p1/8/3R4/r1P5/2KR2PP/8 b - - 18 46
8/r5kp/6p1/8/3R4/2P5/r1KR2PP/8 w - - 19 47
8/r5kp/6p1/8/3R4/2P5/r2R2PP/3K4 b - - 20 47
8/r5kp/6p1/8/3R4/2P5/3R2PP/r2K4 w - - 21 48
8/r5kp/6p1/8/3R4/2P5/3RK1PP/r7 b - - 22 48
8/r6p/6pk/8/3R4/2P5/3RK1PP/r7 w - - 23 49
8/r2R3p/6pk/8/8/2P5/3RK1PP/r7 b - - 24 49
8/3R3p/6pk/r7/8/2P5/3RK1PP/r7 w - - 25 50
8/7p/6pk/r2R4/8/2P5/3RK1PP/r7 b - - 26 50
8/7p/r5pk/3R4/8/2P5/3RK1PP/r7 w - - 27 51
8/7p/r5pk/3R4/8/2PK4/3R2PP/r7 b - - 28 51
8/7p/r5pk/3R4/8/2PK4/3R2PP/2r5 w - - 29 52
8/7p/r5pk/8/3R4/2PK4/3R2PP/2r5 b - - 30 52
8/7p/6pk/8/3R4/r1PK4/3R2PP/2r5 w - - 31 53
8/7p/6pk/8/2R5/r1PK4/3R2PP/2r5 b - - 32 53
8/r6p/6pk/8/2R5/2PK4/3R2PP/2r5 w - - 33 54
8/r6p/6pk/8/2R5/2PK4/4R1PP/2r5 b - - 34 54
8/r6p/6pk/8/2R5/2PK4/4R1PP/3r4 w - - 35 55
8/r6p/6pk/8/2R5/2P5/2K1R1PP/3r4 b - - 36 55
8/r6p/6pk/8/2R5/2P5/2K1R1PP/r7 w - - 37 56
8/r6p/6pk/8/7R/2P5/2K1R1PP/r7 b - - 38 56
8/r5kp/6p1/8/7R/2P5/2K1R1PP/r7 w - - 39 57
8/r5kp/6p1/8/4R3/2P5/2K1R1PP/r7 b - - 40 57
8/r6p/6pk/8/4R3/2P5/2K1R1PP/r7 w - - 41 58
8/r6p/6pk/8/4R3/2P4P/2K1R1P1/r7 b - - 0 58
8/r6p/6pk/8/4R3/2P4P/r1K1R1P1/8 w - - 1 59
8/r6p/6pk/8/4R3/2PK3P/r3R1P1/8 b - - 2 59
8/r6p/6pk/8/4R3/r1PK3P/4R1P1/8 w - - 3 60
8/r3R2p/6pk/8/8/r1PK3P/4R1P1/8 b - - 4 60
8/4R2p/r5pk/8/8/r1PK3P/4R1P1/8 w - - 5 61
8/4R2p/r5pk/8/4R3/r1PK3P/6P1/8 b - - 6 61
8/4R2p/2r3pk/8/4R3/r1PK3P/6P1/8 w - - 7 62
8/4R2p/2r3pk/8/2R5/r1PK3P/6P1/8 b - - 8 62
8/4R2p/3r2pk/8/2R5/r1PK3P/6P1/8 w - - 9 63
8/4R2p/3r2pk/8/3R4/r1PK3P/6P1/8 b - - 10 63
8/4R2p/2r3pk/8/3R4/r1PK3P/6P1/8 w - - 11 64
8/4R2p/2r3pk/8/7R/r1PK3P/6P1/8 b - - 12 64
8/4R2p/2r3p1/6k1/7R/r1PK3P/6P1/8 w - - 13 65
8/4R2p/2r3p1/6k1/2R5/r1PK3P/6P1/8 b - - 14 65
8/4R2p/3r2p1/6k1/2R5/r1PK3P/6P1/8 w - - 15 66
8/4R2p/3r2p1/6k1/2R5/r1P1K2P/6P1/8 b - - 16 66
8/4R3/3r2p1/6kp/2R5/r1P1K2P/6P1/8 w - h6 0 67
8/4R3/3r2p1/6kp/2R4P/r1P1K3/6P1/8 b - - 0 67
8/4R3/3r1kp1/7p/2R4P/r1P1K3/6P1/8 w - - 1 68
8/2R5/3r1kp1/7p/2R4P/r1P1K3/6P1/8 b - - 2 68
8/2R5/4rkp1/7p/2R4P/r1P1K3/6P1/8 w - - 3 69
8/2R5/4rkp1/7p/4R2P/r1P1K3/6P1/8 b - - 4 69
8/2R5/3r1kp1/7p/4R2P/r1P1K3/6P1/8 w - - 5 70
8/8/3r1kp1/2R4p/4R2P/r1P1K3/6P1/8 b - - 6 70
8/8/3r1kp1/2R4p/4R2P/2P1K3/r5P1/8 w - - 7 71
8/8/3r1kp1/2R4p/5R1P/2P1K3/r5P1/8 b - - 8 71
8/8/3rk1p1/2R4p/5R1P/2P1K3/r5P1/8 w - - 9 72
8/8/3rk1p1/6Rp/5R1P/2P1K3/r5P1/8 b - - 10 72
8/3k4/3r2p1/6Rp/5R1P/2P1K3/r5P1/8 w - - 11 73
8/3k4/3r2p1/6Rp/3R3P/2P1K3/r5P1/8 b - - 12 73
8/3k4/r2r2p1/6Rp/3R3P/2P1K3/6P1/8 w - - 13 74
8/3k4/r2r2R1/7p/3R3P/2P1K3/6P1/8 b - - 0 74
8/3k4/r5R1/7p/3r3P/2P1K3/6P1/8 w - - 0 75
8/3k2R1/r7/7p/3r3P/2P1K3/6P1/8 b - - 1 75
8/6R1/r3k3/7p/3r3P/2P1K3/6P1/8 w - - 2 76
8/6R1/r3k3/7p/3P3P/4K3/6P1/8 b - - 0 76
```

## license

MIT License, see LICENSE file.
