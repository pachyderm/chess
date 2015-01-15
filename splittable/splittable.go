package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func main() {
	var err error
	s := bufio.NewReader(os.Stdin)
	empties := 0
	for l, err := s.ReadBytes('\n'); err == nil; l, err = s.ReadBytes('\n') {
		if len(l) == 1 && empties == 0 {
			empties++
		} else if len(l) == 1 && empties == 1 {
			os.Stdout.Write([]byte("===\n"))
			empties = 0
		} else {
			if empties == 1 {
				os.Stdout.Write([]byte("\n"))
				empties = 0
			}
			_, err := os.Stdout.Write(l)
			if err != nil {
				log.Print(err)
			}
		}
	}
	if err != io.EOF {
		log.Print(err)
	}
}
