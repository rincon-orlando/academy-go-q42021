package enum

import (
	"errors"
	"strings"
)

// OddEven - Works as enum to identify Odd or Even numbers
type OddEven int

const (
	Undefined OddEven = iota
	Odd
	Even
)

// ParseOddEven - Takes a string and returns an Odd or Even enum
func ParseOddEven(input string) (OddEven, error) {
	switch strings.ToLower(input) {
	case "odd":
		return Odd, nil
	case "even":
		return Even, nil
	}

	return Undefined, errors.New(input + " is not a valid input. Must be either 'odd' or 'even'")
}
