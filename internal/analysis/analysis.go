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

func (a *Analysis) Start(bpm time.Duration, in *music.Input, f func(int)) {
	notes := []time.Time{}

	newNote := make(chan time.Time)
	newTick := make(chan time.Time)
	analyze := make(chan struct{})

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
	go func() {
		for {
			select {
			case <-newTick:
				i++
				fmt.Println("tick", i, end)
				if i > 1 && i%end == 1 {
					analyze <- sig
				}
			case <-a.quit2:
				a.wg.Done()
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-analyze:
				f(100)
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
