package music

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

func newTab(r io.Reader, bpm int) (*Input, error) {
	var rows []string

	reader := bufio.NewReader(r)
	var i, l int
	for {
		row, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) == 0 {
			continue
		}

		if i == 0 {
			l = len(row)
		} else {
			if len(row) != l {
				return nil, fmt.Errorf("invalid input, all strings of the tab must be the same length")
			}
		}

		i++
		rows = append(rows, string(row))
	}

	if len(rows) != 6 {
		return nil, fmt.Errorf("invalid input, guitar tab must have 6 strings")
	}

	parts := strings.Split(rows[0], "|")
	if len(parts) <= 2 {
		return nil, fmt.Errorf("invalid input")
	}
	parts = parts[1 : len(parts)-1]
	beat := 4
	if len(parts[1])%3 == 0 {
		beat = 3
	}

	timePerBeat := time.Minute / time.Duration(bpm)
	timePerMeasure := time.Duration(beat) * timePerBeat
	timePerDivision := timePerMeasure / time.Duration(len(parts[1]))

	var notes []note
	var j int
	for i := range rows[0] {
		if rows[0][i] == '|' {
			continue
		}
		for _, row := range rows {
			if row[i] != '-' {
				notes = append(notes, note{t: timePerDivision * time.Duration(j)})
				break //TODO: support polyphonic, only one note at a time for now (and probably forever?)
			}
		}
		j++
	}

	return &Input{
		Time:        Time{Beat: beat, Beats: 4},
		Measures:    len(parts),
		timePerLoop: timePerMeasure * time.Duration(len(parts)),
		notes:       notes,
	}, nil
}
