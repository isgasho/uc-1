// The implementation of this package is heavily inspired by Rob Pike's amazing
// talk titled "Lexical Scanning in Go" [1].
//
// [1]: https://www.youtube.com/watch?v=HxaD_trXwRE

// Package lexer implements lexical tokenization of the µC programming language.
// While breaking the input into tokens, the next token is the longest sequence
// of characters that form a valid token.
package lexer

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/mewmew/uc/token"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Parse lexes the input read from r into a slice of tokens. Potential errors
// related to lexing are recorded as error tokens with relevant position
// information.
func Parse(r io.Reader) ([]token.Token, error) {
	br := bufio.NewReader(r)
	ur := newUnicodeReader(br)
	buf, err := ioutil.ReadAll(ur)
	if err != nil {
		return nil, err
	}
	input := string(buf)
	return ParseString(input), nil
}

// ParseFile lexes the input read from path into a slice of tokens. Potential
// errors related to lexing are recorded as error tokens with relevant position
// information.
func ParseFile(path string) ([]token.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

// ParseString lexes the input string into a slice of tokens. Potential errors
// related to lexing are recorded as error tokens with relevant position
// information.
func ParseString(s string) []token.Token {
	l := &lexer{
		input: s,
		// The average token size is roughly 2.5 for the quiet and noisy uC test
		// cases and 3.5 for the quiet, noisy and incorrect uC test cases.
		// Therefore pre-allocate a slice of tokens capable of holding all tokens
		// with the smallest average size.
		tokens: make([]token.Token, 0, len(s)/2),
	}

	// Tokenize the input.
	l.lex()

	return l.tokens
}

// A lexer lexes an input string into a slice of tokens.
type lexer struct {
	// The input string.
	input string
	// Start position of the current token.
	start int
	// Current position in the input.
	cur int
	// Width in byte of the last rune read with next; used by backup.
	width int
	// A slice of scanned tokens.
	tokens []token.Token
}

// lex lexes the input by repeatedly executing the active state function until
// it returns a nil state.
func (l *lexer) lex() {
	// lexToken is the initial state function of the lexer.
	for state := lexToken; state != nil; {
		state = state(l)
	}
}

// TODO: Decide which emit functions should stay (based on usage) after the
// implementation is mature.

// errorf appends an error token at the current token start position.
func (l *lexer) errorf(format string, args ...interface{}) {
	err := fmt.Sprintf(format, args...)
	tok := token.Token{
		Kind: token.Error,
		Val:  err,
		Pos:  l.start,
	}
	l.tokens = append(l.tokens, tok)
}

// errorfCur appends an error token at the current position.
func (l *lexer) errorfCur(format string, args ...interface{}) {
	err := fmt.Sprintf(format, args...)
	_, width := utf8.DecodeLastRuneInString(l.input[:l.cur])
	tok := token.Token{
		Kind: token.Error,
		Val:  err,
		Pos:  l.cur - width,
	}
	l.tokens = append(l.tokens, tok)
}

// emitErrorf emits an error token at the current token start position and
// advances the token start position.
func (l *lexer) emitErrorf(format string, args ...interface{}) {
	l.errorf(format, args...)
	l.ignore()
}

// emitEOF emits an EOF token at the current token start position. It emits an
// "unexpected EOF" error token if there exists unhandled input.
func (l *lexer) emitEOF() {
	if l.start < len(l.input) {
		panic(fmt.Sprintf("unexpected EOF; unhandled input %q", l.input[l.start:]))
	}
	l.emit(token.EOF)
}

// emit emits a token of the specified token type at the current token start
// position and advances the token start position.
func (l *lexer) emit(kind token.Kind) {
	l.emitCustom(kind, l.input[l.start:l.cur])
}

// emitCustom emits a custom token at the current token start position and
// advances the token start position.
func (l *lexer) emitCustom(kind token.Kind, val string) {
	tok := token.Token{
		Kind: kind,
		Val:  val,
		Pos:  l.start,
	}
	l.tokens = append(l.tokens, tok)
	// Advance the token start position.
	l.ignore()
}

// eof is the rune returned by next when no more input is available.
const eof = -1

// next consumes and returns the next rune of the input.
func (l *lexer) next() (r rune) {
	if l.cur >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.cur:])
	l.cur += l.width
	return r
}

// backup backs up one rune in the input. It can only be called once per call to
// next.
func (l *lexer) backup() {
	if l.width == 0 {
		panic("invalid call to backup; no matching call to next")
	}
	l.cur -= l.width
	l.width = 0
}

// accept consumes the next rune if it's from the valid set. It returns true if
// a rune was consumed and false otherwise.
func (l *lexer) accept(valid string) bool {
	r := l.next()
	if r == eof {
		return false
	}
	if strings.IndexRune(valid, r) == -1 {
		l.backup()
		return false
	}
	return true
}

// acceptRun consumes a run of runes from the valid set. It returns true if a
// rune was consumed and false otherwise.
func (l *lexer) acceptRun(valid string) bool {
	consumed := false
	for l.accept(valid) {
		consumed = true
	}
	return consumed
}

// ignore ignores any pending input read since the last token by advancing the
// token start position to the current position.
func (l *lexer) ignore() {
	l.start = l.cur
}

// ignoreRun ignores a run of valid runes.
func (l *lexer) ignoreRun(valid string) {
	if l.acceptRun(valid) {
		l.ignore()
	}
}

// newUnicodeReader wraps r to decode Unicode to UTF-8 as its reads.
func newUnicodeReader(r io.Reader) io.Reader {
	// fallback to r if no BOM sequence is located in the source text.
	t := unicode.BOMOverride(transform.Nop)
	return transform.NewReader(r, t)
}
