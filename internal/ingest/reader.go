package ingest

import (
	"bufio"
	"io"
	"os"
)

// LinesMsg carries a batch of lines read from input.
type LinesMsg struct {
	Lines []string
	Done  bool // true when the input stream is exhausted
}

// Reader reads lines from an io.Reader and sends batches on a channel.
type Reader struct {
	r    io.Reader
	ch   chan LinesMsg
	done chan struct{}
}

// NewReader creates a Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r:    r,
		ch:   make(chan LinesMsg, 16),
		done: make(chan struct{}),
	}
}

// NewFileReader opens a file and creates a Reader for it.
func NewFileReader(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return NewReader(f), nil
}

// Chan returns the channel that receives line batches.
func (rd *Reader) Chan() <-chan LinesMsg {
	return rd.ch
}

// Start begins reading lines in a goroutine. It batches lines for efficiency.
func (rd *Reader) Start() {
	go rd.run()
}

// Stop signals the reader to stop.
func (rd *Reader) Stop() {
	close(rd.done)
}

const batchSize = 1000

func (rd *Reader) run() {
	defer close(rd.ch)

	scanner := bufio.NewScanner(rd.r)
	scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024) // 10MB max line

	batch := make([]string, 0, batchSize)

	for scanner.Scan() {
		select {
		case <-rd.done:
			return
		default:
		}

		batch = append(batch, scanner.Text())

		if len(batch) >= batchSize {
			rd.ch <- LinesMsg{Lines: batch}
			batch = make([]string, 0, batchSize)
		}
	}

	// Send remaining lines
	if len(batch) > 0 {
		rd.ch <- LinesMsg{Lines: batch}
	}

	// Signal completion
	rd.ch <- LinesMsg{Done: true}
}
