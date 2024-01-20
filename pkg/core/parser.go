package core

import (
	"fmt"
	"strconv"

)

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
	ParseProgram() []Node
	ParseNode(int) Node
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
	prefixParseFn func() Node
	infixParseFn  func(Node) Node
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

func (p *V1Parser) ParseProgram() []Node {
	program := []Node{}
	for !p.curTokenIs(EOF) {
		exp := p.ParseNode(LOWEST)
		if exp != nil {
			program = append(program, exp)
		}
		p.nextToken()
	}
	return program
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

func (p *V1Parser) ParseNode(precedence int) Node {
	if p.Debug {
		fmt.Printf("DEBUG: ParseNode(): Precedence: %d, Current Token: %+v\n", precedence, p.curToken)
	}
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(EOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *V1Parser) parseImport() Node {
	if p.Debug {
		fmt.Println("Entering parseImport")
	}

	paths := []string{}

	if !p.peekTokenIs(LPAREN) {
		if !p.expectPeek(STRING) {
			fmt.Printf("Error: bad import statement on Line: %d", p.peekToken.Line)
			return nil
		}
		paths = append(paths, p.curToken.Value)
	} else {
		p.nextToken()

		for !p.peekTokenIs(RPAREN) {
			if !p.expectPeek(STRING) {
				fmt.Printf("Error: bad import statement on Line: %d", p.peekToken.Line)
				return nil
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

	return moduleListNode
}

// parseBooleanLiteral parses the 'true' or 'false' token and creates a BooleanLiteral Node
func (p *V1Parser) parseBooleanLiteral() Node {
	literal := &BooleanLiteral{
		value:  p.curToken.Type == TRUE,
		Line:   p.curToken.Line,
		Column: p.curToken.Column,
	}
	p.nextToken()
	return literal
}

func (p *V1Parser) parseLeftParen() Node {
	if p.Debug {
		fmt.Println("Entering parseLeftParen")
	}

	p.nextToken()

	if p.peekTokenIs(IDENT) {
		if p.Debug {
			fmt.Println("Detected IDENT token")
		}

		ident := p.parseIdentifier().(*IdentifierLiteral)
		if p.Debug {
			fmt.Printf("Parsed IDENT: %v\n", ident.value)
		}

		fc := &FunctionCall{
			Name:      ident.value,
			Arguments: []Node{},
		}

		p.nextToken()

		for !p.curTokenIs(RPAREN) && !p.peekTokenIs(EOF) {
			p.nextToken()
			arg := p.ParseNode(LOWEST)
			fc.Arguments = append(fc.Arguments, arg)

			if p.peekTokenIs(COMMA) {
				p.nextToken()
			}

			if p.curTokenIs(RPAREN) {
				break
			}
		}

		if !p.curTokenIs(RPAREN) {
			return nil
		}

		p.nextToken()

		if p.Debug {
			fmt.Println("Exiting parseLeftParen with function call")
		}

		return fc
	}

	expr := p.ParseNode(LOWEST)

	if p.Debug {
		fmt.Printf("Parsed Node: %v\n", expr)
	}

	if !p.expectPeek(RPAREN) {
		if p.Debug {
			fmt.Println("Missing RPAREN token")
		}

		return nil
	}

	if p.Debug {
		fmt.Println("Exiting parseLeftParen with Node")
	}

	return expr
}

func (p *V1Parser) parseStringLiteral() Node {
	return &StringLiteral{value: p.curToken.Value}
}

func (p *V1Parser) parseIntegerLiteral() Node {
	lit := &IntegerLiteral{}
	value, err := strconv.Atoi(p.curToken.Value)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Value)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.value = value

	return lit
}

func (p *V1Parser) parseFloatLiteral() Node {
	value, err := strconv.ParseFloat(p.curToken.Value, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Value)
		p.errors = append(p.errors, msg)
		return nil
	}

	return &FloatLiteral{value: value}
}

func (p *V1Parser) parseHashLiteral() Node {
	hash := &HashLiteral{Pairs: make(map[Node]Node)}

	for !p.peekTokenIs(RBRACE) {
		p.nextToken()
		key := p.ParseNode(LOWEST)

		if !p.expectPeek(COLON) {
			return nil
		}

		p.nextToken()
		value := p.ParseNode(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(RBRACE) && !p.expectPeek(COMMA) {
			return nil
		}
	}

	if !p.expectPeek(RBRACE) {
		return nil
	}

	return hash
}

func (p *V1Parser) parseIdentifier() Node {
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
		fc := p.parseFunctionCall(ident)

		if p.Debug {
			fmt.Printf("Parsed function call: %v\n", fc)
			fmt.Println("Exiting parseIdentifier with function call")
		}

		return fc
	}

	if p.peekTokenIs(DOT) {
		if p.Debug {
			fmt.Println("Detected DOT token")
		}

		// Note: Here, instead of nextToken(), we're using a new function parseDotNotation()
		// which will handle the parsing of everything following the DOT.
		return p.parseDotNotation(ident)
	}

	if p.peekTokenIs(INC) {
		p.nextToken()
		return &IncrementNode{Operand: ident}
	}
	if p.peekTokenIs(DEC) {
		p.nextToken()
		return &DecrementNode{Operand: ident}
	}

	if p.Debug {
		fmt.Println("Exiting parseIdentifier with ident")
	}

	return ident
}

func (p *V1Parser) parseDotNotation(left Node) Node {
	if p.Debug {
		fmt.Println("Entering parseDotNotation")
	}

	// Advancing past the DOT token
	p.nextToken()

	// Creating a new DotNotationNode with the left part (before the dot)
	// and the right part (the next identifier after the dot)
	dotNotationNode := &DotNotationNode{
		Left:   left,
		Right:  p.parseIdentifier(),
		Line:   left.GetLine(),
		Column: left.GetColumn(),
	}

	if p.Debug {
		fmt.Printf("Parsed DotNotation: %v\n", dotNotationNode)
		fmt.Println("Exiting parseDotNotation with DotNotationNode")
	}

	return dotNotationNode
}

func (p *V1Parser) parseFunctionCall(function Node) Node {
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
		fc.Arguments = append(fc.Arguments, p.ParseNode(LOWEST))

		for p.peekTokenIs(COMMA) {
			p.nextToken() // This will consume the comma
			p.nextToken() // This will move us to the next argument
			fc.Arguments = append(fc.Arguments, p.ParseNode(LOWEST))
		}

		if !p.expectPeek(RPAREN) {
			return nil
		}
	}

	if !p.curTokenIs(RPAREN) {
		if p.Debug {
			fmt.Println("Current token is not RPAREN, returning nil")
		}
		return nil
	}

	if p.Debug {
		fmt.Println("Current token is RPAREN, moving to next token")
	}

	p.nextToken()

	if p.Debug {
		fmt.Println("Exiting parseFunctionCall")
	}

	return fc
}

func (p *V1Parser) parseInfixNode(left Node) Node {
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

	Node.Right = p.ParseNode(precedence)

	if p.Debug {
		fmt.Printf("Right Node: %v\n", Node.Right)
		fmt.Println("Exiting parseInfixNode")
	}

	return Node
}

func (p *V1Parser) parsePrefixNode() Node {
	Node := &PrefixNode{
		Operator: p.curToken.Value,
	}

	p.nextToken()

	Node.Right = p.ParseNode(PREFIX)

	return Node
}

func (p *V1Parser) parseEqualityInequalityNode(left Node) Node {
	Node := &InfixNode{
		Left:     left,
		Operator: p.curToken.Value,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	Node.Right = p.ParseNode(precedence)

	return Node
}

func (p *V1Parser) parseAssignNode(left Node) Node {
	Node := &InfixNode{
		Left:     left,
		Operator: p.curToken.Value,
	}

	precedence := p.curPrecedence()
	p.nextToken()

	Node.Right = p.ParseNode(precedence)

	return Node
}

func (p *V1Parser) noPrefixParseFnError(t TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", TokenTypeStr[t])
	p.errors = append(p.errors, msg)
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

func (p *V1Parser) parseFunctionLiteral() Node {

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

	for !p.peekTokenIs(RBRACE) && !p.curTokenIs(LBRACE){

		var expr Node
		if p.curToken.Type == RETURN {
			expr = p.parseReturnStatement()
		} else {
			expr = p.ParseNode(LOWEST)
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
	return fl
}

func (p *V1Parser) parseArrayLiteral() Node {
	array := &ArrayLiteral{}

	if p.peekTokenIs(RBRACKET) {
		p.nextToken()
		return array
	}

	p.nextToken()
	array.Elements = append(array.Elements, p.ParseNode(LOWEST))

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		array.Elements = append(array.Elements, p.ParseNode(LOWEST))
	}

	if !p.expectPeek(RBRACKET) {
		return nil
	}

	return array
}

func (p *V1Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *V1Parser) parseReturnStatement() Node {

	p.nextToken()

	rs := &ReturnStatement{}

	rs.ReturnValue = p.ParseNode(LOWEST)

	if p.peekTokenIs(SEMICOLON) {
		p.nextToken()
	}

	return rs
}

func (p *V1Parser) parseIfStatement() Node {

	p.nextToken()

	ifExp := &IfNode{}

	ifExp.Condition = p.ParseNode(LOWEST)

	if !p.expectPeek(LBRACE) {
		return nil
	}

	ifExp.Consequence = p.parseBlockStatement().Statements

	if p.peekTokenIs(ELSE) {
		p.nextToken()

		if !p.expectPeek(LBRACE) {
			return nil
		}

		ifExp.Alternative = p.parseBlockStatement().Statements
	}

	return ifExp
}

func (p *V1Parser) parseForStatement() Node {

	if p.Debug {
		fmt.Println("Entering parseForStatement")
	}

	forExp := &ForNode{}

	components := []Node{}

	for {
		components = append(components, p.ParseNode(LOWEST))

		if p.peekTokenIs(SEMICOLON) {
			p.nextToken() 
			p.nextToken() 
		} else if p.peekTokenIs(LBRACE) || p.curTokenIs(LBRACE) {
			break
		} else {
			// TODO: how do we raise SyntaxError?
			return nil
		}
	}

	if len(components) == 2 || len(components) > 3 {
		// TODO: how do we raise SyntaxError?
		return nil
	}

	if len(components) == 1 {
		forExp.Condition = components[0]
	} else if len(components) == 3 {
		forExp.Initialisation = components[0]
		forExp.Condition = components[1]
		forExp.Updater = components[2]
	}

	forExp.Body = p.parseBlockStatement().Statements

	if p.Debug {
		fmt.Println("Exiting parseForStatement")
	}

	return forExp
}

func (p *V1Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{}
	block.Statements = []Node{}

	p.nextToken()

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		stmt := p.ParseNode(LOWEST)
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *V1Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}

