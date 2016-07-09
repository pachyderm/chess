package main

import (
	"encoding/json"
	"log"
	"math"
	"os"

	"github.com/freeeve/pgn"
	"github.com/freeeve/uci"
)

var msPerMove int = 5000
var depth int = 15

type Score struct {
	Position        string   `json:"position"`
	PlayedMove      string   `json:"played-move"`
	BestMoves       []string `json:"best-moves"`
	PlayedMoveScore int      `json:"played-move-score"`
	BestMoveScore   int      `json:"best-move-score"`
	White           string   `json:"white"`
	Black           string   `json:"black"`
	WhiteElo        string   `json:"white-elo"`
	BlackElo        string   `json:"black-elo"`
	HalfMoves       int      `json:"half-moves"`
	TotalHalfMoves  int      `json:"total-half-moves"`
	Event           string   `json:"event"`
	Result          string   `json:"result"`
	Mover           string   `json:"mover"`
	Date            string   `json:"date"`
}

func getScore(s uci.ScoreResult) int {
	if s.Mate && s.Score > 0 {
		return math.MaxInt32
	} else if s.Mate && s.Score < 0 {
		return math.MinInt32
	} else {
		return s.Score
	}
}

func main() {
	ps := pgn.NewPGNScanner(os.Stdin)
	eng, err := uci.NewEngine("/usr/games/stockfish")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	// set some engine options
	err = eng.SetOptions(uci.Options{
		Hash:    128,
		Ponder:  false,
		OwnBook: true,
		MultiPV: 1,
		Threads: 4,
	})
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	encoder := json.NewEncoder(os.Stdout)

	// while there's more to read in the file
	for ps.Next() {
		// scan the next game
		game, err := ps.Scan()
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		if game == nil {
			continue
		}

		var score Score
		// make a new board so we can get FEN positions
		b := pgn.NewBoard()

		// we want one dummy move after the last one to make our loop easier
		for i, move := range append(game.Moves, pgn.Move{}) {
			// set the position on the board
			err := eng.SetFEN(b.String())
			if err != nil {
				log.Print(err)
				os.Exit(1)
			}
			// print out FEN for each move in the game
			results, err := eng.GoDepth(depth, uci.HighestDepthOnly)
			if err != nil {
				log.Print(err)
				os.Exit(1)
			}
			scoreResult := results.Results[0]

			if i > 0 {
				score.PlayedMoveScore = (-1 * getScore(scoreResult))
				encoder.Encode(score)
			}
			if i < len(game.Moves) {

				// Set the score for the last move

				score =
					Score{
						Position:       b.String(),
						PlayedMove:     move.String(),
						BestMoves:      scoreResult.BestMoves,
						BestMoveScore:  getScore(scoreResult),
						White:          game.Tags["White"],
						Black:          game.Tags["Black"],
						WhiteElo:       game.Tags["WhiteElo"],
						BlackElo:       game.Tags["BlackElo"],
						HalfMoves:      i,
						TotalHalfMoves: len(game.Moves),
						Event:          game.Tags["Event"],
						Result:         game.Tags["Result"],
						Date:           game.Tags["Date"],
					}
				if i%2 == 0 {
					score.Mover = game.Tags["White"]
				} else {
					score.Mover = game.Tags["Black"]
				}

				// make the move on the board
				err := b.MakeMove(move)
				if err != nil {
					log.Print(err)
					os.Exit(1)
				}
			}
		}
	}
}