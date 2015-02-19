package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"

	"github.com/jdoliner/pgn"
	"github.com/jdoliner/uci"
)

var msPerMove int = 10000

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

func handler(w http.ResponseWriter, r *http.Request) {
	ps := pgn.NewPGNScanner(r.Body)

	eng, err := uci.NewEngine("/usr/games/stockfish")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// set some engine options
	err = eng.SetOptions(uci.Options{
		Hash:    128,
		Ponder:  false,
		OwnBook: true,
		MultiPV: 1,
		Threads: 20,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	encoder := json.NewEncoder(w)

	// while there's more to read in the file
	for ps.Next() {
		// scan the next game
		game, err := ps.Scan()
		if err != nil {
			http.Error(w, err.Error(), 500)
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
				http.Error(w, err.Error(), 500)
				return
			}
			// print out FEN for each move in the game
			resultOpts := uci.HighestDepthOnly
			results, err := eng.Go("movetime", msPerMove, resultOpts)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			scoreResult := results.Results[0]
			log.Printf("%+v", scoreResult)

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
					http.Error(w, err.Error(), 500)
					return
				}
			}
		}
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	log.Print("Listening on port 80...")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
