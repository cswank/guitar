package music

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
		notes       []note
	}

	note struct {
		// expected time the note should be hit
		t    time.Duration
		freq int
	}
)

func (i *Input) Score(start, ts time.Time, freq int) time.Duration {
	min := time.Hour
	t := ts.Sub(start)

	var dur time.Duration
	if t > 0 {
		dur = i.timePerLoop % t
	}

	for _, note := range i.notes {
		a := abs(dur, note.t)
		if a < min {
			min = a
		}
	}
	return min
}

func abs(t1, t2 time.Duration) time.Duration {
	x := (t1 - t2)
	if x < 0 {
		return x * -1
	}
	return x
}

func New(r io.Reader, typ string, bpm int) (*Input, error) {
	switch typ {
	case "tab":
		return newTab(r, bpm)
	default:
		return nil, fmt.Errorf("unknown input type: %s", typ)
	}
}
