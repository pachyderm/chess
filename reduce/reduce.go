package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Game struct {
	Date     string `json:"date"`
	Elo      int    `json:"elo"`
	Blunders []int  `json:"blunders"`
}

var blundersByRating map[string][10]int = make(map[string][10]int)

func bucket(rating int) string {
	return strconv.Itoa((rating / 25) * 25)
}

func handler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)

	for {
		var g Game
		if err := dec.Decode(&g); err == io.EOF {
			break
		} else if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		blunders := blundersByRating[bucket(g.Elo)]
		for i, value := range g.Blunders {
			blunders[i] += value
		}
		blundersByRating[bucket(g.Elo)] = blunders
	}
	err := json.NewEncoder(w).Encode(blundersByRating)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func main() {
	log.Print("Listening on port 80...")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
