package cinnamon

import (
	"strings"
	"unicode"
)

// lexerStateFunc is a recursive function that returns the next state of the
// lexer or nil if the Lexer should stop.
type lexerStateFunc func() lexerStateFunc

// begin is the first state of the Lexer.
func (l *Lexer) begin() lexerStateFunc { return l.lexPrefix }

// lexPrefix reads the prefix from l's reader. If something goes wrong, nil is
// returned and TokenTypeEOF is sent to l's token channel.
func (l *Lexer) lexPrefix() lexerStateFunc {
	if len(l.opts.Prefixes) == 0 {
		return l.lexLabel
	}

	length, err := l.firstMatching(l.opts.Prefixes, l.opts.PrefixIgnoreCase)
	if err != nil {
		return l.eof("")
	}
	if length == 0 {
		return l.eof("prefix not found")
	}

	bs, err := l.readLimited(length)
	if len(bs) > 0 {
		l.send(TokenTypePrefix, string(bs))
	}
	if err != nil {
		return l.eof("")
	}

	return l.lexLabel
}

// lexPrefix reads the label from l's reader. If something goes wrong, nil is
// returned and TokenTypeEOF is sent to l's token channel.
func (l *Lexer) lexLabel() lexerStateFunc {
	if len(l.opts.Labels) == 0 {
		return l.lexRemaining
	}

	length, err := l.firstMatching(l.opts.Labels, l.opts.LabelIgnoreCase)
	if err != nil {
		return l.eof("")
	}
	if length == 0 {
		return l.eof("label not found")
	}

	bs, err := l.readLimited(length)
	if len(bs) > 0 {
		if r, _, _ := l.readRune(); unicode.IsSpace(r) {
			bs = append(bs, byte(r))
		} else {
			_ = l.unreadRune()
		}
		l.send(TokenTypeLabel, string(bs))
	}
	if err != nil {
		l.send(TokenTypeEOF, "")
		return nil
	}

	return l.lexRemaining
}

// lexRemaining reads the remaining bytes from l's reader. If l's reader returns
// an error, nil is returned and TokenTypeEOF is sent to l's token channel.
func (l *Lexer) lexRemaining() lexerStateFunc {
	r, _, err := l.readRune()
	if err != nil {
		return l.eof("")
	}
	_ = l.unreadRune()
	switch {
	case r == '-':
		return l.lexFlag
	case isQuote(r):
		return l.lexQuote
	default:
		return l.lexArgument
	}
}

// lexFlag reads a single flag from l's reader. If something goes wrong, nil is
// returned and TokenTypeEOF is sent to l's token channel.
func (l *Lexer) lexFlag() lexerStateFunc {
	runes, err := l.readRunes(func(r rune) bool { return unicode.IsSpace(r) })
	if len(runes) > 0 {
		l.send(TokenTypeFlag, string(runes))
	}
	if err != nil {
		return l.eof("")
	}
	return l.lexRemaining
}

func (l *Lexer) lexQuote() lexerStateFunc {
	r, _, err := l.readRune()
	if err != nil {
		return l.eof("")
	}
	switch {
	case l.openingQuote == 0:
		l.openingQuote = r
		l.send(TokenTypeOpeningQuote, string(r))
		return l.lexRemaining
	case l.openingQuote == r:
		l.openingQuote = 0
		l.send(TokenTypeClosingQuote, string(r))
		return l.lexRemaining
	default:
		_ = l.unreadRune()
		return l.lexArgument
	}
}

// lexArgument reads a single argument from l's reader. If something goes wrong,
// nil is returned and TokenTypeEOF is sent to l's token channel.
func (l *Lexer) lexArgument() lexerStateFunc {
	rs, err := l.readRunes(func(r rune) bool { return unicode.IsSpace(r) })
	spaces, _ := l.readUntilNonSpace()
	rs = append(rs, spaces...)

	if len(rs) > 0 {
		var quote rune

		trim := []rune(strings.TrimRight(string(rs), " "))
		if len(trim) > 0 && trim[len(trim)-1] == l.openingQuote {
			i := len(trim) - 1
			quote = rs[i]
			rs = rs[:i]
		}

		if len(rs) > 0 {
			l.send(TokenTypeArgument, string(rs))
		}
		if quote != 0 {
			l.send(TokenTypeClosingQuote, string(quote))
		}
	}
	if err != nil {
		return l.eof("")
	}

	return l.lexRemaining
}

func isQuote(r rune) bool {
	return r == '"' || r == '\'' || r == '`'
}
