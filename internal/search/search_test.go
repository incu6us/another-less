package search

import (
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("empty pattern returns nil", func(t *testing.T) {
		s := New("")
		if s != nil {
			t.Error("expected nil for empty pattern")
		}
	})

	t.Run("valid pattern", func(t *testing.T) {
		s := New("hello")
		if s == nil {
			t.Fatal("expected non-nil")
		}
		if s.Pattern != "hello" {
			t.Errorf("got pattern %q, want %q", s.Pattern, "hello")
		}
	})

	t.Run("invalid regex falls back to literal", func(t *testing.T) {
		s := New("[invalid")
		if s == nil {
			t.Fatal("expected non-nil")
		}
		if !s.HasMatch("[invalid") {
			t.Error("expected literal match")
		}
	})
}

func TestSmartCase(t *testing.T) {
	s := New("error")
	if !s.HasMatch("ERROR in line") {
		t.Error("lowercase pattern should match uppercase (smart case)")
	}

	s = New("Error")
	if s.HasMatch("error in line") {
		t.Error("mixed case pattern should be case-sensitive")
	}
	if !s.HasMatch("Error in line") {
		t.Error("mixed case pattern should match exact case")
	}
}

func TestFindMatches(t *testing.T) {
	s := New("ab")
	matches := s.FindMatches("abcabc")
	if len(matches) != 2 {
		t.Fatalf("got %d matches, want 2", len(matches))
	}
	if matches[0].Start != 0 || matches[0].End != 2 {
		t.Errorf("match[0] = %+v, want {0 2}", matches[0])
	}
	if matches[1].Start != 3 || matches[1].End != 5 {
		t.Errorf("match[1] = %+v, want {3 5}", matches[1])
	}
}

func TestFindNextMatch(t *testing.T) {
	lines := []string{"foo", "bar", "baz", "bar", "qux"}

	s := New("bar")
	idx := s.FindNextMatch(lines, 0)
	if idx != 1 {
		t.Errorf("got %d, want 1", idx)
	}

	idx = s.FindNextMatch(lines, 1)
	if idx != 3 {
		t.Errorf("got %d, want 3", idx)
	}

	// Wraps around
	idx = s.FindNextMatch(lines, 3)
	if idx != 1 {
		t.Errorf("got %d, want 1 (wrap)", idx)
	}
}

func TestFindPrevMatch(t *testing.T) {
	lines := []string{"foo", "bar", "baz", "bar", "qux"}

	s := New("bar")
	idx := s.FindPrevMatch(lines, 4)
	if idx != 3 {
		t.Errorf("got %d, want 3", idx)
	}

	// Wraps around
	idx = s.FindPrevMatch(lines, 1)
	if idx != 3 {
		t.Errorf("got %d, want 3 (wrap)", idx)
	}
}

func TestNilState(t *testing.T) {
	var s *State
	if s.HasMatch("anything") {
		t.Error("nil state should not match")
	}
	if s.FindNextMatch([]string{"a"}, 0) != -1 {
		t.Error("nil state should return -1")
	}
	if matches := s.FindMatches("test"); matches != nil {
		t.Error("nil state should return nil matches")
	}
}
