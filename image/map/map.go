package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"

	"github.com/jdoliner/uci"
	"github.com/wfreeman/pgn"
)

var msPerMove int = 2000

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
	Mover           string   `json:"result"`
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

func handler(w http.ResponseWriter, r *http.Request) {
	ps := pgn.NewPGNScanner(r.Body)

	eng, err := uci.NewEngine("/usr/games/stockfish")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	// set some engine options
	eng.SetOptions(uci.Options{
		Hash:    128,
		Ponder:  false,
		OwnBook: true,
		MultiPV: 4,
	})

	encoder := json.NewEncoder(w)

	// while there's more to read in the file
	for ps.Next() {
		// scan the next game
		game, err := ps.Scan()
		if err != nil {
			http.Error(w, err.Error(), 500)
		}

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
				http.Error(w, err.Error(), 500)
			}
			scoreResult := results.Results[0]

			if i > 0 {
				score.PlayedMoveScore = (-1 * getScore(scoreResult))
				if score.BestMoves[0] == score.PlayedMove {
					score.BestMoveScore = (-1 * getScore(scoreResult))
				}
				encoder.Encode(score)
			}
			if i < len(game.Moves) {
				// make the move on the board
				b.MakeMove(move)

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
					}
				if i%2 == 0 {
					score.Mover = "White"
				} else {
					score.Mover = "Black"
				}
			}
		}
	}
}

func main() {
	log.Print("Listening on port 8080...")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
