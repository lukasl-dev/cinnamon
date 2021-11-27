package cinnamon

// LexerOptions holds the configuration options for the Lexer.
type LexerOptions struct {
	// Prefixes is a slice of strings that are used as prefixes. If the prefixes
	// are empty, the lexer will not use any prefixes.
	Prefixes []string `json:"prefixes,omitempty"`

	// PrefixIgnoreCase specifies whether the prefixes are case-insensitive.
	PrefixIgnoreCase bool `json:"prefixIgnoreCase,omitempty"`

	// Labels is a slice of strings that are used as labels. If the labels are
	// empty, the lexer will not use any labels.
	Labels []string `json:"labels,omitempty"`

	// LabelIgnoreCase specifies whether the labels are case-insensitive.
	LabelIgnoreCase bool `json:"labelIgnoreCase,omitempty"`

	// NoFlags indicates whether the lexer should handle flags as arguments.
	NoFlags bool `json:"noFlags,omitempty"`
}
