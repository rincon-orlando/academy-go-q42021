package model

import "fmt"

// Pokemon - General information about a Pokemon
type Pokemon struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// IsEven - Identifies whether a pokemon is ever or odd
func (p Pokemon) IsEven() bool {
	return p.ID%2 == 0
}

// String - Helps formatting a pokemon as string
func (p Pokemon) String() string {
	return fmt.Sprintf("ID = %d, Name = %s", p.ID, p.Name)
}
