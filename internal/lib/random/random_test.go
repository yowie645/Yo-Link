package random

import (
	"strings"
	"testing"
)

func TestNewRandomStringLength(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"size 0", 0},
		{"size 1", 1},
		{"size 10", 10},
		{"size 100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRandomString(tt.size)
			if len(got) != tt.size {
				t.Errorf("NewRandomString(%d) length = %d, want %d", tt.size, len(got), tt.size)
			}
		})
	}
}

func TestNewRandomStringContent(t *testing.T) {
	const size = 100
	got := NewRandomString(size)

	for _, c := range got {
		if !strings.ContainsRune(letterBytes, c) {
			t.Errorf("NewRandomString(%d) contains invalid character %q", size, c)
		}
	}
}

func TestNewRandomStringUniqueness(t *testing.T) {
	const size = 20
	const iterations = 100
	results := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		s := NewRandomString(size)
		if results[s] {
			t.Errorf("NewRandomString(%d) produced duplicate string %q", size, s)
		}
		results[s] = true
	}
}

func BenchmarkNewRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewRandomString(32)
	}
}
