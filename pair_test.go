package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input map[string]int
		want  PairList
	}{
		{
			name:  "Input is a valid map",
			input: map[string]int{"key1": 10, "key2": 0, "key3": 2, "key4": 9, "key5": 6},
			want:  PairList{{"key2", 0}, {"key3", 2}, {"key5", 6}, {"key4", 9}, {"key1", 10}},
		},
		{
			name:  "Input is nil",
			input: nil,
			want:  PairList{},
		},
		{
			name:  "Input is empty",
			input: map[string]int{},
			want:  PairList{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := SortMap(tt.input)
			if assert.NotNil(t, got) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
