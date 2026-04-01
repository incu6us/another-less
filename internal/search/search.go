package search

import (
	"regexp"
	"strings"
)

// Match represents a single match position within a line.
type Match struct {
	Start int
	End   int
}

// State holds the current search state.
type State struct {
	Pattern   string
	Regex     *regexp.Regexp
	Active    bool
	CaseSmart bool // auto case-insensitive if pattern is all lowercase
}

// New creates a new search state from a pattern string.
// Returns nil if the pattern is empty.
func New(pattern string) *State {
	if pattern == "" {
		return nil
	}

	s := &State{
		Pattern:   pattern,
		CaseSmart: true,
		Active:    true,
	}

	// Smart case: if all lowercase, make case-insensitive
	actual := pattern
	if s.CaseSmart && actual == strings.ToLower(actual) {
		actual = "(?i)" + actual
	}

	re, err := regexp.Compile(actual)
	if err != nil {
		// Fall back to literal match
		re = regexp.MustCompile(regexp.QuoteMeta(pattern))
	}
	s.Regex = re
	return s
}

// FindMatches returns all match positions in a line.
func (s *State) FindMatches(line string) []Match {
	if s == nil || s.Regex == nil {
		return nil
	}

	locs := s.Regex.FindAllStringIndex(line, -1)
	matches := make([]Match, len(locs))
	for i, loc := range locs {
		matches[i] = Match{Start: loc[0], End: loc[1]}
	}
	return matches
}

// HasMatch returns true if the line contains a match.
func (s *State) HasMatch(line string) bool {
	if s == nil || s.Regex == nil {
		return false
	}
	return s.Regex.MatchString(line)
}

// FindNextMatch finds the index of the next line with a match starting from (but not including) the given index.
// It wraps around. Returns -1 if no match found.
func (s *State) FindNextMatch(lines []string, from int) int {
	if s == nil || s.Regex == nil {
		return -1
	}
	n := len(lines)
	for i := 1; i <= n; i++ {
		idx := (from + i) % n
		if s.Regex.MatchString(lines[idx]) {
			return idx
		}
	}
	return -1
}

// FindPrevMatch finds the index of the previous line with a match starting from (but not including) the given index.
// It wraps around. Returns -1 if no match found.
func (s *State) FindPrevMatch(lines []string, from int) int {
	if s == nil || s.Regex == nil {
		return -1
	}
	n := len(lines)
	for i := 1; i <= n; i++ {
		idx := (from - i + n) % n
		if s.Regex.MatchString(lines[idx]) {
			return idx
		}
	}
	return -1
}
