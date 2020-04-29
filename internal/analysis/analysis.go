package analysis

import (
	"fmt"
	"sync"
	"time"

	"github.com/cswank/guitar/internal/input"
	"github.com/cswank/guitar/internal/metronome"
	"github.com/cswank/guitar/internal/music"
)

var (
	ticks []time.Time
	notes []time.Time
	wg    sync.WaitGroup
	quit  chan struct{}
	sig   struct{}
)

func init() {
	quit = make(chan struct{})
}

func Start(bpm time.Duration, in *music.Input) {
	ticks = []time.Time{}
	notes = []time.Time{}

	newNote := make(chan time.Time)
	newTick := make(chan time.Time)

	input.Start(func(ts time.Time) {
		newNote <- ts
	})

	metronome.Start(bpm, func(ts time.Time) {
		newTick <- ts
	})

	wg.Add(1)
	go func() {
		for {
			select {
			case ts := <-newNote:
				notes = append(notes, ts)
			case ts := <-newTick:
				ticks = append(ticks, ts)
			case <-quit:
				wg.Done()
				return
			}
		}
	}()

}

func Stop() {
	input.Stop()
	metronome.Stop()

	quit <- sig
	wg.Wait()
	for i, t := range ticks {
		var n time.Time
		if i < len(notes) {
			n = notes[i]
		}
		fmt.Printf("tick: %s, note: %s\n", t, n)
	}
}
