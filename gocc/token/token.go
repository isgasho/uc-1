// generated by gocc; DO NOT EDIT.

package token

import (
	"fmt"
)

type Token struct {
	Type
	Lit []byte
	Pos
}

type Type int

const (
	INVALID Type = iota
	EOF
)

type Pos struct {
	Offset int
	Line   int
	Column int
}

func (this Pos) String() string {
	return fmt.Sprintf("Pos(offset=%d, line=%d, column=%d)", this.Offset, this.Line, this.Column)
}

type TokenMap struct {
	typeMap []string
	idMap   map[string]Type
}

func (this TokenMap) Id(tok Type) string {
	if int(tok) < len(this.typeMap) {
		return this.typeMap[tok]
	}
	return "unknown"
}

func (this TokenMap) Type(tok string) Type {
	if typ, exist := this.idMap[tok]; exist {
		return typ
	}
	return INVALID
}

func (this TokenMap) TokenString(tok *Token) string {
	//TODO: refactor to print pos & token string properly
	return fmt.Sprintf("%s(%d,%s)", this.Id(tok.Type), tok.Type, tok.Lit)
}

func (this TokenMap) StringType(typ Type) string {
	return fmt.Sprintf("%s(%d)", this.Id(typ), typ)
}

var TokMap = TokenMap{
	typeMap: []string{
		"INVALID",
		"$",
		"empty",
		";",
		"ident",
		"(",
		")",
		"[",
		"]",
		"int_lit",
		"char_lit",
		"typedef",
		",",
		"return",
		"{",
		"}",
		"if",
		"else",
		"while",
		"=",
		"&&",
		"==",
		"!=",
		"<",
		">",
		"<=",
		">=",
		"+",
		"-",
		"*",
		"/",
		"!",
	},

	idMap: map[string]Type{
		"INVALID":  0,
		"$":        1,
		"empty":    2,
		";":        3,
		"ident":    4,
		"(":        5,
		")":        6,
		"[":        7,
		"]":        8,
		"int_lit":  9,
		"char_lit": 10,
		"typedef":  11,
		",":        12,
		"return":   13,
		"{":        14,
		"}":        15,
		"if":       16,
		"else":     17,
		"while":    18,
		"=":        19,
		"&&":       20,
		"==":       21,
		"!=":       22,
		"<":        23,
		">":        24,
		"<=":       25,
		">=":       26,
		"+":        27,
		"-":        28,
		"*":        29,
		"/":        30,
		"!":        31,
	},
}
