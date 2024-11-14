package icl

import "fmt"

type prefixParser func() Expression
type infixParser func(Expression) Expression

type Parser struct {
	lex *Lexer

	errors    []error
	curToken  Token
	peekToken Token

	prefixParsers map[TokenType]prefixParser
	infixParsers  map[TokenType]infixParser
}

// New creates a new parser for the provided lexer
func NewParser(lex *Lexer) *Parser {
	p := &Parser{
		lex:           lex,
		prefixParsers: make(map[TokenType]prefixParser),
		infixParsers:  make(map[TokenType]infixParser),
	}

	p.registerPrefixParser(TknIdent, p.parseIdentifier)
	p.registerPrefixParser(TknInt, p.parseIntegerLiteral)
	p.registerPrefixParser(TknNull, p.parseNullLiteral)
	p.registerPrefixParser(TknTrue, p.parseBooleanLiteral)
	p.registerPrefixParser(TknFalse, p.parseBooleanLiteral)
	p.registerPrefixParser(TknString, p.parseStringLiteral)
	p.registerPrefixParser(TknLBracket, p.parseArrayLiteral)

	// p.registerInfixParser(token.Equal, p.parseInfixExpression)

	// populate cur/next token fields
	p.nextToken()
	p.nextToken()

	return p
}

// ParseProgram parses the tokens in the lexer into an AST
func (p *Parser) Parse() *Ast {
	program := &Ast{}

	for p.curToken.Type != TknEof {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

// Errors returns a slice of errors generated by the parser
func (p *Parser) Errors() []error {
	return p.errors
}

// registerPrefixParser registers a prefix parser funcion for the token type
func (p *Parser) registerPrefixParser(tknType TokenType, parser prefixParser) {
	p.prefixParsers[tknType] = parser
}

// registerInfixParser registers an infix parser funtion for the token type
func (p *Parser) registerInfixParser(tknType TokenType, parser infixParser) {
	p.infixParsers[tknType] = parser
}

func (p *Parser) errorf(format string, args ...any) {
	lineKey := fmt.Sprintf(" -- [line(%d) pos(%d)]", p.peekToken.Line, p.peekToken.Pos)
	p.errors = append(p.errors, fmt.Errorf(format+lineKey, args...))
}

// nextToken advances the lexer to the next token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

// parseStatement parses the next statement in the lexers token stream
func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	// TODO
	case TknIdent:
		return p.parseAssignStatement()
	// case TknReturn:
	// 	return p.parseReturnStatement()
	case TknComment:
		return nil
	default:
		return nil
		// TODO
		//return p.parseExpressionStatement()
	}
}

// peekTokenIs checks if the next token in the stream is of the provided type
func (p *Parser) peekTokenIs(tknType TokenType) bool {
	return p.peekToken.Type == tknType
}

// curTokenIs checks if the current token is of the given type
func (p *Parser) curTokenIs(tknType TokenType) bool {
	return p.curToken.Type == tknType
}

// peekError checks that the next token is of the given type
// if not then an error will be generated
func (p *Parser) peekError(tknType TokenType) (err error) {
	if !p.peekTokenIs(tknType) {
		p.errorf("Unexpected token type: expected(%s) found(%s)", tknType, p.peekToken.Type)
	}

	return
}

// expectPeek runs the peekError method and advances the token streabm if no erro is found
func (p *Parser) expectPeek(tknType TokenType) bool {
	if p.peekError(tknType) != nil {
		return false
	}

	p.nextToken()
	return true
}
