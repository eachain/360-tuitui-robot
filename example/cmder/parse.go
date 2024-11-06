package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func ParseTokens(s string) ([]string, error) {
	var rs []rune
	var tokens []string
	var inSingleQuot, inDoubleQuot, inString, escape bool
	var pos, singlePos, doublePos, escapePos int

	t := s
	for s != "" {
		r, n := utf8.DecodeRuneInString(s)
		s = s[n:]
		pos += n

		if escape {
			switch r {
			case 'n':
				rs = append(rs, '\n')
			case 'r':
				rs = append(rs, '\r')
			case 't':
				rs = append(rs, '\t')
			default:
				rs = append(rs, r)
			}
			escape = false
			continue
		}

		if r == '\\' {
			if inSingleQuot {
				rs = append(rs, r)
				continue
			}
			escape = true
			escapePos = pos - n
			continue
		}

		if r == '\'' {
			if inSingleQuot {
				inSingleQuot = false
				if !inString {
					tokens = append(tokens, string(rs))
					rs = rs[:0]
				}
				continue
			}
			if inDoubleQuot {
				rs = append(rs, r)
				continue
			}
			inSingleQuot = true
			singlePos = pos - n
			continue
		}

		if r == '"' {
			if inDoubleQuot {
				inDoubleQuot = false
				if !inString {
					tokens = append(tokens, string(rs))
					rs = rs[:0]
				}
				continue
			}
			if inSingleQuot {
				rs = append(rs, r)
				continue
			}
			inDoubleQuot = true
			doublePos = pos - n
			continue
		}

		if unicode.IsSpace(r) {
			if inSingleQuot || inDoubleQuot || escape {
				rs = append(rs, r)
				continue
			}
			if len(rs) > 0 {
				tokens = append(tokens, string(rs))
				rs = rs[:0]
			}
			inString = false
			continue
		}

		rs = append(rs, r)
		inString = true
	}

	if inSingleQuot {
		return nil, fmt.Errorf("parse token: pos %v: not closed single quotation mark('): %v", singlePos, t[singlePos:])
	}
	if inDoubleQuot {
		return nil, fmt.Errorf("parse token: pos %v: not closed double quotation mark(\"): %v", doublePos, t[doublePos:])
	}
	if escape {
		return nil, fmt.Errorf("parse token: pos %v: unexpected end of line(\\)", escapePos)
	}
	if len(rs) > 0 {
		tokens = append(tokens, string(rs))
	}
	return tokens, nil
}
