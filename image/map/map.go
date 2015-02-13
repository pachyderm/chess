package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jdoliner/uci"
	"github.com/wfreeman/pgn"
)

var msPerMove int = 2000

func handler(w http.ResponseWriter, r *http.Request) {
}

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
}

func main() {
	//log.Print("Listening on port 80...")
	//http.HandleFunc("/", handler)
	//log.Fatal(http.ListenAndServe(":80", nil))
	f, err := os.Open("../../data/sample/xx0005")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	ps := pgn.NewPGNScanner(f)

	eng, err := uci.NewEngine("/usr/games/stockfish")
	if err != nil {
		log.Fatal(err)
	}

	// set some engine options
	eng.SetOptions(uci.Options{
		Hash:    128,
		Ponder:  false,
		OwnBook: true,
		MultiPV: 4,
	})

	encoder := json.NewEncoder(os.Stdout)

	// while there's more to read in the file
	for ps.Next() {
		// scan the next game
		game, err := ps.Scan()
		if err != nil {
			log.Fatal(err)
		}

		// print out tags
		fmt.Println(game.Tags)

		var score Score
		// make a new board so we can get FEN positions
		b := pgn.NewBoard()
		for i, move := range append(game.Moves, pgn.Move{}) {
			// set the position on the board
			eng.SetFEN(b.String())
			// print out FEN for each move in the game
			resultOpts := uci.HighestDepthOnly
			results, err := eng.Go("movetime", msPerMove, resultOpts)
			if err != nil {
				log.Fatal(err)
			}
			scoreResult := results.Results[0]

			if i > 0 {
				score.PlayedMoveScore = (-1 * scoreResult.Score)
				encoder.Encode(score)
				log.Print("\n")
			}
			if i < len(game.Moves) {
				// make the move on the board
				b.MakeMove(move)

				// Set the score for the last move

				score =
					Score{
						Position:      b.String(),
						PlayedMove:    move.String(),
						BestMoves:     scoreResult.BestMoves,
						BestMoveScore: scoreResult.Score,
						White:         game.Tags["White"],
						Black:         game.Tags["Black"],
						WhiteElo:      game.Tags["WhiteElo"],
						BlackElo:      game.Tags["BlackElo"],
						HalfMoves:     i,
					}
			}
		}
	}
}
