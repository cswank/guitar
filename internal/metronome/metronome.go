package metronome

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/hajimehoshi/oto"
)

var (
	ch     chan struct{}
	quit   chan struct{}
	player *oto.Player
	signal struct{}
	Low    []byte
	High   []byte
)

func init() {
	ch = make(chan struct{})
	quit = make(chan struct{})
	ctx, err := oto.NewContext(44100, 2, 2, 2048)
	if err != nil {
		log.Fatal(err)
	}
	player = ctx.NewPlayer()
}

func Start(bpm int) {
	d := time.Minute / time.Duration(bpm)
	tk := time.NewTicker(d)
	l := bytes.NewReader(Low)
	h := bytes.NewReader(High)
	var i int
	go func() {
		for {
			select {
			case <-tk.C:
				fmt.Println("tick")
				if i%4 == 0 {
					io.Copy(player, h)
					h.Seek(0, 0)
				} else {
					io.Copy(player, l)
					l.Seek(0, 0)
				}
				i++
			case <-quit:
				return
			}
		}
	}()
}

func Tick() {
	ch <- signal
}

func Stop() {
	quit <- signal
	player.Close()
}
