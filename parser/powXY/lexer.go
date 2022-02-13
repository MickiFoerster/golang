// Lexer for simple assignments operations on right-handside
package main

type token struct {
	id  int
	val string
}

const (
	UNDEF = iota
	WS
	EQ
	IDENTIFIER
	NUMBER
	OP
	SEQ
	ASSIGN
	EXPR
)

func lexer(s []byte) (token, []byte) {
	c := s[0]

	if 'a' <= c && c <= 'z' {
		identifier := string(c)
		var i uint = 1
		for {
			c = s[i]
			if 'a' <= c && c <= 'z' {
				identifier += string(c)
				i++
				continue
			}
			break
		}
		return token{
			id:  IDENTIFIER,
			val: identifier,
		}, s[i:]
	}

	if '0' <= c && c <= '9' {
		number := string(c)
		var i uint = 1
		for {
			c = s[i]
			if '0' <= c && c <= '9' {
				number += string(c)
				i++
				continue
			}
			break
		}
		return token{
			id:  NUMBER,
			val: number,
		}, s[i:]
	}

	switch c {
	case '\n', '\t', ' ':
		return token{
			id:  WS,
			val: string(c),
		}, s[1:]
	case '=':
		return token{
			id:  EQ,
			val: string(c),
		}, s[1:]
	case '*':
		return token{
			id:  OP,
			val: string(c),
		}, s[1:]
	}

	return token{id: UNDEF, val: ""}, nil
}
