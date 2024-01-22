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
	p.registerInfix(ASSIGN, p.parseAssignNode)
	p.registerInfix(ASSIGN_INF, p.parseAssignNode)
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
	p.registerPrefix(LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(TRUE, p.parseBooleanLiteral)
	p.registerPrefix(FALSE, p.parseBooleanLiteral)
	p.registerPrefix(IMPORT, p.parseImport)

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

func (p *V1Parser) parseImport() (Node, error) {
	if p.Debug {
		fmt.Println("Entering parseImport")
	}

	paths := []string{}

	if !p.peekTokenIs(LPAREN) {
		if !p.expectPeek(STRING) {
			return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
		}
		paths = append(paths, p.curToken.Value)
	} else {
		p.nextToken()

		for !p.peekTokenIs(RPAREN) {
			if !p.expectPeek(STRING) {
				return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
			}
			paths = append(paths, p.curToken.Value)
			p.nextToken()
		}

		p.nextToken()
	}

	moduleListNode := &ModuleListNode{
		Modules: []Node{},
	}

	for _, path := range paths {
		moduleListNode.Modules = append(moduleListNode.Modules, NewStringLiteral(path, p.curToken.Line, p.curToken.Column))
	}

	if p.Debug {
		fmt.Println("Exiting parseImport")
	}

	return moduleListNode, nil
}

// parseBooleanLiteral parses the 'true' or 'false' token and creates a BooleanLiteral Node
func (p *V1Parser) parseBooleanLiteral() (Node, error) {
	literal := &BooleanLiteral{
		value:  p.curToken.Type == TRUE,
		Line:   p.curToken.Line,
		Column: p.curToken.Column,
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
			fmt.Printf("Parsed IDENT: %v\n", ident.Value())
		}

		fc := &FunctionCall{
			Name:      ident.Value().(string),
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
	return &StringLiteral{value: p.curToken.Value}, nil
}

func (p *V1Parser) parseIntegerLiteral() (Node, error) {
	lit := &IntegerLiteral{}
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

	return &FloatLiteral{value: value}, nil
}

func (p *V1Parser) parseHashLiteral() (Node, error) {
	hash := &HashLiteral{Pairs: make(map[Node]Node)}

	for !p.peekTokenIs(RBRACE) {
		p.nextToken()
		key, err := p.ParseNode(LOWEST)
		if err != nil {
			return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
		}

		if !p.expectPeek(COLON) {
			return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
		}

		p.nextToken()
		value, err := p.ParseNode(LOWEST)
		if err != nil {
			return nil, err
		}

		hash.Pairs[key] = value

		if !p.peekTokenIs(RBRACE) && !p.expectPeek(COMMA) {
			return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
		}
	}

	if !p.expectPeek(RBRACE) {
		return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
	}

	return hash, nil
}

func (p *V1Parser) parseIdentifier() (Node, error) {
	if p.Debug {
		fmt.Println("Entering parseIdentifier")
	}

	ident := &IdentifierLiteral{value: p.curToken.Value}
	p.lastParsedIdent = p.curToken.Value

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

	// if p.peekTokenIs(DOT) {
	// 	if p.Debug {
	// 		fmt.Println("Detected DOT token")
	// 	}

	// 	// Note: Here, instead of nextToken(), we're using a new function parseDotNotation()
	// 	// which will handle the parsing of everything following the DOT.
	// 	return p.parseDotNotation(ident)
	// }

	if p.peekTokenIs(INC) {
		p.nextToken()
		return &IncrementNode{Operand: ident}, nil
	}
	if p.peekTokenIs(DEC) {
		p.nextToken()
		return &DecrementNode{Operand: ident}, nil
	}

	if p.Debug {
		fmt.Println("Exiting parseIdentifier with ident")
	}

	return ident, nil
}


func (p *V1Parser) parseFunctionCall(function Node) (Node, error) {
	if p.Debug {
		fmt.Println("Entering parseFunctionCall")
		fmt.Printf("Function: %v\n", function)
	}

	fc := &FunctionCall{
		Name:     p.lastParsedIdent,
		Function: function,
	}

	p.nextToken()

	if !p.curTokenIs(RPAREN) {
		node, err := p.ParseNode(LOWEST)
		if err != nil {
			return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
		}
		fc.Arguments = append(fc.Arguments, node)

		for p.peekTokenIs(COMMA) {
			p.nextToken() // This will consume the comma
			p.nextToken() // This will move us to the next argument
			node, err := p.ParseNode(LOWEST)
			if err != nil {
				return nil, err
			}
			fc.Arguments = append(fc.Arguments, node)
		}

		if !p.expectPeek(RPAREN) {
			return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
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

func (p *V1Parser) parseEqualityInequalityNode(left Node) (Node, error) {
	Node := &InfixNode{
		Left:     left,
		Operator: p.curToken.Value,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	right, err := p.ParseNode(precedence)
	if err != nil {
		return nil, err
	}
	Node.Right = right

	return Node, nil
}

func (p *V1Parser) parseAssignNode(left Node) (Node, error) {
	Node := &InfixNode{
		Left:     left,
		Operator: p.curToken.Value,
	}

	precedence := p.curPrecedence()
	p.nextToken()

	right, err := p.ParseNode(precedence)
	if err != nil {
		return nil, err
	}
	Node.Right = right
	return Node, nil
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

func (p *V1Parser) parseFunctionLiteral() (Node, error) {

	if p.Debug {
		fmt.Println("Entering function literal ")
	}
	fl := &FunctionLiteral{}
	p.nextToken()

	fl.Name = p.curToken.Value

	p.nextToken()

	if !p.peekTokenIs(RPAREN) && !p.curTokenIs(LBRACE) {
		p.nextToken()
		for !p.peekTokenIs(RPAREN) {
			fmt.Println(p.curToken.Value)
			param := &VariableDeclaration{Identifier: &IdentifierLiteral{value: p.curToken.Value}}
			p.nextToken()
			param.Type = p.curToken
			p.nextToken()
			fmt.Println(p.curToken.Value)
			// p.nextToken()
			// fl.Parameters = append(fl.Parameters, param)
		}
		p.nextToken()

	}

	p.nextToken()
	p.nextToken()

	for !p.peekTokenIs(RBRACE) && !p.curTokenIs(LBRACE) {

		var expr Node
		var err error
		if p.curToken.Type == RETURN {
			expr, err = p.parseReturnStatement()
			if err != nil {
				return nil, err
			}
		} else {
			expr, err = p.ParseNode(LOWEST)
			if err != nil {
				return nil, err
			}
		}
		fl.Body = append(fl.Body, expr)
	}

	if len(fl.Body) > 0 {
		if _, ok := fl.Body[len(fl.Body)-1].(*ReturnStatement); !ok {
			fl.Body = append(fl.Body, &ReturnStatement{})
		}
	} else {
		fl.Body = append(fl.Body, &ReturnStatement{})
	}
	return fl, nil
}

func (p *V1Parser) parseArrayLiteral() (Node, error) {
	array := &ArrayLiteral{}

	if p.peekTokenIs(RBRACKET) {
		p.nextToken()
		return array, nil
	}

	p.nextToken()
	node, err := p.ParseNode(LOWEST)
	if err != nil {
		return nil, err
	}
	array.Elements = append(array.Elements, node)

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		node, err = p.ParseNode(LOWEST)
		if err != nil {
			return nil, err
		}
		array.Elements = append(array.Elements, node)
	}

	if !p.expectPeek(RBRACKET) {
		return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
	}

	return array, nil
}

func (p *V1Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
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

	for {
		node, err := p.ParseNode(LOWEST)
		if err != nil {
			return nil, err
		}

		components = append(components, node)

		if p.peekTokenIs(SEMICOLON) {
			p.nextToken()
			p.nextToken()
		} else if p.peekTokenIs(LBRACE) || p.curTokenIs(LBRACE) {
			break
		} else {
			return nil, fmt.Errorf(SYNTAX_ERROR_MSG, p.curToken.Line)
		}
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
	forExp.Body = block.Statements

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

func (p *V1Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}
