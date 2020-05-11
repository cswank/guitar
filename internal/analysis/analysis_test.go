package analysis_test

import (
	"bytes"
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
		bpm    int
		notes  []time.Duration
		tab    string
		scores []int
	}{
		{
			bpm:   60,
			notes: []time.Duration{0, time.Second, 2 * time.Second, 3 * time.Second, 4 * time.Second, 5 * time.Second, 6 * time.Second, 7 * time.Second, 8 * time.Second},
			loops: 1,
			tab: `
|--------------------------------|--------3-------3---------------|
|--------------------------------|6-----------------------6-------|
|------------------------5-------|--------------------------------|
|--------3-------5---------------|--------------------------------|
|3-------------------------------|--------------------------------|
|--------------------------------|--------------------------------|`,
			scores: []int{100},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			start := time.Now()
			notes := make([]time.Time, len(tc.notes))
			for i, dur := range tc.notes {
				notes[i] = start.Add(dur)
			}
			in, err := music.New(bytes.NewBufferString(tc.tab), "tab", tc.bpm)
			if !assert.NoError(t, err) {
				return
			}

			met := &metronome{
				loops: tc.loops,
				in:    *in,
				start: notes[0],
				bpm:   time.Duration(tc.bpm),
			}
			a := analysis.New(met, &input{notes: notes})
			var scores []int
			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				a.Start(60, in, func(score int) {
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
	start time.Time
	bpm   time.Duration
}

func (m *metronome) Start(ts time.Duration, f func(time.Time)) {
	for i := 0; i <= (m.loops * m.in.Time.Beats * m.in.Measures); i++ {
		f(m.start.Add(time.Duration(i) * time.Minute / m.bpm))
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
