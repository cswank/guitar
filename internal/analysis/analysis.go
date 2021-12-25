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
		wg                 sync.WaitGroup
		quit, quit2, quit3 chan struct{}
		metronome          ticker
		input              listener
	}
)

func New(met ticker, in listener) *Analysis {
	return &Analysis{
		quit:      make(chan struct{}),
		quit2:     make(chan struct{}),
		quit3:     make(chan struct{}),
		metronome: met,
		input:     in,
	}
}

func (a *Analysis) Start(bpm time.Duration, in *music.Input, f func(time.Duration)) {
	notes := []time.Time{}

	newNote := make(chan time.Time)
	newTick := make(chan time.Time)
	analyze := make(chan time.Time)

	a.wg.Add(3)
	go func() {
		for {
			select {
			case ts := <-newNote:
				notes = append(notes, ts)
			case <-a.quit:
				a.wg.Done()
				return
			}
		}
	}()

	end := in.Time.Beats * in.Measures
	var i int
	var start, prev time.Time

	go func() {
		for {
			select {
			case tick := <-newTick:
				if i == 0 {
					start = tick
				}
				if i > 0 && i%end == 0 {
					analyze <- start
					start = prev
				}
				i++
				prev = tick
			case <-a.quit2:
				a.wg.Done()
				return
			}
		}
	}()

	go func() {
		var n int
		for {
			select {
			case start := <-analyze:
				end := start.Add(in.TimePerLoop)
				var score time.Duration
				chunk := notes[n:]
				for _, note := range chunk {
					if end.Sub(note) < 0 {
						break
					}
					n++
					s := in.Score(start, note, 0)
					fmt.Println("score", s, len(chunk))
					score += s
				}

				fmt.Println("final score", score, n)
				f(score)
			case <-a.quit3:
				a.wg.Done()
				return
			}
		}
	}()

	a.input.Start(func(ts time.Time) {
		newNote <- ts
	})

	a.metronome.Start(bpm, func(ts time.Time) {
		newTick <- ts
	})
}

func (a *Analysis) Stop() {
	a.input.Stop()
	a.metronome.Stop()

	a.quit <- sig
	a.quit2 <- sig
	a.wg.Wait()
}
