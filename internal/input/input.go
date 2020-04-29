package input

import (
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
)

var (
	stream *portaudio.Stream
	quit   chan struct{}
	sig    struct{}
	wg     sync.WaitGroup
)

func init() {
	quit = make(chan struct{})
}

func Start(cb func(ts time.Time)) error {
	wg.Add(1)

	portaudio.Initialize()

	in := make([]float32, 1024)
	var err error
	stream, err = portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	if err != nil {
		return err
	}

	if err := stream.Start(); err != nil {
		return err
	}

	var stop bool
	var prev float32
	prevTime := time.Now()
	go func() {
		for !stop {
			select {
			case <-quit:
				stop = true
			default:
				stream.Read()
				m := max(in)
				now := time.Now()
				if (m-prev) > 0.1 && now.Sub(prevTime) > time.Millisecond*10 {
					prevTime = now
					cb(prevTime)
				}
				prev = m
			}
		}
		wg.Done()
	}()

	return nil
}

func Stop() {
	quit <- sig
	wg.Wait()
	portaudio.Terminate()
	stream.Close()
}

func max(in []float32) float32 {
	var out float32
	for _, f := range in {
		if f > out {
			out = f
		}
	}
	return out
}
