package parser

import (
	"fmt"
	"strconv"

	"github.com/hyperioxx/goscript/pkg/core/lexer"
)

const (
	_ int = iota
	LOWEST
	ASSIGN
	IF
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X

	ARRAY
	CALL // foo(X)
)

var precedences = map[lexer.TokenType]int{
	lexer.EQ:         EQUALS,
	lexer.NOT_EQ:     EQUALS,
	lexer.LT:         LESSGREATER,
	lexer.GT:         LESSGREATER,
	lexer.ADD:        SUM,
	lexer.SUB:        SUM,
	lexer.MUL:        PRODUCT,
	lexer.DIV:        PRODUCT,
	lexer.LPAREN:     CALL,
	lexer.ASSIGN:     ASSIGN,
	lexer.ASSIGN_INF: ASSIGN,
	lexer.IF:         IF,
	lexer.LBRACKET:   ARRAY,
}

type Parser interface {
	ParseProgram() []Node
	ParseNode(int) Node
}

type V1Parser struct {
	l lexer.Lexer

	curToken        lexer.Token
	peekToken       lexer.Token
	errors          []string
	Debug           bool
	lastParsedIdent string
	prefixParseFns  map[lexer.TokenType]prefixParseFn
	infixParseFns   map[lexer.TokenType]infixParseFn
}

type (
	prefixParseFn func() Node
	infixParseFn  func(Node) Node
)

func NewV1Parser(l lexer.Lexer, debug bool) Parser {
	p := &V1Parser{
		l:      l,
		errors: []string{},
		Debug:  debug,
	}

	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)

	p.registerInfix(lexer.ADD, p.parseInfixNode)
	p.registerInfix(lexer.SUB, p.parseInfixNode)
	p.registerInfix(lexer.MUL, p.parseInfixNode)
	p.registerInfix(lexer.DIV, p.parseInfixNode)
	p.registerInfix(lexer.EQ, p.parseInfixNode)
	p.registerInfix(lexer.NOT_EQ, p.parseInfixNode)
	p.registerInfix(lexer.GT, p.parseInfixNode)
	p.registerInfix(lexer.LT, p.parseInfixNode)
	p.registerInfix(lexer.GT_EQ, p.parseInfixNode)
	p.registerInfix(lexer.LT_EQ, p.parseInfixNode)
	p.registerInfix(lexer.ASSIGN, p.parseAssignNode)
	p.registerInfix(lexer.ASSIGN_INF, p.parseAssignNode)
	// prefix expressions
	p.registerPrefix(lexer.INT, p.parseIntegerLiteral)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.FUNC, p.parseFunctionLiteral)
	p.registerPrefix(lexer.LPAREN, p.parseLeftParen)
	p.registerPrefix(lexer.RETURN, p.parseReturnStatement)
	p.registerPrefix(lexer.IF, p.parseIfStatement)
	p.registerPrefix(lexer.FOR, p.parseForStatement)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(lexer.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(lexer.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.IMPORT, p.parseImport)
	p.registerPrefix(lexer.STRUCT, p.parseStructLiteral)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *V1Parser) ParseProgram() []Node {
	program := []Node{}
	for !p.curTokenIs(lexer.EOF) {
		exp := p.ParseNode(LOWEST)
		if exp != nil {
			program = append(program, exp)
		}
		p.nextToken()
	}
	return program
}

func (p *V1Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *V1Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
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
	for !p.peekTokenIs(lexer.EOF) && precedence < p.peekPrecedence() {
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

	if !p.peekTokenIs(lexer.LPAREN) {
		if !p.expectPeek(lexer.STRING) {
			fmt.Printf("Error: bad import statement on Line: %d", p.peekToken.Line)
			return nil
		}
		paths = append(paths, p.curToken.Value)
	} else {
		p.nextToken()

		for !p.peekTokenIs(lexer.RPAREN) {
			if !p.expectPeek(lexer.STRING) {
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
		value:  p.curToken.Type == lexer.TRUE,
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

	if p.peekTokenIs(lexer.IDENT) {
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

		for !p.curTokenIs(lexer.RPAREN) && !p.peekTokenIs(lexer.EOF) {
			p.nextToken()
			arg := p.ParseNode(LOWEST)
			fc.Arguments = append(fc.Arguments, arg)

			if p.peekTokenIs(lexer.COMMA) {
				p.nextToken()
			}

			if p.curTokenIs(lexer.RPAREN) {
				break
			}
		}

		if !p.curTokenIs(lexer.RPAREN) {
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

	if !p.expectPeek(lexer.RPAREN) {
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

	p.nextToken()

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

	for !p.peekTokenIs(lexer.RBRACE) {
		p.nextToken()
		key := p.ParseNode(LOWEST)

		if !p.expectPeek(lexer.COLON) {
			return nil
		}

		p.nextToken()
		value := p.ParseNode(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(lexer.RBRACE) && !p.expectPeek(lexer.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(lexer.RBRACE) {
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

	if p.peekTokenIs(lexer.LPAREN) {
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

	if p.peekTokenIs(lexer.DOT) {
		if p.Debug {
			fmt.Println("Detected DOT token")
		}

		// Note: Here, instead of nextToken(), we're using a new function parseDotNotation()
		// which will handle the parsing of everything following the DOT.
		return p.parseDotNotation(ident)
	}

	if p.peekTokenIs(lexer.INC) {
		p.nextToken()
		return &IncrementNode{Operand: ident}
	}
	if p.peekTokenIs(lexer.DEC) {
		p.nextToken()
		return &DecrementNode{Operand: ident}
	}

	if p.peekTokenIs(lexer.INT) || p.peekTokenIs(lexer.STRING) || p.peekTokenIs(lexer.FLOAT){
       
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

	if !p.curTokenIs(lexer.RPAREN) {
		fc.Arguments = append(fc.Arguments, p.ParseNode(LOWEST))

		for p.peekTokenIs(lexer.COMMA) {
			p.nextToken() // This will consume the comma
			p.nextToken() // This will move us to the next argument
			fc.Arguments = append(fc.Arguments, p.ParseNode(LOWEST))
		}

		if !p.expectPeek(lexer.RPAREN) {
			return nil
		}
	}

	if !p.curTokenIs(lexer.RPAREN) {
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

func (p *V1Parser) noPrefixParseFnError(t lexer.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", lexer.TokenTypeStr[t])
	p.errors = append(p.errors, msg)
}

func (p *V1Parser) peekTokenIs(t lexer.TokenType) bool {
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
	fl := &FunctionLiteral{}
	p.nextToken()
	// Check if the current token is an opening parenthesis
	if p.curTokenIs(lexer.LPAREN) {
		fl.Name = "_anonymous" // placeholder name for nameless functions we would have to generate a unique name as you can only have 1 nameless function per scope
	} else {
		fl.Name = p.curToken.Value
		p.nextToken()
	}

	if !p.peekTokenIs(lexer.RPAREN) {
		for !p.peekTokenIs(lexer.RPAREN) {
			p.nextToken()
			param := &IdentifierLiteral{value: p.curToken.Value}
			fl.Parameters = append(fl.Parameters, param)
			p.nextToken()
			if p.curToken.Type == lexer.COMMA {
				p.nextToken()
			}
		}

		// final arg
		param := &IdentifierLiteral{value: p.curToken.Value}
		fl.Parameters = append(fl.Parameters, param)
	}

	p.nextToken()
	p.nextToken() // should be this symbol {

	for !p.peekTokenIs(lexer.RBRACE) {
		p.nextToken()
		var expr Node
		if p.curToken.Type == lexer.RETURN {
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

	if p.peekTokenIs(lexer.RBRACKET) {
		p.nextToken()
		return array
	}

	p.nextToken()
	array.Elements = append(array.Elements, p.ParseNode(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		array.Elements = append(array.Elements, p.ParseNode(LOWEST))
	}

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

	return array
}

func (p *V1Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *V1Parser) parseReturnStatement() Node {

	p.nextToken()

	rs := &ReturnStatement{}

	rs.ReturnValue = p.ParseNode(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return rs
}

func (p *V1Parser) parseIfStatement() Node {

	p.nextToken()

	ifExp := &IfNode{}

	ifExp.Condition = p.ParseNode(LOWEST)

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	ifExp.Consequence = p.parseBlockStatement().Statements

	if p.peekTokenIs(lexer.ELSE) {
		p.nextToken()

		if !p.expectPeek(lexer.LBRACE) {
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
	p.nextToken()

	forExp := &ForNode{}

	components := []Node{}

	// gathers up all semicolon-delimited Nodes preceding the LBRACE
	for {
		components = append(components, p.ParseNode(LOWEST))

		if p.peekTokenIs(lexer.SEMICOLON) {
			p.nextToken() // whatever ended the Node
			p.nextToken() // the semicolon
		} else if p.peekTokenIs(lexer.LBRACE) || p.curTokenIs(lexer.LBRACE) {
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

	// nothing to do if `len(components) == 0`, VM will understand what this means
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

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		stmt := p.ParseNode(LOWEST)
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *V1Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *V1Parser) parseStructLiteral() Node {
	structLiteral := &StructLiteral{
		Fields: make(map[string]Node),
		Line:   p.curToken.Line,
		Column: p.curToken.Column,
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	p.nextToken()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		if !p.curTokenIs(lexer.IDENT) {
			msg := fmt.Sprintf("expected identifier, got %s", p.curToken.Type)
			p.errors = append(p.errors, msg)
			return nil
		}

		fieldName := p.curToken.Value

		if !p.expectPeek(lexer.COLON) {
			return nil
		}

		p.nextToken()

		fieldValue := p.ParseNode(LOWEST)

		structLiteral.Fields[fieldName] = fieldValue

		if !p.peekTokenIs(lexer.RBRACE) && !p.expectPeek(lexer.COMMA) {
			return nil
		}

		p.nextToken()
	}

	if !p.curTokenIs(lexer.RBRACE) {
		msg := fmt.Sprintf("expected }, got %s", p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	p.nextToken()

	return structLiteral
}
