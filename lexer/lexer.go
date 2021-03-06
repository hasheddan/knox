// Based on The Interpreter Book.

package lexer

import (
	"fmt"
	"knox/token"
)

// Lexer object.
type Lexer struct {
	position     int    //current character position
	readPosition int    //next character position
	ch           rune   //current character
	characters   []rune //rune slice of input string
}

// New a Lexer instance from string input.
func New(input string) *Lexer {
	l := &Lexer{characters: []rune(input)}
	l.readChar()
	return l
}

// read one forward character
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.characters) {
		l.ch = rune(0)
	} else {
		l.ch = l.characters[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// NextToken to read next token, skipping the white space.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	l.checkComments()
	switch l.ch {
	case rune('='):
		tok = newToken(token.ASSIGN, l.ch)
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case rune(';'):
		tok = newToken(token.SEMICOLON, l.ch)
	case rune('('):
		tok = newToken(token.LPAREN, l.ch)
	case rune(')'):
		tok = newToken(token.RPAREN, l.ch)
	case rune(','):
		tok = newToken(token.COMMA, l.ch)
	case rune('+'):
		tok = newToken(token.PLUS, l.ch)
	case rune('{'):
		tok = newToken(token.LBRACE, l.ch)
	case rune('}'):
		tok = newToken(token.RBRACE, l.ch)
	case rune('-'):
		tok = newToken(token.MINUS, l.ch)
	case rune('/'):
		tok = newToken(token.SLASH, l.ch)
	case rune('*'):
		tok = newToken(token.ASTERISK, l.ch)
	case rune('%'):
		tok = newToken(token.PERCENT, l.ch)
	case rune('^'):
		tok = newToken(token.CARET, l.ch)
	case rune('.'):
		tok = newToken(token.DOT, l.ch)
	case rune('<'):
		tok = newToken(token.LT, l.ch)
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case rune('>'):
		tok = newToken(token.GT, l.ch)
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case rune('&'):
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch)}
		}
	case rune('|'):
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch)}
		}
	case rune('!'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOTEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case rune('"'):
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case rune('['):
		tok = newToken(token.LBRACKET, l.ch)
	case rune(']'):
		tok = newToken(token.RBRACKET, l.ch)
	case rune(':'):
		tok = newToken(token.COLON, l.ch)
	//case rune('.'):
	//return l.readFloat()
	case rune(0):
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isDigit(l.ch) {
			return l.readDecimal()
		} else {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		}
	}
	l.readChar()
	return tok
}

// return new token
func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// read Identifier
func (l *Lexer) readIdentifier() string {
	// Already checked first ch for isDigit.
	position := l.position
	for isIdentifier(l.ch) {
		l.readChar()
	}
	return string(l.characters[position:l.position])
}

// skip white space
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) checkComments() {
	for l.ch == rune('/') && l.peekChar() == '/' {
		l.skipComment()
		l.skipWhitespace()
	}
}

// Ignore comments.
func (l *Lexer) skipComment() {
	for l.ch != rune('\n') {
		l.readChar()
	}
}

// read number
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) || isUnderscore(l.ch) {
		l.readChar()
	}
	return string(l.characters[position:l.position])
}

// read until white space
func (l *Lexer) readUntilWhitespace() string {
	position := l.position
	for !isWhitespace(l.ch) {
		l.readChar()
	}
	return string(l.characters[position:l.position])
}

// read decimal
func (l *Lexer) readDecimal() token.Token {
	integer := l.readNumber()
	if l.ch == rune('.') {
		l.readChar()
		fraction := l.readNumber()
		if isEmpty(l.ch) || isWhitespace(l.ch) || IsOperator(l.ch) || isComparison(l.ch) || isCompound(l.ch) || isBracket(l.ch) || isBrace(l.ch) || isParen(l.ch) {
			return token.Token{Type: token.FLOAT, Literal: integer + "." + fraction}
		} else {
			illegalPart := l.readUntilWhitespace()
			return token.Token{Type: token.ILLEGAL, Literal: integer + "." + fraction + illegalPart}
		}
	} else if isEmpty(l.ch) || isWhitespace(l.ch) || IsOperator(l.ch) || isComparison(l.ch) || isCompound(l.ch) || isBracket(l.ch) || isBrace(l.ch) || isParen(l.ch) {
		return token.Token{Type: token.INT, Literal: integer}
	} else {
		illegalPart := l.readUntilWhitespace()
		return token.Token{Type: token.ILLEGAL, Literal: integer + illegalPart}
	}
}

// read string
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' {
			break
		}
		if l.ch == rune(0) {
			fmt.Println("End of string literal not found.")
			panic("Aborted.")
		}
	}
	return string(l.characters[position:l.position])
}

// peek character
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.characters) {
		return rune(0)
	} else {
		return l.characters[l.readPosition]
	}
}

// determinate ch is identifier or not
func isIdentifier(ch rune) bool {
	return !isWhitespace(ch) && !isBrace(ch) && !IsOperator(ch) && !isComparison(ch) && !isCompound(ch) && !isBrace(ch) && !isParen(ch) && !isBracket(ch) && !isEmpty(ch) && !isDot(ch)
}

func isDot(ch rune) bool {
	return ch == rune('.')
}

// is white space
func isWhitespace(ch rune) bool {
	return ch == rune(' ') || ch == rune('\t') || ch == rune('\n') || ch == rune('\r')
}

// IsOperator returns
func IsOperator(ch rune) bool {
	return ch == rune('+') || ch == rune('-') || ch == rune('/') || ch == rune('*') || ch == rune('%') || ch == rune('^')
}

// is comparison
func isComparison(ch rune) bool {
	return ch == rune('=') || ch == rune('!') || ch == rune('>') || ch == rune('<')
}

// is compound
func isCompound(ch rune) bool {
	return ch == rune(',') || ch == rune(':') || ch == rune('"') || ch == rune(';')
}

// is brace
func isBrace(ch rune) bool {
	return ch == rune('{') || ch == rune('}')
}

// is bracket
func isBracket(ch rune) bool {
	return ch == rune('[') || ch == rune(']')
}

// is parenthesis
func isParen(ch rune) bool {
	return ch == rune('(') || ch == rune(')')
}

// is empty
func isEmpty(ch rune) bool {
	return rune(0) == ch
}

// is Digit
func isDigit(ch rune) bool {
	return rune('0') <= ch && ch <= rune('9')
}

func isUnderscore(ch rune) bool {
	return ch == rune('_')
}
