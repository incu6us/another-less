package ingest

import (
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	input := "line1\nline2\nline3\n"
	r := NewReader(strings.NewReader(input))
	r.Start()

	var allLines []string
	for msg := range r.Chan() {
		if msg.Done {
			break
		}
		allLines = append(allLines, msg.Lines...)
	}

	if len(allLines) != 3 {
		t.Fatalf("got %d lines, want 3", len(allLines))
	}
	expected := []string{"line1", "line2", "line3"}
	for i, line := range allLines {
		if line != expected[i] {
			t.Errorf("line[%d] = %q, want %q", i, line, expected[i])
		}
	}
}

func TestReaderEmpty(t *testing.T) {
	r := NewReader(strings.NewReader(""))
	r.Start()

	for msg := range r.Chan() {
		if msg.Lines != nil {
			t.Errorf("expected no lines, got %d", len(msg.Lines))
		}
		if msg.Done {
			break
		}
	}
}

func TestReaderBatching(t *testing.T) {
	// Create input larger than batchSize
	var b strings.Builder
	for i := 0; i < batchSize+10; i++ {
		b.WriteString("line\n")
	}
	r := NewReader(strings.NewReader(b.String()))
	r.Start()

	var totalLines int
	var batchCount int
	for msg := range r.Chan() {
		if msg.Done {
			break
		}
		batchCount++
		totalLines += len(msg.Lines)
	}

	if totalLines != batchSize+10 {
		t.Errorf("got %d total lines, want %d", totalLines, batchSize+10)
	}
	if batchCount < 2 {
		t.Errorf("expected at least 2 batches, got %d", batchCount)
	}
}
