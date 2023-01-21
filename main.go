package main

import (
	"log"
	"os"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/cswank/guitar/internal/analysis"
	"github.com/cswank/guitar/internal/input"
	"github.com/cswank/guitar/internal/metronome"
	"github.com/cswank/guitar/internal/music"
)

var (
	in  = kingpin.Arg("input", "input file").Required().String()
	typ = kingpin.Flag("type", "input file type").Enum("tab")
	bpm = kingpin.Flag("bpm", "beats per minute").Default("60").Int()
)

func main() {
	kingpin.Parse()

	f, err := os.Open(*in)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	m, err := music.New(f, *typ, *bpm)
	if err != nil {
		log.Fatal(err)
	}

	met, err := metronome.New()
	if err != nil {
		log.Fatal(err)
	}

	in := input.New()
	a := analysis.New(met, in)

	a.Start(time.Duration(*bpm), m, func(d time.Duration) {})

	time.Sleep(10 * time.Second)

	a.Stop()
}
