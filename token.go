package icl

type TokenType string

const (
	// Special
	TknIllegal TokenType = "ILLEGAL"
	TknEof     TokenType = "EOF"

	// Identifiers & literals
	TknIdent  TokenType = "IDENT"
	TknInt    TokenType = "INT"
	TknString TokenType = "STRING"

	// Operators
	TknAssign      TokenType = "="
	TknHash        TokenType = "#"
	TknLessThan    TokenType = "<"
	TknGreaterThan TokenType = ">"

	// Delimiters
	TknComma TokenType = ","
	TknColon TokenType = ":"

	TknLParen   TokenType = "("
	TknRParen   TokenType = ")"
	TknLBrace   TokenType = "{"
	TknRBrace   TokenType = "}"
	TknLBracket TokenType = "["
	TknRBracket TokenType = "]"

	// Keywords
	TknTrue    TokenType = "TRUE"
	TknFalse   TokenType = "FALSE"
	TknNull    TokenType = "NULL"
	TknComment TokenType = "COMMENT"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Pos     int
}

var keywords = map[string]TokenType{
	"true":  TknTrue,
	"false": TknFalse,
	"null":  TknNull,
}

// LookupIdent looks up the identifer in the map of keywords and returns the appropriate
// token type for the string
func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return TknIdent
}
