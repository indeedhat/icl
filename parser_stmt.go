package icl

// parseAssignStatement parses a let statement from the lexers token stream
func (p *Parser) parseAssignStatement() *AssignStatement {
	stmt := &AssignStatement{Token: p.curToken}

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
func (p *Parser) parseBlockStatement() Statement {
	expr := &BlockStatement{Token: p.curToken}

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

	expr.Body = p.parseBlockBodyStatement()

	return expr
}

func (p *Parser) parseBlockBodyStatement() *BlockBodyStatement {
	block := &BlockBodyStatement{Token: p.curToken}

	// advance cursor bast { token
	p.nextToken()

	// loop until either a } or EOF token is found
	for !p.curTokenIs(TknRBrace) && !p.curTokenIs(TknEof) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		// p.parseStatement() leaves the cursor on the final token of the statement so we need to advance
		// the cursor before the next parse
		p.nextToken()
	}

	return block
}
