package input

import (
	"fmt"
	"io"
	"time"
)

type (
	Scorer func(start, ts time.Time, freq, bpm int) Score

	Score struct {
		Diff time.Duration
	}

	Time struct {
		Beat  int
		Beats int
	}

	Input struct {
		Time        Time
		Measures    int
		timePerLoop time.Duration
		currentNote int
		notes       []note
	}

	note struct {
		// expected time the note should be hit
		t    time.Duration
		freq int
	}
)

func (i *Input) Score(start, ts time.Time, freq int) Score {
	if i.currentNote >= len(i.notes) {
		return Score{Diff: time.Minute}
	}
	n := i.notes[i.currentNote]
	diff := ts.Sub(start) % i.timePerLoop
	i.currentNote++
	return Score{Diff: diff - n.t}
}

func New(r io.Reader, typ string, bpm int) (*Input, error) {
	switch typ {
	case "tab":
		return newTab(r, bpm)
	default:
		return nil, fmt.Errorf("unknown input type: %s", typ)
	}
}
