package token

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

// Defining literals

const(
	ILLEGAL	= "ILLEGAL"
	EOF = "EOF"

	IDENT = "IDENT"
	INT = "INT"

	ASSIGN = "="
	PLUS = "+"
	MINUS = "-"

	EQ     = "=="
	NOT_EQ = "!="

	BANG     = "!"
        ASTERISK = "*"
        SLASH    = "/"
        LT = "<"
        GT = ">"

	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"

	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	LET = "LET"
    STRING = "STRING"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType {
	"fn":FUNCTION,
	"let": LET,
	"true": TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

