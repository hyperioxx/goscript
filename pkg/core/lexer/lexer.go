package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Lexer interface {
	NextToken() Token
}

type V1Lexer struct {
	input         string
	position      int
	readPosition  int
	line          int
	column        int
	ch            rune
	currentIndent int
	indentStack   []int
	Debug         bool
}

func NewV1Lexer(input string) Lexer {
	l := &V1Lexer{
		input:       input + "\n",
		indentStack: []int{0},
		line:        1,
		column:      0,
	}
	l.readChar()
	return l
}

func (l *V1Lexer) readChar() {
	if l.ch == '\n' {
		l.line++
		l.column = -1
	}

	if l.readPosition >= len(l.input) {
		l.ch = -1
	} else {
		l.ch = []rune(string(l.input[l.readPosition]))[0]
		l.position = l.readPosition
		l.readPosition++
		l.column++
	}
}

func (l *V1Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '+':
		if l.peekChar() == '+' {
			tok = newToken(INC, "++")
			l.readChar()
		} else if l.peekChar() == '=' {
			tok = newToken(ADD_ASSIGN, "+=")
			l.readChar()
		} else {
			tok = newToken(ADD, "+")
		}
	case '-':
		if l.peekChar() == '-' {
			tok = newToken(DEC, "--")
			l.readChar()
		} else if isDigit(l.peekChar()) {
			tok.Type = INT
			tok.Value = "-"
			l.readChar() // Skip the "-"
			l.skipWhitespace()
			if isDigit(l.ch) {
				tok.Value += l.readNumber()
				return tok
			}
		} else if l.peekChar() == '=' {
			tok = newToken(SUB_ASSIGN, "-=")
			l.readChar()
		} else {
			tok = newToken(SUB, "-")
		}
	case '*':
		if l.peekChar() == '*' {
			tok = newToken(EXP, "**")
			l.readChar()
		} else {
			tok = newToken(MUL, "*")
		}
	case '/':
		if l.peekChar() == '/' {
			// this case catches comments
			// thinking about catching this and bringing it into the parser think golang tags
			l.skipComment()
			return l.NextToken()
		} else {
			tok = newToken(DIV, "/")
		}
	case '%':
		tok = newToken(REM, "%")
	case '|':
		if l.peekChar() == '=' {
			tok = newToken(OR, "or")
		}
	case '^':
		if l.peekChar() == '=' {
			tok = newToken(XOR, "^")
		}
	case '<':
		if l.peekChar() == '=' {
			tok = newToken(LT_EQ, "<=")
			l.readChar()
		} else if l.peekChar() == '<' {
			tok = newToken(LEFT_SHIFT, "<<")
			l.readChar()
		} else {
			tok = newToken(LT, "<")
		}
	case '>':
		if l.peekChar() == '=' {
			tok = newToken(GT_EQ, ">=")
			l.readChar()
		} else if l.peekChar() == '>' {
			tok = newToken(RIGHT_SHIFT, ">>")
			l.readChar()
		} else {
			tok = newToken(GT, ">")
		}
	case '=':
		if l.peekChar() == '=' {
			tok = newToken(EQ, "==")
			l.readChar()
		} else {
			tok = newToken(ASSIGN, "=")
		}
	case '!':
		if l.peekChar() == '=' {
			tok = newToken(NOT_EQ, "!=")
			l.readChar()
		}
	case ';':
		tok = newToken(SEMICOLON, ";")
	case '\n':
		l.readChar()
		if l.ch == -1 {
			return newToken(EOF, "EOF")
		}

		return l.NextToken()
	case '[':
		tok = newToken(LBRACKET, "[")
	case ']':
		tok = newToken(RBRACKET, "]")
	case '(':
		tok = newToken(LPAREN, "(")
	case ')':
		tok = newToken(RPAREN, ")")
	case ',':
		tok = newToken(COMMA, ",")
	case ':':
		if l.peekChar() == '=' {
            tok = newToken(ASSIGN_INF, ":=")
            l.readChar()
        } else {
            tok = newToken(COLON, ":")
        }
	case '{':
		tok = newToken(LBRACE, "{")
	case '}':
		tok = newToken(RBRACE, "}")
	case '"':
		tok.Type = STRING
		tok.Value = l.readString()
	case '#':
		l.skipComment()
		return l.NextToken()
	case '.':
		tok = newToken(DOT, ".")
	default:
		if isLetter(l.ch) {
			tok.Value = l.readIdentifier()
			if keywordType, isKeyword := keywordLookup[tok.Value]; isKeyword {
				tok.Type = keywordType
			} else {
				tok.Type = IDENT
			}
			return tok
		} else if isDigit(l.ch) {
			tok.Value = l.readNumber()
			if strings.Contains(tok.Value, ".") {
				tok.Type = FLOAT
			} else {
				tok.Type = INT
			}
			return tok
		} else if l.ch == -1 {
			tok.Type = EOF
			tok.Value = "EOF"
			return tok
		} else {
			tok = Token{
				Type:   ERROR,
				Value:  string(l.ch),
				Line:   l.line,
				Column: l.column,
				Error:  fmt.Sprintf("Unexpected character: %q", l.ch),
			}
		}
	}

	// this is a kinda rough way to fast forward on remaining whitespace
	for l.ch == ' ' {
		l.readChar()
	}

	if l.Debug {
		fmt.Printf("NextToken: %v\n", tok)
	}

	l.readChar()
	return tok
}

func (l *V1Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == -1 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *V1Lexer) skipComment() {
	for l.ch != '\n' && l.ch != -1 {
		l.readChar()
	}
}

func (l *V1Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return -1
	} else {
		ch, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
		return ch
	}
}

func newToken(tokenType TokenType, ch string) Token {
	return Token{Type: tokenType, Value: ch}
}

func (l *V1Lexer) skipWhitespace() {

	for l.ch == ' ' || l.ch == '\t' {
		l.readChar()
	}

}

func (l *V1Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *V1Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
