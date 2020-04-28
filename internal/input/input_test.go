package input_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/cswank/guitar/internal/input"
	"github.com/stretchr/testify/assert"
)

func TestInput(t *testing.T) {
	testCases := []struct {
		input  string
		typ    string
		output input.Input
		bpm    int
		score  func(t *testing.T, i *input.Input)
	}{
		{
			input: `
|--------------------------------|--------3---6---3---------------|--------------------------------|
|----------------------------4---|6---7---------------7---6---4---|--------------------------------|
|--------------------3---5-------|--------------------------------|5---3---------------------------|
|--------3---4---5---------------|--------------------------------|--------5---4---3---------------|
|3---6---------------------------|--------------------------------|--------------------6---3-------|
|--------------------------------|--------------------------------|----------------------------6---|`,
			typ: "tab",
			bpm: 60,
			output: input.Input{
				Time: input.Time{
					Beat:  4,
					Beats: 4,
				},
				Measures: 3,
			},
			score: func(t *testing.T, in *input.Input) {
				ts := time.Now()
				for i := time.Duration(0); i < 24; i++ {
					assert.Equal(t, in.Score(ts, ts.Add(i*500*time.Millisecond), 100), input.Score{Diff: 0})
				}
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			out, err := input.New(bytes.NewBuffer([]byte(tc.input)), tc.typ, tc.bpm)
			assert.NoError(t, err)
			assert.Equal(t, tc.output.Time, out.Time)
			tc.score(t, out)
		})
	}
}
