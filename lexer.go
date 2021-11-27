package cinnamon

import (
	"bufio"
	"io"
	"strings"
	"sync"
	"unicode"
)

// Lexer is a lexical reader.
type Lexer struct {
	// mu prevents concurrent access.
	mu *sync.Mutex

	// opts are the pre-configured options.
	opts LexerOptions

	// rd is the reader to read from.
	rd *bufio.Reader

	// tokens is the channel of recent lexed tokens.
	tokens chan Token

	// openingQuote is the opening quote in the current state.
	openingQuote rune
}

// NewLexer returns a new Lexer that is configured by the first given LexerOptions.
func NewLexer(opts ...LexerOptions) *Lexer {
	l := &Lexer{mu: new(sync.Mutex)}
	if len(opts) > 0 {
		l.opts = opts[0]
	}
	return l
}

// Lex lets the Lexer read from rd and returns a channel of tokens.
func (l *Lexer) Lex(rd io.Reader) <-chan Token {
	l.reset(rd)
	go l.start()
	return l.tokens
}

// start resets l's state and starts the lexical reading.
func (l *Lexer) start() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for fn := l.begin(); fn != nil; {
		fn = fn()
	}
	close(l.tokens)
}

// reset resets l's state and sets rd as the reader to read from.
func (l *Lexer) reset(rd io.Reader) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.rd = bufio.NewReader(rd)
	l.tokens = make(chan Token)
	l.openingQuote = 0
}

// send sends a token to l's tokens channel.
func (l *Lexer) send(typ TokenType, val string) {
	l.tokens <- tokenOf(typ, val)
}

// eof sends TokenTypeEOF with val to l's token channel and returns nil.
func (l *Lexer) eof(val string) lexerStateFunc {
	l.send(TokenTypeEOF, val)
	return nil
}

// readRunes reads runes from l's reader until pred returns true or an error
// occurs.
func (l *Lexer) readRunes(pred func(r rune) bool) (runes []rune, err error) {
	for {
		r, _, err := l.rd.ReadRune()
		if err != nil {
			return runes, err
		}
		runes = append(runes, r)
		if pred(r) {
			return runes, nil
		}
	}
}

func (l *Lexer) readUntilNonSpace() (rs []rune, err error) {
	for {
		r, _, err := l.rd.ReadRune()
		if err != nil {
			return rs, err
		}
		if !unicode.IsSpace(r) {
			_ = l.unreadRune()
			return rs, nil
		}
		rs = append(rs, r)
	}
}

// readRune reads the next rune from l's reader.
func (l *Lexer) readRune() (r rune, size int, err error) {
	return l.rd.ReadRune()
}

// unreadRune unreads the last rune read from l's reader.
func (l *Lexer) unreadRune() error {
	return l.rd.UnreadRune()
}

// readLimited reads n bytes from l's reader.
func (l *Lexer) readLimited(n int) ([]byte, error) {
	return io.ReadAll(io.LimitReader(l.rd, int64(n)))
}

// firstMatching returns the length of the first matching string in available.
func (l *Lexer) firstMatching(available []string, ignoreCase bool) (int, error) {
	var length int
	for _, s := range available {
		bs, err := l.rd.Peek(len(s))
		if err != nil {
			return 0, err
		}
		if (ignoreCase && strings.EqualFold(string(bs), s)) || string(bs) == s {
			length = len(s)
			break
		}
	}
	return length, nil
}
