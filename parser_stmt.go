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
