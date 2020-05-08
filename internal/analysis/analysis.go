package analysis

import (
	"fmt"
	"sync"
	"time"

	"github.com/cswank/guitar/internal/music"
)

var (
	sig struct{}
)

type (
	ticker interface {
		Start(time.Duration, func(time.Time))
		Stop()
	}

	listener interface {
		Start(func(time.Time)) error
		Stop()
	}

	Analysis struct {
		ticks     []time.Time
		notes     []time.Time
		wg        sync.WaitGroup
		quit      chan struct{}
		quit2     chan struct{}
		metronome ticker
		input     listener
	}
)

func New(met ticker, in listener) *Analysis {
	return &Analysis{
		quit:      make(chan struct{}),
		quit2:     make(chan struct{}),
		metronome: met,
		input:     in,
	}
}

func (a *Analysis) Start(bpm time.Duration, in *music.Input) {
	a.ticks = []time.Time{}
	a.notes = []time.Time{}

	newNote := make(chan time.Time)
	newTick := make(chan time.Time)

	a.input.Start(func(ts time.Time) {
		newNote <- ts
	})

	a.metronome.Start(bpm, func(ts time.Time) {
		newTick <- ts
	})

	a.wg.Add(2)
	go func() {
		for {
			select {
			case ts := <-newNote:
				a.notes = append(a.notes, ts)
			case <-a.quit:
				a.wg.Done()
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case ts := <-newTick:
				a.ticks = append(a.ticks, ts)
			case <-a.quit2:
				a.wg.Done()
				return
			}
		}
	}()
}

func (a *Analysis) Stop() {
	a.input.Stop()
	a.metronome.Stop()

	a.quit <- sig
	a.quit2 <- sig
	a.wg.Wait()
	for i, t := range a.ticks {
		var n time.Time
		if i < len(a.notes) {
			n = a.notes[i]
		}
		fmt.Printf("diff: %s, tick: %s, note: %s\n", t.Sub(n), t, n)
	}
}
