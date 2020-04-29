package main

import (
	"log"
	"os"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/alecthomas/kingpin"
	"github.com/cswank/guitar/internal/analysis"
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

	box := rice.MustFindBox("sounds")

	if err := metronome.Init(box); err != nil {
		log.Fatal(err)
	}

	analysis.Start(time.Duration(*bpm), m)

	time.Sleep(10 * time.Second)

	analysis.Stop()
}
