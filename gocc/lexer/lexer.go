// generated by gocc; DO NOT EDIT.

package lexer

import (

	// "fmt"
	// "github.com/mewmew/uc/gocc/util"

	"io/ioutil"
	"unicode/utf8"

	"github.com/mewmew/uc/gocc/token"
)

const (
	NoState    = -1
	NumStates  = 64
	NumSymbols = 82
)

type Lexer struct {
	src    []byte
	pos    int
	line   int
	column int
}

func NewLexer(src []byte) *Lexer {
	lexer := &Lexer{
		src:    src,
		pos:    0,
		line:   1,
		column: 1,
	}
	return lexer
}

func NewLexerFile(fpath string) (*Lexer, error) {
	src, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return NewLexer(src), nil
}

func (this *Lexer) Scan() (tok *token.Token) {

	// fmt.Printf("Lexer.Scan() pos=%d\n", this.pos)

	tok = new(token.Token)
	if this.pos >= len(this.src) {
		tok.Type = token.EOF
		tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = this.pos, this.line, this.column
		return
	}
	start, end := this.pos, 0
	tok.Type = token.INVALID
	state, rune1, size := 0, rune(-1), 0
	for state != -1 {

		// fmt.Printf("\tpos=%d, line=%d, col=%d, state=%d\n", this.pos, this.line, this.column, state)

		if this.pos >= len(this.src) {
			rune1 = -1
		} else {
			rune1, size = utf8.DecodeRune(this.src[this.pos:])
			this.pos += size
		}
		switch rune1 {
		case '\n':
			this.line++
			this.column = 1
		case '\r':
			this.column = 1
		case '\t':
			this.column += 4
		default:
			this.column++
		}

		// Production start
		if rune1 != -1 {
			state = TransTab[state](rune1)
		} else {
			state = -1
		}
		// Production end

		// Debug start
		// nextState := -1
		// if rune1 != -1 {
		// 	nextState = TransTab[state](rune1)
		// }
		// fmt.Printf("\tS%d, : tok=%s, rune == %s(%x), next state == %d\n", state, token.TokMap.Id(tok.Type), util.RuneToString(rune1), rune1, nextState)
		// fmt.Printf("\t\tpos=%d, size=%d, start=%d, end=%d\n", this.pos, size, start, end)
		// if nextState != -1 {
		// 	fmt.Printf("\t\taction:%s\n", ActTab[nextState].String())
		// }
		// state = nextState
		// Debug end

		if state != -1 {
			switch {
			case ActTab[state].Accept != -1:
				tok.Type = ActTab[state].Accept
				// fmt.Printf("\t Accept(%s), %s(%d)\n", string(act), token.TokMap.Id(tok), tok)
				end = this.pos
			case ActTab[state].Ignore != "":
				// fmt.Printf("\t Ignore(%s)\n", string(act))
				start = this.pos
				state = 0
				if start >= len(this.src) {
					tok.Type = token.EOF
				}

			}
		} else {
			if tok.Type == token.INVALID {
				end = this.pos
			}
		}
	}
	if end > start {
		this.pos = end
		tok.Lit = this.src[start:end]
	} else {
		tok.Lit = []byte{}
	}
	tok.Pos.Offset = start
	tok.Pos.Column = this.column
	tok.Pos.Line = this.line
	return
}

func (this *Lexer) Reset() {
	this.pos = 0
}

/*
Lexer symbols:
0: '''
1: '"'
2: '''
3: ';'
4: '('
5: ')'
6: '['
7: ']'
8: 't'
9: 'y'
10: 'p'
11: 'e'
12: 'd'
13: 'e'
14: 'f'
15: ','
16: 'r'
17: 'e'
18: 't'
19: 'u'
20: 'r'
21: 'n'
22: '{'
23: '}'
24: 'i'
25: 'f'
26: 'e'
27: 'l'
28: 's'
29: 'e'
30: 'w'
31: 'h'
32: 'i'
33: 'l'
34: 'e'
35: '='
36: '&'
37: '&'
38: '='
39: '='
40: '!'
41: '='
42: '<'
43: '>'
44: '<'
45: '='
46: '>'
47: '='
48: '+'
49: '-'
50: '*'
51: '/'
52: '!'
53: '_'
54: '/'
55: '/'
56: '\n'
57: '#'
58: '\n'
59: '/'
60: '*'
61: '*'
62: '*'
63: '/'
64: '\'
65: 'n'
66: ' '
67: '\t'
68: '\v'
69: '\f'
70: '\r'
71: '\n'
72: \u0001-'\t'
73: '\v'-'\f'
74: \u000e-'!'
75: '#'-'&'
76: '('-'['
77: ']'-\u007f
78: 'a'-'z'
79: 'A'-'Z'
80: '0'-'9'
81: .

*/
