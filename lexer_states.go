package cinnamon

import (
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
		l.send(TokenTypeEOF, "")
		return nil
	}
	if length == 0 {
		l.send(TokenTypeEOF, "prefix not found")
		return nil
	}

	bs, err := l.readLimited(length)
	if len(bs) > 0 {
		l.send(TokenTypePrefix, string(bs))
	}
	if err != nil {
		l.send(TokenTypeEOF, "")
		return nil
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
		l.send(TokenTypeEOF, "")
		return nil
	}
	if length == 0 {
		l.send(TokenTypeEOF, "label not found")
		return nil
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
		l.send(TokenTypeEOF, "")
		return nil
	}
	_ = l.unreadRune()
	switch r {
	case '-':
		return l.lexFlag
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
		l.send(TokenTypeEOF, "")
		return nil
	}
	return l.lexRemaining
}

// lexArgument reads a single argument from l's reader. If something goes wrong,
// nil is returned and TokenTypeEOF is sent to l's token channel.
func (l *Lexer) lexArgument() lexerStateFunc {
	runes, err := l.readRunes(func(r rune) bool { return unicode.IsSpace(r) })
	if len(runes) > 0 {
		l.send(TokenTypeArgument, string(runes))
	}
	if err != nil {
		l.send(TokenTypeEOF, "")
		return nil
	}
	return l.lexRemaining
}
