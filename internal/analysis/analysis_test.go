package analysis_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/cswank/guitar/internal/analysis"
	"github.com/cswank/guitar/internal/music"
	"github.com/stretchr/testify/assert"
)

func TestAnalysis(t *testing.T) {
	testCases := []struct {
		loops  int
		notes  []time.Time
		input  music.Input
		scores []int
	}{
		{
			notes: []time.Time{},
			loops: 1,
			input: music.Input{
				Time:     music.Time{Beat: 4, Beats: 4},
				Measures: 2,
			},
			scores: []int{100},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			a := analysis.New(&metronome{loops: tc.loops, in: tc.input}, &input{notes: tc.notes})
			var scores []int
			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				a.Start(60, &tc.input, func(score int) {
					scores = append(scores, score)
					if len(scores) == tc.loops {
						wg.Done()
					}
				})
			}()

			wg.Wait()
			go a.Stop()
			assert.Equal(t, tc.scores, scores)
		})
	}
}

type metronome struct {
	loops int
	in    music.Input
}

func (m *metronome) Start(ts time.Duration, f func(time.Time)) {
	for i := 0; i <= (m.loops * m.in.Time.Beats * m.in.Measures); i++ {
		t := time.Now()
		f(t)
	}
}

func (m *metronome) Stop() {}

type input struct {
	notes []time.Time
}

func (i *input) Start(f func(time.Time)) error {
	for _, n := range i.notes {
		f(n)
	}
	return nil
}

func (m *input) Stop() {}
