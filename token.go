package cinnamon

import (
	"fmt"
)

// Token is a lexical token.
type Token struct {
	// Type is the type of the Token.
	Type TokenType `json:"type,omitempty"`

	// Value is the literal value.
	Value string `json:"value,omitempty"`
}

// String returns the string representation of t in the following format:
// <type> [<value> | "<empty>"]
func (t Token) String() string {
	value := t.Value
	if value == "" {
		value = "<empty>"
	}
	return fmt.Sprintf("%s: %s", t.Type, value)
}
