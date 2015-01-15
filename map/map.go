package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func runCrafty(game string, out io.Writer) error {
	log.Print("runCrafty")
	{
		c := exec.Command("/usr/games/crafty")
		stdin, err := c.StdinPipe()
		if err != nil {
			return err
		}
		stderr, err := c.StderrPipe()
		if err != nil {
			return err
		}
		err = c.Start()
		if err != nil {
			return err
		}
		fmt.Fprintf(stdin, "annotate %s wb 1-999 1 2\r\n", game)
		fmt.Fprint(stdin, "quit\r\n")

		buf := new(bytes.Buffer)
		buf.ReadFrom(stderr)
		if buf.Len() != 0 {
			log.Print("Command had output on stderr.\n Cmd: ", strings.Join(c.Args, " "), "\nstderr: ", buf)
		}

		log.Print("Waiting for crafty...")
		err = c.Wait()
		if err != nil {
			return err
		}
		log.Print("Done")
	}

	{
		c := exec.Command("/bin/crafty-to-json", game+".can")
		stdout, err := c.StdoutPipe()
		if err != nil {
			return err
		}
		stderr, err := c.StderrPipe()
		if err != nil {
			return err
		}
		err = c.Start()
		if err != nil {
			return err
		}

		_, err = io.Copy(out, stdout)
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(stderr)
		if buf.Len() != 0 {
			log.Print("Command had output on stderr.\n Cmd: ", strings.Join(c.Args, " "), "\nstderr: ", buf)
		}

		log.Print("Waiting for crafty...")
		err = c.Wait()
		if err != nil {
			return err
		}
		log.Print("Done")
	}

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Start")
	f, err := ioutil.TempFile("", "game")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer f.Close()
	//defer os.Remove(f.Name())
	_, err = io.Copy(f, r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = runCrafty(f.Name(), w)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func main() {
	log.Print("Listening on port 80...")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
