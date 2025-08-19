// Package lexer
package lexer

import (
	"monkey/token"
	"strings"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) Reset() {
	l.position = 0
	l.readPosition = 0
	l.readChar()
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	// Operators
	case '=':
		literal := l.readTwoCharToken('=')
		if len(literal) > 1 {
			tok.Literal = literal
			tok.Type = token.EQ
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		literal := l.readTwoCharToken('+')
		if len(literal) > 1 {
			tok.Literal = literal
			tok.Type = token.INCREMENT
		} else {
			tok = newToken(token.PLUS, l.ch)
		}
	case '-':
		literal := l.readTwoCharToken('-')
		if len(literal) > 1 {
			tok.Literal = literal
			tok.Type = token.DECREMENT
		} else {
			tok = newToken(token.MINUS, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '!':
		literal := l.readTwoCharToken('=')
		if len(literal) > 1 {
			tok.Literal = literal
			tok.Type = token.NOTEQ
		} else {
			tok = newToken(token.BANG, l.ch)
		}

	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)

		// Delimiters
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)

	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch, 10) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	base := 10
	if l.ch == '0' && isBasePrefix(l.peekChar()) {
		l.readChar()
		prefix := l.ch
		switch prefix {
		case 'b':
			base = 2
		case 'o':
			base = 8
		case 'x':
			base = 16
		}
		l.readChar()
	}
	for isDigit(l.ch, base) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	var out strings.Builder
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			switch l.peekChar() {
			case 't':
				out.WriteRune('\t')
			case 'n':
				out.WriteRune('\n')
			case 'r':
				out.WriteRune('\r')
			case '"':
				out.WriteRune('"')
			default:
				out.WriteRune('\\')
				out.WriteByte(l.peekChar())
			}
			l.readChar()
		} else {
			out.WriteByte(l.ch)
		}
	}
	return out.String()
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func (l *Lexer) readTwoCharToken(targetChar byte) string {
	result := string(l.ch)
	if l.peekChar() == targetChar {
		l.readChar()
		result += string(l.ch)
	}

	return result
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '$'
}

func isDigit(ch byte, base int) bool {
	switch base {
	case 2:
		return '0' <= ch && ch <= '1'
	case 8:
		return '0' <= ch && ch <= '7'
	case 10:
		return '0' <= ch && ch <= '9'
	case 16:
		return '0' <= ch && ch <= '9' || 'A' <= ch && ch <= 'F' || 'a' <= ch && ch <= 'f'
	default:
		return '0' <= ch && ch <= '9'
	}
}

func isBasePrefix(ch byte) bool {
	return ch == 'b' || ch == 'x' || ch == 'o'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
