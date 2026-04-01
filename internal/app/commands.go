package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/incu6us/another-less/internal/ingest"
)

// waitForLines returns a tea.Cmd that waits for the next batch of lines from the reader.
func waitForLines(ch <-chan ingest.LinesMsg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return ingest.LinesMsg{Done: true}
		}
		return msg
	}
}
