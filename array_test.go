package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetByIndexOrNil(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		slice []int
		index int
		want  []int
	}{
		{
			name:  "Index within range",
			slice: []int{1, 2, 3, 4, 5},
			index: 2,
			want:  []int{3, 4, 5},
		},
		{
			name:  "Index out of range",
			slice: []int{1, 2, 3, 4, 5},
			index: 7,
			want:  nil,
		},
		{
			name:  "Empty slice",
			slice: nil,
			index: 0,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := SliceFromIndexOrNil(tt.slice, tt.index)
			require.Equal(t, tt.want, actual)
		})
	}
}

func TestMergeWithoutDuplicates(t *testing.T) {
	t.Parallel()

	type args struct {
		a []string
		b []string
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "should work as expected",
			args: args{
				a: []string{"a", "b", "a", "c"},
				b: []string{"a", "d", "c", "e"},
			},
			want: []string{"a", "b", "a", "c", "d", "e"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := MergeWithoutDuplicates(tt.args.a, tt.args.b)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestSafeAppendToSlice(t *testing.T) {
	t.Parallel()

	t.Run("append to nil slice", func(t *testing.T) {
		t.Parallel()

		var nilSlice *[]int

		result := SafeAppendToSlice(nilSlice, 1)

		require.NotNil(t, result)
		require.Equal(t, []int{1}, *result)
	})

	t.Run("append to existing slice", func(t *testing.T) {
		t.Parallel()

		existingSlice := &[]int{1, 2}
		result := SafeAppendToSlice(existingSlice, 3)

		require.Equal(t, []int{1, 2, 3}, *result)
		require.Same(t, existingSlice, result)
	})
}
