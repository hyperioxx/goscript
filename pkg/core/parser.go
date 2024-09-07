package core

import (
	"fmt"
	"strconv"
)

const SYNTAX_ERROR_MSG = "syntax error on line: %d"

const (
	_ int = iota
	LOWEST
	ASSIGN_P
	IF_P
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X

	ARRAY_P
	CALL_P // foo(X)
)

var precedences = map[TokenType]int{
	EQ:         EQUALS,
	NOT_EQ:     EQUALS,
	LT:         LESSGREATER,
	GT:         LESSGREATER,
	ADD:        SUM,
	SUB:        SUM,
	MUL:        PRODUCT,
	DIV:        PRODUCT,
	LPAREN:     CALL_P,
	ASSIGN:     ASSIGN_P,
	ASSIGN_INF: ASSIGN_P,
	IF:         IF_P,
	LBRACKET:   ARRAY_P,
}

type Parser interface {
	ParseProgram() (Node, error)
	ParseNode(int) (Node, error)
}

type V1Parser struct {
	l Lexer

	curToken        Token
	peekToken       Token
	errors          []string
	Debug           bool
	lastParsedIdent string
	prefixParseFns  map[TokenType]prefixParseFn
	infixParseFns   map[TokenType]infixParseFn
}

type (
	prefixParseFn func() (Node, error)
	infixParseFn  func(Node) (Node, error)
)

func NewV1Parser(l Lexer, debug bool) Parser {
	p := &V1Parser{
		l:      l,
		errors: []string{},
		Debug:  debug,
	}

	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.infixParseFns = make(map[TokenType]infixParseFn)

	p.registerInfix(ADD, p.parseInfixNode)
	p.registerInfix(SUB, p.parseInfixNode)
	p.registerInfix(MUL, p.parseInfixNode)
	p.registerInfix(DIV, p.parseInfixNode)
	p.registerInfix(EQ, p.parseInfixNode)
	p.registerInfix(NOT_EQ, p.parseInfixNode)
	p.registerInfix(GT, p.parseInfixNode)
	p.registerInfix(LT, p.parseInfixNode)
	p.registerInfix(GT_EQ, p.parseInfixNode)
	p.registerInfix(LT_EQ, p.parseInfixNode)
	p.registerInfix(ASSIGN, p.parseInfixNode)
	p.registerInfix(ASSIGN_INF, p.parseInfixNode)
	// prefix expressions
	p.registerPrefix(INT, p.parseIntegerLiteral)
	p.registerPrefix(IDENT, p.parseIdentifier)
	p.registerPrefix(FUNC, p.parseFunctionLiteral)
	p.registerPrefix(LPAREN, p.parseLeftParen)
	p.registerPrefix(RETURN, p.parseReturnStatement)
	p.registerPrefix(IF, p.parseIfStatement)
	p.registerPrefix(FOR, p.parseForStatement)
	p.registerPrefix(STRING, p.parseStringLiteral)
	p.registerPrefix(FLOAT, p.parseFloatLiteral)
	p.registerPrefix(BOOL, p.parseBooleanLiteral)
	p.registerPrefix(TRUE, p.parseBooleanLiteral)
	p.registerPrefix(FALSE, p.parseBooleanLiteral)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *V1Parser) ParseProgram() (Node, error) {
	program := []Node{}
	for !p.curTokenIs(EOF) {
		exp, err := p.ParseNode(LOWEST)
		if err != nil {
			return nil, err
		}
		if exp != nil {
			program = append(program, exp)
		}
		p.nextToken()
	}
	return &BlockStatement{program, 0, 0}, nil
}

func (p *V1Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *V1Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *V1Parser) nextToken() {
	if p.Debug {
		fmt.Printf("DEBUG: nextToken(): Current Token: %+v, Peek Token: %+v\n", p.curToken, p.peekToken)
	}
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *V1Parser) ParseNode(precedence int) (Node, error) {
	if p.Debug {
		fmt.Printf("DEBUG: ParseNode(): Precedence: %d, Current Token: %+v\n", precedence, p.curToken)
	}
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil, nil
	}
	leftExp, err := prefix()
	if err != nil {
		return nil, err
	}
	for !p.peekTokenIs(EOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp, nil
		}

		p.nextToken()

		leftExp, err = infix(leftExp)
	}

	return leftExp, err
}

func (p *V1Parser) parseBooleanLiteral() (Node, error) {
	literal := &Boolean{
		value: p.curToken.Value == "true",
	}
	p.nextToken()
	return literal, nil
}

func (p *V1Parser) parseLeftParen() (Node, error) {
	if p.Debug {
		fmt.Println("Entering parseLeftParen")
	}

	p.nextToken()

	if p.peekTokenIs(IDENT) {
		if p.Debug {
			fmt.Println("Detected IDENT token")
		}

		ident, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}

		if p.Debug {
			fmt.Printf("Parsed IDENT: %v\n", ident.String().value)
		}

		fc := &FunctionCall{
			Name:      ident.String().value,
			Arguments: []Node{},
		}

		p.nextToken()

		for !p.curTokenIs(RPAREN) && !p.peekTokenIs(EOF) {
			p.nextToken()
			arg, err := p.ParseNode(LOWEST)
			if err != nil {
				return nil, err
			}

			fc.Arguments = append(fc.Arguments, arg)

			if p.peekTokenIs(COMMA) {
				p.nextToken()
			}

			if p.curTokenIs(RPAREN) {
				break
			}
		}

		if !p.curTokenIs(RPAREN) {
			return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
		}

		p.nextToken()

		if p.Debug {
			fmt.Println("Exiting parseLeftParen with function call")
		}

		return fc, nil
	}

	expr, err := p.ParseNode(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.Debug {
		fmt.Printf("Parsed Node: %v\n", expr)
	}

	if !p.expectPeek(RPAREN) {
		if p.Debug {
			fmt.Println("Missing RPAREN token")
		}

		return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
	}

	if p.Debug {
		fmt.Println("Exiting parseLeftParen with Node")
	}

	return expr, nil
}

func (p *V1Parser) parseStringLiteral() (Node, error) {
	return &String{value: p.curToken.Value}, nil
}

func (p *V1Parser) parseIntegerLiteral() (Node, error) {
	lit := &Integer{}
	value, err := strconv.Atoi(p.curToken.Value)
	if err != nil {
		return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
	}

	lit.value = value

	return lit, nil
}

func (p *V1Parser) parseFloatLiteral() (Node, error) {
	value, err := strconv.ParseFloat(p.curToken.Value, 64)
	if err != nil {
		return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
	}

	return &Float{value: value}, nil
}

func (p *V1Parser) parseIdentifier() (Node, error) {
	if p.Debug {
		fmt.Println("Entering parseIdentifier")
	}

	ident := &IdentifierLiteral{value: p.curToken.Value}

	if p.Debug {
		fmt.Printf("Parsed IDENT: %v\n", ident.value)
	}

	if p.peekTokenIs(LPAREN) {
		if p.Debug {
			fmt.Println("Detected LPAREN token")
		}

		p.nextToken()
		fc, err := p.parseFunctionCall(ident)

		if err != nil {
			return nil, err
		}

		if p.Debug {
			fmt.Printf("Parsed function call: %v\n", fc)
			fmt.Println("Exiting parseIdentifier with function call")
		}

		return fc, nil
	}

	return ident, nil
}

func (p *V1Parser) parseFunctionLiteral() (Node, error) {

	if p.Debug {
		fmt.Println("Entering function literal ")
	}
	fl := &FunctionLiteral{}
	p.nextToken()

	fl.Name = p.curToken.Value

	p.nextToken()

	if !p.peekTokenIs(RPAREN) && !p.curTokenIs(LBRACE) {
		for !p.curTokenIs(RPAREN) {
			p.nextToken()
			param := &IdentifierLiteral{value: p.curToken.Value}
			p.nextToken()
			fl.Arguments = append(fl.Arguments, param)
		}
		p.nextToken()

	} else { // function has 0 arguments so we skip ()
		p.nextToken()
		p.nextToken()
	}

	block, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	fl.Body = block

	return fl, nil
}

func (p *V1Parser) parseFunctionCall(function Node) (Node, error) {
	if p.Debug {
		fmt.Println("Entering parseFunctionCall")
		fmt.Printf("Function: %v\n", function)
	}

	fc := &FunctionCall{
		Name:     function.String().value,
		Function: function,
	}

	if !p.curTokenIs(RPAREN) {
		for !p.curTokenIs(RPAREN) {
			p.nextToken()
			param, err := p.ParseNode(0)
			if err != nil {
				return nil, err
			}
			p.nextToken()
			fc.Arguments = append(fc.Arguments, param)
		}
	}

	if !p.curTokenIs(RPAREN) {
		if p.Debug {
			fmt.Println("Current token is not RPAREN, returning nil")
		}
		return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
	}

	if p.Debug {
		fmt.Println("Current token is RPAREN, moving to next token")
	}

	p.nextToken()

	if p.Debug {
		fmt.Println("Exiting parseFunctionCall")
	}

	return fc, nil
}

func (p *V1Parser) parseInfixNode(left Node) (Node, error) {
	if p.Debug {
		fmt.Println("Entering parseInfixNode")
		fmt.Printf("Left Node: %v\n", left)
		fmt.Printf("Operator: %s\n", p.curToken.Value)
	}

	Node := &InfixNode{
		Left:     left,
		Operator: p.curToken.Value,
	}

	precedence := p.curPrecedence()

	if p.Debug {
		fmt.Printf("Current Precedence: %d\n", precedence)
	}

	p.nextToken()

	if p.Debug {
		fmt.Println("Parsing Right Node")
	}

	right, err := p.ParseNode(precedence)

	if err != nil {
		return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
	}

	Node.Right = right

	if p.Debug {
		fmt.Printf("Right Node: %v\n", Node.Right)
		fmt.Println("Exiting parseInfixNode")
	}
	

	return Node, nil
}

func (p *V1Parser) parsePrefixNode() (Node, error) {
	Node := &PrefixNode{
		Operator: p.curToken.Value,
	}

	p.nextToken()

	right, err := p.ParseNode(PREFIX)

	if err != nil {
		return nil, err
	}

	Node.Right = right

	return Node, nil
}

func (p *V1Parser) parseReturnStatement() (Node, error) {

	p.nextToken()

	rs := &ReturnStatement{}

	node, err := p.ParseNode(LOWEST)
	if err != nil {
		return nil, err
	}

	rs.ReturnValue = node

	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return rs, nil
}

func (p *V1Parser) parseIfStatement() (Node, error) {

	p.nextToken()

	ifExp := &IfNode{}

	condition, err := p.ParseNode(LOWEST)
	if err != nil {
		return nil, err
	}

	switch c := condition.(type) {

	case *InfixNode:
		switch c.Operator {
		case "+", "-", "/", "*", "%":
			return nil, fmt.Errorf("operator %s is non-logical on line: %d", c.Operator, c.Line)
		}

	}

	ifExp.Condition = condition

	if !p.expectPeek(LBRACE) {
		return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
	}

	block, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	ifExp.Consequence = block

	if p.peekTokenIs(ELSE) {
		p.nextToken()

		if !p.expectPeek(LBRACE) {
			return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
		}
		block, err := p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
		ifExp.Alternative = block
	}

	return ifExp, nil
}

func (p *V1Parser) parseForStatement() (Node, error) {

	if p.Debug {
		fmt.Println("Entering parseForStatement")
	}

	forExp := &ForNode{}

	components := []Node{}


	for !p.curTokenIs(LBRACE) && len(components) <= 3 {
		p.nextToken()
		node, err := p.ParseNode(LOWEST)
		if err != nil {
			return nil, err
		}
		p.nextToken()

		components = append(components, node)
	}
	if len(components) == 2 || len(components) > 3 {
		return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
	}

	if len(components) == 1 {
		forExp.Condition = components[0]
	} else if len(components) == 3 {
		forExp.Initialisation = components[0]
		forExp.Condition = components[1]
		forExp.Updater = components[2]
	}

	block, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	forExp.Body = block

	if p.Debug {
		fmt.Println("Exiting parseForStatement")
	}

	return forExp, nil
}

func (p *V1Parser) parseBlockStatement() (*BlockStatement, error) {
	block := &BlockStatement{}
	block.Statements = []Node{}

	p.nextToken()

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		stmt, err := p.ParseNode(LOWEST)
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block, nil
}

func (p *V1Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *V1Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *V1Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *V1Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *V1Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}
