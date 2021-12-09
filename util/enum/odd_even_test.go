package enum

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OddEven(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedResult OddEven
		hasError       bool
		error          error
	}{
		{
			name:           "test odd translation",
			input:          "odd",
			expectedResult: Odd,
			hasError:       false,
			error:          nil,
		},
		{
			name:           "test even translation",
			input:          "even",
			expectedResult: Even,
			hasError:       false,
			error:          nil,
		},
		{
			name:           "test mixed case translation",
			input:          "EvEn",
			expectedResult: Even,
			hasError:       false,
			error:          nil,
		},
		{
			name:           "test wrong case translation",
			input:          "Par",
			expectedResult: Undefined,
			hasError:       true,
			error:          errors.New("Par is not a valid input. Must be either 'odd' or 'even'"),
		},
	}

	for _, tc := range testCases {
		result, err := ParseOddEven(tc.input)
		assert.Equal(t, tc.expectedResult, result)
		if tc.hasError {
			assert.EqualError(t, err, tc.error.Error())
		}
	}
}
