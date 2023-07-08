package parser


import ( "necronet.info/interpreter/ast"
       "necronet.info/interpreter/lexer"
       "necronet.info/interpreter/token"
)


type Parser struct { 
	l *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser { 
	p := &Parser{l: l}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p 
}

func (p *Parser) nextToken() { 
	p.curToken = p.peekToken 
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program { 
	// program = newProgramASTNode()
	// advanceTokens()
	// for (currentToken() != EOF_TOKEN) {
	// 	statement = null
		
	// 	if (currentToken() == LET_TOKEN) { 
	// 		statement = parseLetStatement()
	// 	} else if (currentToken() == RETURN_TOKEN) { 
	// 		statement = parseReturnStatement()
	// 	} else if (currentToken() == IF_TOKEN) { 
	// 		statement = parseIfStatement()
	// 	}

	// 	if (statement != null) { 
	// 		program.Statements.push(statement)
	// 	}
	// 	advanceTokens()
	// }

	return nil
}

