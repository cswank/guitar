package metronome

import (
	"bytes"
	"io"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/hajimehoshi/oto"
)

var (
	sig struct{}
)

type Metronome struct {
	quit   chan struct{}
	player *oto.Player
	low    []byte
	high   []byte
}

func New(box *rice.Box) (*Metronome, error) {
	var err error
	low, err := box.Bytes("low")
	if err != nil {
		return nil, err
	}

	high, err := box.Bytes("high")
	if err != nil {
		return nil, err
	}

	ctx, err := oto.NewContext(44100, 2, 2, 128)
	if err != nil {
		return nil, err
	}

	return &Metronome{
		player: ctx.NewPlayer(),
		quit:   make(chan struct{}),
		low:    low,
		high:   high,
	}, nil
}

func (m Metronome) Start(bpm time.Duration, cb func(ts time.Time)) {
	d := time.Minute / bpm
	tk := time.NewTicker(d)
	l := bytes.NewReader(m.low)
	h := bytes.NewReader(m.high)
	var i int
	go func() {
		for {
			select {
			case <-tk.C:
				switch {
				case i < 4:
					io.Copy(m.player, h)
					h.Seek(0, 0)
				case i%4 == 0:
					cb(time.Now())
					io.Copy(m.player, h)
					h.Seek(0, 0)
				default:
					cb(time.Now())
					io.Copy(m.player, l)
					l.Seek(0, 0)
				}
				i++
			case <-m.quit:
				return
			}
		}
	}()
}

func (m Metronome) Stop() {
	m.quit <- sig
	m.player.Close()
}
