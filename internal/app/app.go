package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/slavka/another-less/internal/ingest"
	"github.com/slavka/another-less/internal/search"
)

type mode int

const (
	modeNormal mode = iota
	modeSearch
)

const horizontalScrollStep = 8

// Model is the top-level Bubble Tea model.
type Model struct {
	lines   []string
	linesCh <-chan ingest.LinesMsg
	done    bool // input stream exhausted

	// viewport
	width  int
	height int
	top    int // first visible line index
	left   int // horizontal scroll offset

	// search
	search      *search.State
	searchInput textinput.Model
	mode        mode

	// pipe mode: if true, Q writes buffer to stdout
	pipeOutput bool
	pipeBuffer string

	keys  KeyMap
	theme Theme
}

// Config holds initialization parameters for the app model.
type Config struct {
	LinesCh    <-chan ingest.LinesMsg
	PipeOutput bool
}

// New creates a new app model.
func New(cfg Config) Model {
	ti := textinput.New()
	ti.Prompt = "/"
	ti.CharLimit = 256

	return Model{
		linesCh:     cfg.LinesCh,
		pipeOutput:  cfg.PipeOutput,
		keys:        VimKeys(),
		theme:       DefaultTheme(),
		searchInput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return waitForLines(m.linesCh)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 1 // reserve 1 line for status bar
		m.clampViewport()
		return m, nil

	case ingest.LinesMsg:
		if msg.Lines != nil {
			m.lines = append(m.lines, msg.Lines...)
		}
		if msg.Done {
			m.done = true
			return m, nil
		}
		return m, waitForLines(m.linesCh)

	case tea.KeyMsg:
		if m.mode == modeSearch {
			return m.updateSearchMode(msg)
		}
		return m.updateNormalMode(msg)
	}

	return m, nil
}

func (m Model) updateNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keys.PipeQuit):
		if m.pipeOutput {
			m.pipeBuffer = strings.Join(m.lines, "\n")
		}
		return m, tea.Quit

	case key.Matches(msg, m.keys.Down):
		m.top++
		m.clampViewport()

	case key.Matches(msg, m.keys.Up):
		m.top--
		m.clampViewport()

	case key.Matches(msg, m.keys.Right):
		m.left += horizontalScrollStep
		m.clampViewport()

	case key.Matches(msg, m.keys.Left):
		m.left -= horizontalScrollStep
		m.clampViewport()

	case key.Matches(msg, m.keys.HalfPageDn):
		m.top += m.height / 2
		m.clampViewport()

	case key.Matches(msg, m.keys.HalfPageUp):
		m.top -= m.height / 2
		m.clampViewport()

	case key.Matches(msg, m.keys.PageDown):
		m.top += m.height
		m.clampViewport()

	case key.Matches(msg, m.keys.PageUp):
		m.top -= m.height
		m.clampViewport()

	case key.Matches(msg, m.keys.Top):
		m.top = 0
		m.left = 0

	case key.Matches(msg, m.keys.Bottom):
		m.top = len(m.lines) - m.height
		m.clampViewport()

	case key.Matches(msg, m.keys.LineStart):
		m.left = 0

	case key.Matches(msg, m.keys.LineEnd):
		m.left = m.maxLineWidth() - m.width
		m.clampViewport()

	case key.Matches(msg, m.keys.Search):
		m.mode = modeSearch
		m.searchInput.SetValue("")
		m.searchInput.Focus()
		return m, textinput.Blink

	case key.Matches(msg, m.keys.NextMatch):
		if m.search != nil {
			if idx := m.search.FindNextMatch(m.lines, m.top); idx >= 0 {
				m.top = idx
				m.clampViewport()
			}
		}

	case key.Matches(msg, m.keys.PrevMatch):
		if m.search != nil {
			if idx := m.search.FindPrevMatch(m.lines, m.top); idx >= 0 {
				m.top = idx
				m.clampViewport()
			}
		}
	}

	return m, nil
}

func (m Model) updateSearchMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		pattern := m.searchInput.Value()
		m.search = search.New(pattern)
		m.mode = modeNormal
		m.searchInput.Blur()
		// Jump to first match from current position
		if m.search != nil {
			if idx := m.search.FindNextMatch(m.lines, m.top-1); idx >= 0 {
				m.top = idx
				m.clampViewport()
			}
		}
		return m, nil

	case tea.KeyEsc:
		m.mode = modeNormal
		m.searchInput.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}

func (m *Model) clampViewport() {
	maxTop := len(m.lines) - m.height
	if maxTop < 0 {
		maxTop = 0
	}
	if m.top > maxTop {
		m.top = maxTop
	}
	if m.top < 0 {
		m.top = 0
	}
	if m.left < 0 {
		m.left = 0
	}
}

func (m Model) maxLineWidth() int {
	max := 0
	for _, line := range m.lines {
		if len(line) > max {
			max = len(line)
		}
	}
	return max
}

// PipeBuffer returns the text to write to stdout when pipe-quitting.
func (m Model) PipeBuffer() string {
	return m.pipeBuffer
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	var b strings.Builder

	viewHeight := m.height
	end := m.top + viewHeight
	if end > len(m.lines) {
		end = len(m.lines)
	}

	for i := m.top; i < end; i++ {
		line := m.lines[i]

		// Apply horizontal scroll
		if m.left < len(line) {
			line = line[m.left:]
		} else {
			line = ""
		}

		// Truncate to viewport width
		if len(line) > m.width {
			line = line[:m.width]
		}

		// Apply search highlighting
		if m.search != nil && m.search.Active {
			line = m.highlightLine(line)
		}

		b.WriteString(line)

		if i < end-1 {
			b.WriteByte('\n')
		}
	}

	// Pad empty lines
	rendered := end - m.top
	for i := rendered; i < viewHeight; i++ {
		b.WriteString(m.theme.LineNumber.Render("~"))
		if i < viewHeight-1 {
			b.WriteByte('\n')
		}
	}

	// Status bar
	b.WriteByte('\n')
	b.WriteString(m.renderStatusBar())

	return b.String()
}

func (m Model) highlightLine(line string) string {
	matches := m.search.FindMatches(line)
	if len(matches) == 0 {
		return line
	}

	var b strings.Builder
	prev := 0
	for _, match := range matches {
		if match.Start > prev {
			b.WriteString(line[prev:match.Start])
		}
		b.WriteString(m.theme.MatchHL.Render(line[match.Start:match.End]))
		prev = match.End
	}
	if prev < len(line) {
		b.WriteString(line[prev:])
	}
	return b.String()
}

func (m Model) renderStatusBar() string {
	total := len(m.lines)
	pos := ""
	if total == 0 {
		pos = "empty"
	} else if m.top+m.height >= total {
		pos = "END"
	} else if m.top == 0 {
		pos = "TOP"
	} else {
		pct := (m.top * 100) / total
		pos = fmt.Sprintf("%d%%", pct)
	}

	status := fmt.Sprintf(" %d lines | %s", total, pos)
	if !m.done {
		status += " | streaming..."
	}
	if m.search != nil && m.search.Active {
		status += fmt.Sprintf(" | /%s", m.search.Pattern)
	}

	if m.mode == modeSearch {
		prompt := m.searchInput.View()
		return lipgloss.PlaceHorizontal(m.width, lipgloss.Left, prompt)
	}

	return m.theme.StatusBar.Width(m.width).Render(status)
}
