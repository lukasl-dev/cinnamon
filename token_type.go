package cinnamon

//go:generate stringer -type=TokenType -trimprefix=TokenType -output token_type_string.go

// TokenType is used to differentiate between tokens.
type TokenType uint8

const (
	// TokenTypeInvalid is the zero value of TokenType. It is used to indicate
	// that a token is invalid.
	TokenTypeInvalid TokenType = iota

	// TokenTypeEOF is used to indicate that the end of the input has been
	// reached.
	TokenTypeEOF

	// TokenTypePrefix is the token type for prefixes.
	TokenTypePrefix

	// TokenTypeLabel is the token type for labels.
	TokenTypeLabel

	// TokenTypeFlag is the token type for flag name. For example: '--help', '-v'
	TokenTypeFlag

	// TokenTypeArgument is the token type for arguments. Arguments are space-separated
	// strings (spaces included).
	TokenTypeArgument

	// TokenTypeOpeningQuote is the token type for opening quotes. It indicates
	// that a quote is expected.
	TokenTypeOpeningQuote

	// TokenTypeEndQuote is the token type for closing quotes. It indicates
	// that a quote ends.
	TokenTypeEndQuote
)
