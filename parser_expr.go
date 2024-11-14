package icl

import (
	"slices"
	"strconv"
)

func (p *Parser) parseExpression(allowed ...TokenType) Expression {
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

func (p *Parser) parseExpressionList(closeToken TokenType) (list []Expression) {
	p.nextToken()

	if p.curTokenIs(closeToken) {
		return list
	}

	// i dont like this procedure but i cant think of a better way atm
	list = append(list, p.parseExpression())
	for p.peekTokenIs(TknComma) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression())
	}

	if !p.expectPeek(closeToken) {
		return nil
	}

	return list
}

func (p *Parser) parseMapBody() map[Expression]Expression {
	body := make(map[Expression]Expression)
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

	return body
}

// parseNullLiteral parses a token as a null literal expression
func (p *Parser) parseNullLiteral() Expression {
	return &NullLiteral{Token: p.curToken}
}

// parseBooleanLiteral parses a token as a boolean literal expression
func (p *Parser) parseBooleanLiteral() Expression {
	return &BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(TknTrue)}
}

// parseStringLiteral parses a string literal token as a literal expression
func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerLiteral parses a token as an integer literal expression
func (p *Parser) parseIntegerLiteral() Expression {
	expr := &IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errorf("could not parse %q as integer", p.curToken.Literal)
		return nil
	}

	expr.Value = value

	return expr
}

// parseIdentifier parses an identifier token into an expression
func (p *Parser) parseIdentifier() Expression {
	return &Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseArrayLiteral() Expression {
	return &ArrayLiteral{
		Token:    p.curToken,
		Elements: p.parseExpressionList(TknRBracket),
	}
}

func (p *Parser) parseMapLiteral() Expression {
	return &MapLiteral{
		Token:    p.curToken,
		Elements: p.parseMapBody(),
	}
}
