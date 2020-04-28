package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/cswank/beep/mp3"
)

var (
	input  = kingpin.Flag("input", "input mp3").Required().String()
	output = kingpin.Flag("output", "output file").Required().String()
)

const max = 32767

func main() {
	kingpin.Parse()

	inf, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}

	defer inf.Close()

	streamer, format, err := mp3.Decode(inf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", format)

	f, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	b := make([][2]float64, 1024)
	var buf []byte
	var data []byte
	for {
		n, ok := streamer.Stream(b)
		fmt.Println(n, ok)
		buf = make([]byte, n*4)
		if !ok {
			break
		}
		for i := 0; i < n; i++ {
			for c := range b[i] {
				val := b[i][c]
				if val < -1 {
					val = -1
				}
				if val > +1 {
					val = +1
				}
				valInt16 := int16(val * (1<<15 - 1))
				low := byte(valInt16)
				high := byte(valInt16 >> 8)
				buf[i*4+c*2+0] = low
				buf[i*4+c*2+1] = high
				//data = append(data, byte(int16(b[i][0]*0.3*max)))
			}
		}
		data = append(data, buf...)
	}

	fmt.Println(len(data))
	binary.Write(f, binary.LittleEndian, data)
}
