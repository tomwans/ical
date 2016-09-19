// Package ical implements an extremely lazy iCal format decoder. I
// mean 'lazy' in the sense that, I really didn't put any effort into
// it. I didn't even read the spec. The thing just emits
// semi-structured data for you to use in Go.
//
// The parser itself is also lazy. You can pluck events out of an
// io.Reader lazily.
package ical

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Token represents an iCal token.
type Token struct {
	Type      string
	Value     string
	Subtokens []*Token
	Metadata  map[string]string
}

// Subtoken returns the first token that matches tokenType.
func (t *Token) Subtoken(tokenType string) *Token {
	for i := range t.Subtokens {
		if t.Subtokens[i].Type == tokenType {
			return t.Subtokens[i]
		}
	}
	return nil
}

// String outputs a human-friendly representation of the Token.
func (t *Token) String() string {
	if t == nil {
		return "<NIL: NIL>"
	}

	var attrs []string
	for i := range t.Subtokens {
		attr := t.Subtokens[i]
		var md []string
		for i := range attr.Metadata {
			md = append(md, fmt.Sprintf("%s=%s", i, attr.Metadata[i]))
		}
		var s string
		if len(md) > 0 {
			s = fmt.Sprintf("%s(%s)=%s", attr.Type, strings.Join(md, ","), attr.Value)
		} else {
			s = fmt.Sprintf("%s=%s", attr.Type, attr.Value)
		}
		attrs = append(attrs, s)
	}
	return fmt.Sprintf("<%s: %s>", t.Type, strings.Join(attrs, ", "))
}

// Decoder reads and decodes iCal tokens from an input stream.
type Decoder struct {
	src     io.Reader
	scanner *bufio.Scanner
}

func scanEntries(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	var ret []byte
	complete := false
	i := 0
	for {
		if len(data) == i {
			break
		}

		if data[i] == '\n' {
			// we only have enough up to the \n, but we don't have enough to
			// test to see if this line continues or not, so we need to ask
			// for more data.
			if i == len(data)-1 {
				// we need more
				return 0, nil, nil
			}

			// if the next thing isn't a space, we're done building our
			// token.
			if data[i+1] != ' ' {
				// we don't increase i because we want to
				// use this character in the next line.
				complete = true
				break
			}

			// here, we know that we have a \n + space combo, so we
			// skip 2 chars ahead for the next go-around
			i += 2
			continue
		}

		// we always omit carriage returns.
		if data[i] == '\r' {
			i++
			continue
		}

		// otherwise, this is a valid character for this token
		ret = append(ret, data[i])
		i++
	}

	if !complete {
		return 0, nil, nil
	}

	return i + 1, ret, nil
}

// NewDecoder returns a new Decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	scanner := bufio.NewScanner(r)
	scanner.Split(scanEntries)
	return &Decoder{src: r, scanner: scanner}
}

// Decode decodes the VCALENDAR token.
func (d *Decoder) Decode(tok *Token) error {
	parsedToken, err := d.NextToken("VCALENDAR")
	if err != nil && err != io.EOF {
		return err
	}
	tok.Value = parsedToken.Value
	tok.Type = parsedToken.Type
	tok.Subtokens = parsedToken.Subtokens
	tok.Metadata = parsedToken.Metadata
	return nil
}

// NextToken returns the next token specified by token type.
func (d *Decoder) NextToken(tokentype string) (*Token, error) {
	var stack []*Token
	for {
		// Try to scan through the file, looking for a new line.
		if !d.scanner.Scan() {
			// the file is over, so just emit whatever we found, starting at
			// the bottom of the stack.
			if len(stack) == 0 {
				return nil, io.EOF
			}
			return stack[0], io.EOF
		}

		line := d.scanner.Text()

		// we got nothing, try the next line.
		if line == "" {
			continue
		}

		pieces := strings.SplitN(line, ":", 2)

		if pieces[0] == "BEGIN" {
			// we're beginning a new token onto the stack.
			tok := &Token{}
			// what type is it?
			tok.Type = pieces[1]
			// push onto the stack.
			stack = append(stack, tok)
			continue
		}

		if pieces[0] == "END" {
			// we have an END token, so lets pop off the last entry in our
			// stack.
			tok := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			// is this a subtoken?
			if len(stack) > 0 {
				parentTok := stack[len(stack)-1]
				parentTok.Subtokens = append(parentTok.Subtokens, tok)
			}

			// for END lines, we just need to check
			// if this is the END we're looking for.
			if pieces[1] == tokentype {
				return tok, nil
			}

			continue
		}

		// we have some other value, which signals that it is a Attribute
		attr := &Token{}
		attr.Value = pieces[1]

		// check if there are any metadata
		metadata := strings.Split(pieces[0], ";")
		if len(metadata) > 0 {
			attr.Metadata = make(map[string]string)
			// skip the value before the ;
			attr.Type = metadata[0]
			metadata = metadata[1:]
			for i := range metadata {
				// split this by =
				tags := strings.Split(metadata[i], "=")
				attr.Metadata[tags[0]] = tags[1]
			}
		} else {
			attr.Type = pieces[0]
		}

		// place this attribute onto the token.
		parentTok := stack[len(stack)-1]
		parentTok.Subtokens = append(parentTok.Subtokens, attr)
	}
}
