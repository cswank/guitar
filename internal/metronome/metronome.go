package metronome

import (
	"bytes"
	"io"
	"log"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/hajimehoshi/oto"
)

var (
	quit   chan struct{}
	player *oto.Player
	sig    struct{}
	low    []byte
	high   []byte
)

func Init(box *rice.Box) error {
	var err error
	low, err = box.Bytes("low")
	if err != nil {
		return err
	}

	high, err = box.Bytes("low")
	return err
}

func init() {
	quit = make(chan struct{})
	ctx, err := oto.NewContext(44100, 2, 2, 2048)
	if err != nil {
		log.Fatal(err)
	}
	player = ctx.NewPlayer()
}

func Start(bpm time.Duration, cb func(ts time.Time)) {
	d := time.Minute / bpm
	tk := time.NewTicker(d)
	l := bytes.NewReader(low)
	h := bytes.NewReader(high)
	var i int
	go func() {
		for {
			select {
			case <-tk.C:
				if i%4 == 0 {
					cb(time.Now())
					io.Copy(player, h)
					h.Seek(0, 0)
				} else {
					cb(time.Now())
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

func Stop() {
	quit <- sig
	player.Close()
}
