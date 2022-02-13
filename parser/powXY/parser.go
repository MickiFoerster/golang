// Parser for sequence of assignments
package main

import (
	"fmt"
)

/*
Grammar:
seq : assign seq
    | assign
assign : IDENTIFIER = expr \n
expr : NUMBER OP expr
     | IDENTIFIER OP expr
     | NUMBER
	 | IDENTIFIER
*/

func parse(buf []byte) error {
	tok, buf, err := seq(buf)
	if err != nil {
		return err
	}
	fmt.Printf("Parsing result: %q\n", tok.val)

	if len(buf) > 0 {
		return fmt.Errorf("Parsing consumed not the whole buffer: %s", buf)
	}
	return nil
}

func next(buf []byte) (token, []byte, error) {
	tok, newbuf := lexer(buf)
	if tok.id == UNDEF {
		return token{id: UNDEF, val: ""}, nil, fmt.Errorf("Illegal token found: %q\n", buf)
	}

	return tok, newbuf, nil
}

func seq(buf []byte) (token, []byte, error) {
	tok, newbuf, err := assign(buf)
	if err != nil {
		return token{id: UNDEF}, buf, err
	}
	if tok.id != ASSIGN {
		return token{id: UNDEF}, buf, fmt.Errorf("Expected token ASSIGN, instead: %v", tok)
	}
	buf = newbuf
	seq_val := tok.val

	// Check if there is a longer sequence
	for len(buf) > 0 {
		tok, newbuf, err := next(buf)
		if err != nil {
			return token{id: UNDEF}, buf, err
		}
		if tok.id != WS {
			break
		}
		buf = newbuf
	}

	if len(buf) == 0 {
		return token{id: SEQ, val: seq_val}, buf, nil
	}

	tok, newbuf, err = seq(buf)
	if err != nil {
		return token{id: UNDEF}, buf, err
	}
	if tok.id != SEQ {
		return token{id: UNDEF}, buf, fmt.Errorf("Illegal sequence found: %v", buf)
	}
	buf = newbuf
	seq_val += tok.val

	return token{id: SEQ, val: seq_val}, buf, nil
}

func assign(buf []byte) (token, []byte, error) {
	assign_val := ""
	// IDENTIFIER
	tok, newbuf, err := next(buf)
	if err != nil {
		return token{id: UNDEF}, buf, err
	}
	if tok.id != IDENTIFIER {
		return token{id: UNDEF}, buf, fmt.Errorf("Expected token IDENTIFIER, instead: %v", tok)
	}
	buf = newbuf
	assign_val += tok.val

	// EQ
	tok, newbuf, err = next(buf)
	if err != nil {
		return token{id: UNDEF}, buf, err
	}

	if tok.id != EQ {
		return token{id: UNDEF}, buf, fmt.Errorf("Expected token '=', instead: %v", tok)
	}
	buf = newbuf
	assign_val += tok.val

	// expr
	tok, newbuf, err = expr(buf)
	if err != nil {
		return token{id: UNDEF}, buf, err
	}
	buf = newbuf
	assign_val += tok.val

	// WS
	tok, newbuf, err = next(buf)
	if err != nil {
		return token{id: UNDEF}, buf, err
	}
	if tok.id != WS {
		return token{id: UNDEF}, buf, fmt.Errorf("Expected whitespace token, instead: %v", tok)
	}
	buf = newbuf
	assign_val += ";"

	return token{id: ASSIGN, val: assign_val}, buf, nil
}

func expr(buf []byte) (token, []byte, error) {
	expr_val := ""
	// IDENTIFIER | NUMBER
	tok, newbuf, err := next(buf)
	if err != nil {
		return token{id: UNDEF}, buf, err
	}

	switch tok.id {
	case IDENTIFIER:
	case NUMBER:
	default:
		return token{id: UNDEF}, buf, fmt.Errorf("Unexpected token: %v", tok)
	}
	buf = newbuf
	expr_val += tok.val

	// Now look forward if next token is OP
	tok, newbuf, err = next(buf)
	if err != nil {
		return token{id: UNDEF}, buf, err
	}
	// If next token is not an OP then rewind buffer and return
	if tok.id != OP {
		return token{
			id:  EXPR,
			val: expr_val,
		}, buf, nil
	}
	buf = newbuf
	expr_val += tok.val

	// expr
	tok, newbuf, err = expr(buf)
	if err != nil {
		return token{id: UNDEF}, buf, err
	}
	buf = newbuf
	expr_val += tok.val

	return token{
		id:  EXPR,
		val: expr_val,
	}, buf, nil
}
