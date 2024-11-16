package icl

import (
	"slices"
	"strconv"
)

func (p *Parser) parseExpression(allowed ...TokenType) Node {
	if len(allowed) > 0 && !slices.Contains(allowed, p.curToken.Type) {
		p.errorf("token type %s is not allowed here", p.curToken.Type)
		return nil
	}

	prefix := p.prefixParsers[p.curToken.Type]
	if prefix == nil {
		p.errorf("no prefix parser found for %s", p.curToken.Type)
		return nil
	}

	return prefix()
}

func (p *Parser) parseListEntries(closeToken TokenType) (list []Node) {
	p.nextToken()

	if p.curTokenIs(closeToken) {
		return list
	}

	// i dont like this procedure but i cant think of a better way atm
	list = append(list, p.parseNode())
	for p.peekTokenIs(TknComma) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseNode())
	}

	if !p.expectPeek(closeToken) {
		return nil
	}

	return list
}

func (p *Parser) parseMapBody() map[Node]Node {
	body := make(map[Node]Node)
	closeToken := TknRBrace

	if p.curTokenIs(TknLBrace) {
		p.nextToken()
	}

	if p.curTokenIs(closeToken) {
		return body
	}

	if p.peekToken.Type != TknColon {
		p.errorf("no prefix parser found for %s", p.curToken.Type)
		return nil
	}

	for !p.peekTokenIs(closeToken) {
		for p.curTokenIs(TknComment) {
			p.nextToken()
		}

		if p.curTokenIs(closeToken) {
			break
		}

		key := p.parseExpression(TknIdent, TknString)
		if !p.expectPeek(TknColon) {
			return body
		}

		p.nextToken()
		value := p.parseExpression()

		body[key] = value

		if p.peekToken.Type != closeToken {
			if !p.expectPeek(TknComma) {
				p.errorf("expected : or }")
			}
			p.nextToken()
		}
	}

	// advance past closing }
	p.nextToken()

	return body
}

// parseNullNode parses a token as a null literal expression
func (p *Parser) parseNullNode() Node {
	return &NullNode{Token: p.curToken}
}

// parseBooleanNode parses a token as a boolean literal expression
func (p *Parser) parseBooleanNode() Node {
	return &BooleanNode{Token: p.curToken, Value: p.curTokenIs(TknTrue)}
}

// parseStringNode parses a string literal token as a literal expression
func (p *Parser) parseStringNode() Node {
	return &StringNode{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerNode parses a token as an integer literal expression
func (p *Parser) parseIntegerNode() Node {
	expr := &IntegerNode{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errorf("could not parse %q as integer", p.curToken.Literal)
		return nil
	}

	expr.Value = value

	return expr
}

// parseIdentifier parses an identifier token into an expression
func (p *Parser) parseIdentifier() Node {
	if p.curToken.Literal == "env" && p.peekTokenIs(TknLParen) {
		n := EnvarNode{
			Token: p.curToken,
		}

		// advance past (
		p.nextToken()

		if !p.expectPeek(TknIdent) {
			return nil
		}

		n.Identifier = p.parseIdentifier().(*Identifier)

		if !p.expectPeek(TknRParen) {
			return nil
		}

		return &n
	}

	return &Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseArrayNode() Node {
	return &ArrayNode{
		Token:    p.curToken,
		Elements: p.parseListEntries(TknRBracket),
	}
}

func (p *Parser) parseMapNode() Node {
	return &MapNode{
		Token:    p.curToken,
		Elements: p.parseMapBody(),
	}
}

// parseAssignNode parses a let statement from the lexers token stream
func (p *Parser) parseAssignNode() *AssignNode {
	stmt := &AssignNode{Token: p.curToken}

	if p.peekToken.Type != TknAssign {
		return nil
	}

	stmt.Name = &Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	p.nextToken()
	p.nextToken()
	stmt.Value = p.parseExpression()

	return stmt
}

// function
func (p *Parser) parseBlockNode() Node {
	expr := &BlockNode{Token: p.curToken}

	// advance past ident
	p.nextToken()

	// expr.Parameters =
	for p.curTokenIs(TknString) || p.curTokenIs(TknIdent) {
		expr.Parameters = append(expr.Parameters, p.curToken)
		p.nextToken()
	}

	if !p.expectPeek(TknLParen) {
		return nil
	}

	expr.Body = p.parseBlockBodyNode()

	return expr
}

func (p *Parser) parseBlockBodyNode() *BlockBodyNode {
	block := &BlockBodyNode{Token: p.curToken}

	// advance cursor bast { token
	p.nextToken()

	// loop until either a } or EOF token is found
	for !p.curTokenIs(TknRBrace) && !p.curTokenIs(TknEof) {
		stmt := p.parseNode()
		if stmt != nil {
			block.Nodes = append(block.Nodes, stmt)
		}

		// p.parseNode() leaves the cursor on the final token of the statement so we need to advance
		// the cursor before the next parse
		p.nextToken()
	}

	return block
}
