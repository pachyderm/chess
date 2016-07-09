# uci 

[![Build Status](https://travis-ci.org/freeeve/uci.png?branch=master)](https://travis-ci.org/freeeve/uci)
[![Coverage Status](https://coveralls.io/repos/freeeve/uci/badge.png)](https://coveralls.io/r/freeeve/uci)

A golang API to interact with UCI chess engines. (should be considered experimental for the time being) [A description of how UCI works is available here.](http://wbec-ridderkerk.nl/html/UCIProtocol.html)

Many chess engines support UCI (Universal Chess Interface). This library is designed for use with Stockfish, but should work with other UCI engines.

[Godoc for UCI](http://godoc.org/github.com/freeeve/uci)

## minimum viable snippet
```Go
package main

import (
	"fmt"
	"log"
	"gopkg.in/freeeve/uci.v1"
)

func main() {
	eng, err := uci.NewEngine("/path/to/stockfish")
	if err != nil {
		log.Fatal(err)
	}
	
	// set some engine options
	eng.SetOptions(uci.Options{
		Hash:128,
		Ponder:false,
		OwnBook:true,
		MultiPV:4,
	})

	// set the starting position
	eng.SetFEN("rnb4r/ppp1k1pp/3bp3/1N3p2/1P2n3/P3BN2/2P1PPPP/R3KB1R b KQ - 4 11")
	
	// set some result filter options
	resultOpts := uci.HighestDepthOnly | uci.IncludeUpperbounds | uci.IncludeLowerbounds
	results := eng.GoDepth(10, resultOpts)

	// print it (String() goes to pretty JSON for now)
	fmt.Println(results)
}
```

produces this output:

```
{
  "BestMove": "c8d7",
  "Results": [
    {
      "Time": 136,
      "Depth": 10,
      "SelDepth": 16,
      "Nodes": 183853,
      "NodesPerSecond": 1351860,
      "MultiPV": 1,
      "Lowerbound": false,
      "Upperbound": false,
      "Score": 20,
      "Mate": false,
      "BestMoves": [
        "c8d7",
        "b5d6",
        "c7d6",
        "f3g5",
        "b8c6",
        "g5e4",
        "f5e4",
        "a1d1",
        "h7h6",
        "f2f4",
        "a8f8"
      ]
    },
    {
      "Time": 136,
      "Depth": 10,
      "SelDepth": 16,
      "Nodes": 183853,
      "NodesPerSecond": 1351860,
      "MultiPV": 2,
      "Lowerbound": false,
      "Upperbound": false,
      "Score": 2,
      "Mate": false,
      "BestMoves": [
        "a7a5",
        "b5d6",
        "c7d6",
        "b4b5",
        "b8d7",
        "f3d2",
        "e4d2",
        "e3d2",
        "b7b6",
        "a3a4"
      ]
    },
    {
      "Time": 136,
      "Depth": 10,
      "SelDepth": 16,
      "Nodes": 183853,
      "NodesPerSecond": 1351860,
      "MultiPV": 3,
      "Lowerbound": false,
      "Upperbound": false,
      "Score": -4,
      "Mate": false,
      "BestMoves": [
        "b8c6",
        "c2c4",
        "a7a6",
        "b5d6",
        "c7d6",
        "f3g5",
        "c8d7",
        "c4c5",
        "d6d5",
        "g5e4",
        "f5e4"
      ]
    },
    {
      "Time": 136,
      "Depth": 10,
      "SelDepth": 16,
      "Nodes": 183853,
      "NodesPerSecond": 1351860,
      "MultiPV": 4,
      "Lowerbound": false,
      "Upperbound": false,
      "Score": -18,
      "Mate": false,
      "BestMoves": [
        "e7f7",
        "h2h4",
        "a7a5",
        "b5d6",
        "c7d6",
        "b4a5",
        "b8d7",
        "f3g5",
        "e4g5",
        "h4g5",
        "a8a5",
        "e3d2"
      ]
    }
  ]
}
```
