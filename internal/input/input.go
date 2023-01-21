package input

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/mattetti/audio/dsp/filters"
	"github.com/mattetti/audio/dsp/windows"
)

var (
	sig struct{}
)

const (
	bufSize = 4096
	//bufSize = 1024
)

type Input struct {
	quit   chan struct{}
	wg     sync.WaitGroup
	filter *filters.FIR
	stream *portaudio.Stream
}

func New() *Input {
	return &Input{
		quit: make(chan struct{}),
		filter: &filters.FIR{
			Sinc: &filters.Sinc{
				CutOffFreq:   5000,
				SamplingFreq: 44100,
				Taps:         10,
				Window:       windows.Blackman,
			},
		},
	}
}

func (i *Input) Start(cb func(ts time.Time)) error {
	i.wg.Add(1)

	portaudio.Initialize()

	in := make([]float32, bufSize)
	var err error
	i.stream, err = portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	if err != nil {
		return err
	}

	if err := i.stream.Start(); err != nil {
		return err
	}

	var stop bool
	var prev float32
	//var prev int
	prevTime := time.Now()
	prevLoop := time.Now()

	go func() {
		for !stop {
			select {
			case <-i.quit:
				stop = true
			default:
				i.stream.Read()
				now := time.Now()
				diff := now.Sub(prevTime)
				n, npeaks, max := peak(in)
				fmt.Printf("max: %f, peaks: %d, diff: %s, time: %s\n", max, npeaks, diff, now)
				if (max-prev) > 0.05 && diff > time.Millisecond*100 {
					fmt.Println("peak!!!!!!")
					diff = now.Sub(prevLoop)
					cb(prevLoop.Add((diff * time.Duration(n)) / bufSize))
					//cb(now)
					prevTime = now
				}
				prev = max
				prevLoop = now
			}
		}
		i.wg.Done()
	}()

	return nil
}

func (i *Input) Stop() {
	i.quit <- sig
	i.wg.Wait()
	portaudio.Terminate()
	i.stream.Close()
}

type stat struct {
	n   int
	max float32
	sum float32
}

/*
However, by setting a minimum change in amplitude δ for maxima and
minima detection the noise in the signal is ignored [6]. The value for
δ is determined by making a pre- liminary sweep with δ = 0.01. For the
true peak detection pass, δ is defined as the average of the squared
magnitudes of the peaks from the first pass divided by an empirically
determined factor = 3.2.
*/

func peak(in []float32) (int, int, float32) {
	in64 := make([]float64, len(in))
	for i, f := range in {
		in64[i] = float64(f)
	}

	_, minmaxvals := peaks(in64, 0.01)

	if len(minmaxvals) == 0 {
		return 0, 0, 0.0
	}

	var sum float64
	var n float64
	for _, val := range minmaxvals {
		sum += val * val
		n += 1
	}

	sig := (sum / n) * 3.2

	minmaxx, minmaxvals2 := peaks(in64, sig)
	if len(minmaxx) == 0 {
		return 0, 0, 0.0
	}

	var max float64
	var idx int
	for i, f := range minmaxvals2 {
		if math.Abs(f) > max {
			max = math.Abs(f)
			idx = i
		}
	}

	i := minmaxx[idx]
	return i, len(minmaxx), float32(max)
}

func peaks(vals []float64, delta float64) ([]int, []float64) {
	minmaxx := make([]int, 0, 0)
	out := make([]float64, 0, 0)

	if delta < 0 {
		return minmaxx, vals
	}

	mn := math.Inf(1)
	mx := math.Inf(-1)
	mnpos := int(math.NaN())
	mxpos := int(math.NaN())

	lookForMax := true

	for i, val := range vals {
		if val > mx {
			mx = val
			mxpos = i
		}
		if val < mn {
			mn = val
			mnpos = i
		}

		if lookForMax {
			if val < mx-delta {
				minmaxx = append(minmaxx, mxpos)
				out = append(out, mx)

				mn = val
				mnpos = i
				lookForMax = false
			}
		} else {
			if val > mn+delta {
				minmaxx = append(minmaxx, mnpos)
				out = append(out, mn)

				mx = val
				mxpos = i
				lookForMax = true
			}
		}
	}

	return minmaxx, out
}
