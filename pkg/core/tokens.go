package core

const (
	// Special tokens
	ERROR TokenType = iota
	ILLEGAL
	EOF
	WS
	NEWLINE

	// Literals
	IDENT  // main, foo, bar, x, y, etc.
	INT    // int
	FLOAT  // 123.456
	STRING // "abc", 'abc'
	BOOL   // true
	ARRAY  // [1, 2]
	STRUCT // { a int }

	// Operators
	ADD         // +
	SUB         // -
	MUL         // *
	DIV         // /
	REM         // %
	EXP         // **
	ASSIGN      // =
	ASSIGN_INF  // :=
	LEFT_SHIFT  // <<
	RIGHT_SHIFT // >>
	XOR         //^
	ADD_ASSIGN  // +=
	SUB_ASSIGN  // -=
	INC         // ++
	DEC         // --

	// Comparators
	EQ     // ==
	NOT_EQ // !=
	GT     // >
	LT     // <
	GT_EQ  // >=
	LT_EQ  // <=
	OR

	// Delimiters
	LPAREN    // (
	RPAREN    // )
	LBRACKET  // [
	RBRACKET  // ]
	LBRACE    // {
	RBRACE    // }
	COMMA     // ,
	DOT       // .
	COLON     // :
	SEMICOLON // ;

	// Keywords
	FUNC
	VAR
	CLASS
	RETURN
	IF
	ELIF
	ELSE
	FOR
	FOREVER
	BREAK
	CONTINUE
	IMPORT
	TRUE
	FALSE
	CALL
	ASYNC
	AWAIT
)

var keywordLookup = map[string]TokenType{
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
	"import":   IMPORT,
	"true":     TRUE,
	"false":    FALSE,
	"func":     FUNC,
	"return":   RETURN,
	"int":      INT,
	"string":   STRING,
	"float":    FLOAT,
	"bool":     BOOL,
	"struct":   STRUCT,
	"async":    ASYNC,
	"await":    AWAIT,
}

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
	Error  string
}

type TokenType int

var TokenTypeStr = map[TokenType]string{
	ERROR:       "ERROR",
	ILLEGAL:     "ILLEGAL",
	EOF:         "EOF",
	WS:          "WS",
	IDENT:       "IDENT",
	INT:         "INT",
	FLOAT:       "FLOAT",
	STRING:      "STRING",
	ARRAY:       "ARRAY",
	BOOL:        "BOOL",
	ADD:         "+",
	SUB:         "-",
	MUL:         "*",
	DIV:         "/",
	REM:         "%",
	EXP:         "^",
	ASSIGN:      "=",
	ASSIGN_INF:  ":=",
	LEFT_SHIFT:  "<<",
	RIGHT_SHIFT: ">>",
	ADD_ASSIGN:  "+=",
	SUB_ASSIGN:  "-=",
	INC:         "++",
	DEC:         "--",
	EQ:          "==",
	NOT_EQ:      "!=",
	GT:          ">",
	LT:          "<",
	GT_EQ:       ">=",
	LT_EQ:       "<=",
	LPAREN:      "(",
	RPAREN:      ")",
	LBRACKET:    "[",
	RBRACKET:    "]",
	LBRACE:      "{",
	RBRACE:      "}",
	COMMA:       ",",
	DOT:         ".",
	COLON:       ":",
	SEMICOLON:   ";",
	FUNC:        "FUNC",
	RETURN:      "RETURN",
	IF:          "IF",
	ELSE:        "ELSE",
	FOR:         "FOR",
	BREAK:       "BREAK",
	CONTINUE:    "CONTINUE",
	IMPORT:      "IMPORT",
	TRUE:        "TRUE",
	FALSE:       "FALSE",
	NEWLINE:     "NEWLINE",
	CALL:        "CALL",
}
