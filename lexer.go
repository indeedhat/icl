package icl

import (
	"bytes"
)

type Lexer struct {
	input string
	// current position of input (position of char)
	pos int
	// current reading pos (char + 1)
	readPos int
	// current char under exam
	char byte

	// line of the input
	line int
	// cursor pos on current line
	linePos int
}

// newLexer creates a new Lexer instance with the provided input string
func newLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}

func (l *Lexer) NextToken() Token {
	defer l.readChar()
	l.consumeWhitespace()

	switch l.char {
	case ',':
		return l.token(TknComma, string(l.char))
	case '(':
		return l.token(TknLParen, string(l.char))
	case ')':
		return l.token(TknRParen, string(l.char))
	case '{':
		return l.token(TknLBrace, string(l.char))
	case '}':
		return l.token(TknRBrace, string(l.char))
	case '[':
		return l.token(TknLBracket, string(l.char))
	case ']':
		return l.token(TknRBracket, string(l.char))
	case '=':
		return l.token(TknAssign, string(l.char))
	case '#':
		return l.token(TknComment, l.readLineComment())
	case '"':
		str := l.readStringLiteral()
		if str == nil {
			return l.token(TknIllegal, string(l.char))
		}
		return l.token(TknString, *str)
	case 0:
		return Token{Type: TknEof}
	default:
		if isIdentChar(l.char) {
			ident := l.readIdentifier()
			return l.token(lookupIdent(ident), ident)
		}
		if isDigit(l.char) {
			return l.token(TknInt, l.readNumber())
		}

		return l.token(TknIllegal, string(l.char))
	}
}

// New initializes a new token
func (l *Lexer) token(tokenType TokenType, char string) Token {
	return Token{
		Type:    tokenType,
		Literal: char,
		Line:    l.line,
		Pos:     l.linePos - len(char),
	}
}

// readChar reads the next char in the input string
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPos]
	}

	l.pos = l.readPos
	l.readPos++
	l.linePos++
}

// readChar reads the next char in the input string
func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}

	return l.input[l.readPos]
}

// multiReadChar reads multiple characters as a string
func (l *Lexer) multiReadChar(n int) string {
	pos := l.pos
	for i := 0; i < n-1; i++ {
		l.readChar()
	}

	return l.input[pos:l.readPos]
}

// consumeWhitespace keeps reading characters until the current char is not a valid whitespace character
func (l *Lexer) consumeWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		if l.char == '\n' {
			l.line++
			l.linePos = 0
		}
		l.readChar()
	}
}

// readIdentifier reads a string of consecutive valid identifier characters
func (l *Lexer) readIdentifier() string {
	pos := l.pos

	for isIdentChar(l.peekChar(), true) {
		l.readChar()
	}

	return l.input[pos:l.readPos]
}

// readStringLiteral reads a string literal
// this will only accept double quoted strings and can have double quotes escaped with a \
func (l *Lexer) readStringLiteral() *string {
	var buf bytes.Buffer

	for {
		char := l.input[l.pos]
		peekChar := l.peekChar()
		// if we dont find a closing " then its an invalid string
		if peekChar == 0 {
			return nil
		}

		if char != '\\' || peekChar != '"' {
			buf.WriteString(string(char))
		}

		l.readChar()
		if peekChar == '"' && char != '\\' {
			break
		}
	}

	// skip final "
	if l.peekChar() == '"' {
		l.readChar()
	}

	// dont include the " on each side in the strings value
	str := buf.String()[1:]
	return &str
}

func (l *Lexer) readLineComment() string {
	var buf bytes.Buffer

	for {
		char := l.input[l.pos]
		peekChar := l.peekChar()

		if peekChar == 0 {
			break
		} else if char == '\r' && peekChar != '\n' {
			break
		} else if char == '\n' {
			break
		}
		buf.WriteString(string(char))

		l.readChar()
	}

	return buf.String()
}

func (l *Lexer) readBlockComment() string {
	var buf bytes.Buffer

	for {
		char := l.input[l.pos]
		peekChar := l.peekChar()
		if peekChar == 0 {
			break
		}

		buf.WriteString(string(char))

		if char == '*' && peekChar == '/' {
			l.readChar()
			buf.WriteString("/")
			break
		}

		l.readChar()
	}

	return buf.String()
}

// readNumber reads a numeric value from the input string
func (l *Lexer) readNumber() string {
	pos := l.pos

	for isDigit(l.peekChar()) {
		l.readChar()
	}

	return l.input[pos:l.readPos]
}

// isIdentChar checks if the provided byte is a valid character for an identifier
func isIdentChar(char byte, subsequent ...bool) bool {
	// numbers are allowed in idents but not as the first character
	if len(subsequent) > 0 && subsequent[0] && char >= '0' && char <= '9' {
		return true
	}

	return char >= 'a' && char <= 'z' ||
		char >= 'A' && char <= 'Z' ||
		char == '_'
}

// isDigit checks if the byte is a numeric character
func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}
