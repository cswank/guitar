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
		notes  func() []time.Duration
		tab    string
		scores []time.Duration
	}{
		{
			bpm: 60,
			notes: func() []time.Duration {
				out := make([]time.Duration, 17)
				for i := range out {
					out[i] = time.Duration(i) * time.Second
				}
				return out
			},
			loops: 2,
			tab: `
|--------------------------------|--------3-------3---------------|
|--------------------------------|6-----------------------6-------|
|------------------------5-------|--------------------------------|
|--------3-------5---------------|--------------------------------|
|3-------------------------------|--------------------------------|
|--------------------------------|--------------------------------|`,
			scores: []time.Duration{0, 0},
		},
		{
			bpm: 60,
			notes: func() []time.Duration {
				out := make([]time.Duration, 17)
				for i := range out {
					d := time.Duration(i) * time.Second
					if i > 7 {
						d += time.Millisecond
					}
					out[i] = d
				}
				return out
			},
			loops: 2,
			tab: `
|--------------------------------|--------3-------3---------------|
|--------------------------------|6-----------------------6-------|
|------------------------5-------|--------------------------------|
|--------3-------5---------------|--------------------------------|
|3-------------------------------|--------------------------------|
|--------------------------------|--------------------------------|`,
			scores: []time.Duration{0, 8},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			start := time.Now()
			in := tc.notes()
			notes := make([]time.Time, len(in))
			for i, dur := range in {
				notes[i] = start.Add(dur)
			}
			tab, err := music.New(bytes.NewBufferString(tc.tab), "tab", tc.bpm)
			if !assert.NoError(t, err) {
				return
			}

			met := &metronome{
				loops: tc.loops,
				music: *tab,
				start: notes[0],
				bpm:   time.Duration(tc.bpm),
			}

			a := analysis.New(met, &input{notes: notes})
			var scores []time.Duration
			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				a.Start(60, tab, func(score time.Duration) {
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
	music music.Input
	start time.Time
	bpm   time.Duration
	lock  *sync.Mutex
}

func (m *metronome) Start(ts time.Duration, f func(time.Time)) {
	for i := 0; i <= (m.loops * m.music.Time.Beats * m.music.Measures); i++ {
		f(m.start.Add(time.Duration(i) * time.Minute / m.bpm))
	}
}

func (m *metronome) Stop() {}

type input struct {
	notes []time.Time
	lock  *sync.Mutex
}

func (i *input) Start(f func(time.Time)) error {
	for _, n := range i.notes {
		f(n)
	}
	return nil
}

func (m *input) Stop() {}
