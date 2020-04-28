package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/alecthomas/kingpin"
	"github.com/gordonklaus/portaudio"
)

var (
	tab = kingpin.Arg("input", "input file").Required().String()
	typ = kingpin.Flag("type", "input file type").Enum("tab")
)

func main() {
	kingpin.Parse()

	listen()
}

func listen() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	portaudio.Initialize()
	defer portaudio.Terminate()
	in := make([]float32, 1024)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	chk(err)
	defer stream.Close()

	chk(stream.Start())
	var stop bool
	var prev float32
	var i int
	for !stop {
		i++
		stream.Read()
		m := max(in)
		if (m - prev) > 0.1 {
			fmt.Println(m)
		}
		prev = m
		select {
		case <-sig:
			fmt.Println("break")
			stop = true
		default:
		}
	}

	chk(stream.Stop())
	fmt.Println(i)
	// out := make([]float64, len(buf))
	// for i, val := range buf {
	// 	out[i] = float64(val)
	// }

	// t := fft.FFTReal(out)
	// fmt.Println(t)
}

func max(in []float32) float32 {
	var out float32
	for _, f := range in {
		if f > out {
			out = f
		}
	}
	return out
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
