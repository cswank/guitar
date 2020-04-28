package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/cswank/guitar/internal/metronome"
)

func main() {
	f, err := os.Open("low")
	if err != nil {
		log.Fatal(err)
	}

	d, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	f.Close()

	metronome.Low = d

	f, err = os.Open("high")
	if err != nil {
		log.Fatal(err)
	}

	d, err = ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	f.Close()

	metronome.High = d

	metronome.Start(120)
	time.Sleep(10 * time.Second)
	metronome.Stop()
}
