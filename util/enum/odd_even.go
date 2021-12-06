package enum

import (
	"errors"
	"strings"
)

type OddEven int

const (
	Undefined OddEven = iota
	Odd
	Even
)

func ParseOddEven(input string) (OddEven, error) {
	switch strings.ToLower(input) {
	case "odd":
		return Odd, nil
	case "even":
		return Even, nil
	}

	return Undefined, errors.New(input + " is not a valid input. Must be either 'odd' or 'even'")
}
