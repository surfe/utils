package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSafeString(t *testing.T) {
	t.Parallel()

	s := "A"
	n := 1

	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{
			name:  "Type String",
			input: s,
			want:  "A",
		},
		{
			name:  "Type Number",
			input: n,
			want:  "1",
		},
		{
			name:  "Type String Ptr",
			input: &s,
			want:  "A",
		},
		{
			name:  "Type Number Ptr",
			input: &n,
			want:  "1",
		},
		{
			name:  "Nil",
			input: nil,
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, SafeString(tt.input))
		})
	}
}

func TestGetFloat(t *testing.T) {
	t.Parallel()

	var (
		intValue             = 1
		float32Value float32 = 0.4
		float64Value         = 0.5
	)

	tests := []struct {
		name       string
		value      any
		want       float64
		isErrorNil bool
	}{
		{
			name:       "value is float64",
			value:      float64Value,
			want:       float64Value,
			isErrorNil: true,
		},
		{
			name:       "value is int",
			value:      intValue,
			want:       float64(intValue),
			isErrorNil: true,
		},
		{
			name:       "value is float32",
			value:      float32Value,
			want:       float64(float32Value),
			isErrorNil: true,
		},
		{
			name:       "value is not a number",
			value:      "hello",
			want:       0,
			isErrorNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := GetFloat(tt.value)
			require.InDelta(t, tt.want, got, 0.000001)
			require.Equal(t, tt.isErrorNil, err == nil)
		})
	}
}

func TestFindLooseStringInSlice(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		slice       []string
		item        string
		expectedIdx int
	}{
		{
			name:        "Item found (case-sensitive match)",
			slice:       []string{"apple", "banana", "cherry", "date"},
			item:        "banana",
			expectedIdx: 1,
		},
		{
			name:        "Item found (case-insensitive match)",
			slice:       []string{"apple", "banana", "cherry", "date"},
			item:        "BaNaNa",
			expectedIdx: 1,
		},
		{
			name:        "Item not found",
			slice:       []string{"apple", "banana", "cherry", "date"},
			item:        "grape",
			expectedIdx: -1,
		},
		{
			name:        "Item found with space separators (case-sensitive match)",
			slice:       []string{"Hello World", "foo bar", "Test Item"},
			item:        "foo bar",
			expectedIdx: 1,
		},
		{
			name:        "Item found with space separators (case-insensitive match)",
			slice:       []string{"Hello World", "foo bar", "Test Item"},
			item:        "Foo Bar",
			expectedIdx: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			index := FindLooseStringInSlice(tc.slice, tc.item)
			if index != tc.expectedIdx {
				t.Errorf("Expected index %d, but got %d", tc.expectedIdx, index)
			}
		})
	}
}

func TestGetSafeType(t *testing.T) {
	t.Parallel()

	type T struct {
		SomeString  string
		SomeInt     int
		SomePointer *T
	}

	tests := []struct {
		name string
		ptr  *T
		want T
	}{
		{
			name: "Returns empty value of type if nil is passes",
			ptr:  nil,
			want: T{},
		},
		{
			name: "returns dereferenced value if passed pointer is not nil",
			ptr: &T{
				SomeString:  "some-test",
				SomeInt:     1949,
				SomePointer: &T{},
			},
			want: T{
				SomeString:  "some-test",
				SomeInt:     1949,
				SomePointer: &T{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GetSafeType(tt.ptr)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetStringSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input interface{}
		want  []string
	}{
		{
			name:  "Null",
			input: nil,
			want:  []string{""},
		},
		{
			name:  "Strings slice",
			input: interface{}([]string{"val1", "val2"}),
			want:  []string{"val1", "val2"},
		},
		{
			name:  "Single string",
			input: interface{}("val1"),
			want:  []string{"val1"},
		},
		{
			name:  "Interface slice",
			input: interface{}([]interface{}{float64(10), float64(11), float64(12)}),
			want:  []string{"10", "11", "12"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GetStringSlice(tt.input)

			require.Equal(t, tt.want, got)
		})
	}
}
func TestGetPointerOrNil(t *testing.T) {
	t.Parallel()

	type testCase[T comparable] struct {
		name     string
		input    T
		expected *T
	}

	stringCases := []testCase[string]{
		{
			name:     "Non-empty string returns pointer",
			input:    "hello",
			expected: PtrTo("hello"),
		},
		{
			name:     "Empty string returns nil",
			input:    "",
			expected: nil,
		},
	}

	intCases := []testCase[int]{
		{
			name:     "Non-zero int returns pointer",
			input:    42,
			expected: PtrTo(42),
		},
		{
			name:     "Zero int returns nil",
			input:    0,
			expected: nil,
		},
	}

	boolCases := []testCase[bool]{
		{
			name:     "True returns pointer",
			input:    true,
			expected: PtrTo(true),
		},
		{
			name:     "False returns nil",
			input:    false,
			expected: nil,
		},
	}

	for _, tc := range stringCases {
		t.Run("string/"+tc.name, func(t *testing.T) {
			t.Parallel()

			got := GetPointerOrNil(tc.input)
			require.Equal(t, tc.expected, got)
		})
	}

	for _, tc := range intCases {
		t.Run("int/"+tc.name, func(t *testing.T) {
			t.Parallel()

			got := GetPointerOrNil(tc.input)
			require.Equal(t, tc.expected, got)
		})
	}

	for _, tc := range boolCases {
		t.Run("bool/"+tc.name, func(t *testing.T) {
			t.Parallel()

			got := GetPointerOrNil(tc.input)
			require.Equal(t, tc.expected, got)
		})
	}
}
