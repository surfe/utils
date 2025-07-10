package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestStructA struct {
	Field1 string
	Field2 int
	Field3 *bool
	Field4 []string
}

type TestStructB struct {
	Field1 string
	Field2 *int
	Field3 map[string]int
	Field4 float64
}

func TestCountNonEmptyFields(t *testing.T) {
	t.Parallel()

	t.Run("Struct with non-empty fields", func(t *testing.T) {
		t.Parallel()

		trueValue := true
		testStruct := TestStructA{
			Field1: "value",
			Field2: 42,
			Field3: &trueValue,
			Field4: []string{"item1", "item2"},
		}
		result := CountNonEmptyFields(testStruct)
		require.Equal(t, 4, result)
	})

	t.Run("Struct with some empty fields", func(t *testing.T) {
		t.Parallel()

		testStruct := TestStructA{
			Field1: "",
			Field2: 0,
			Field3: nil,
			Field4: []string{},
		}
		result := CountNonEmptyFields(testStruct)
		require.Equal(t, 0, result)
	})

	t.Run("Struct with pointer and map fields", func(t *testing.T) {
		t.Parallel()

		intValue := 10
		testStruct := TestStructB{
			Field1: "non-empty",
			Field2: &intValue,
			Field3: map[string]int{"key": 1},
			Field4: 0.0,
		}
		result := CountNonEmptyFields(testStruct)
		require.Equal(t, 3, result)
	})

	t.Run("Struct with all fields empty", func(t *testing.T) {
		t.Parallel()

		testStruct := TestStructB{}
		result := CountNonEmptyFields(testStruct)
		require.Equal(t, 0, result)
	})

	t.Run("Struct with an empty list", func(t *testing.T) {
		t.Parallel()

		testStruct := TestStructA{
			Field1: "",
			Field2: 0,
			Field3: nil,
			Field4: []string{},
		}
		result := CountNonEmptyFields(testStruct)
		require.Equal(t, 0, result)
	})
}
