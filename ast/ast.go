// Package ast
package ast

import (
	"fmt"
	"io"
	"monkey/token"
	"strings"
)

type Node interface {
	TokenType() token.TokenType
	TokenLiteral() string
	String() string
	ASTDebugString(out io.Writer, depth int)
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) TokenType() token.TokenType {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenType()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var b strings.Builder

	for _, stmt := range p.Statements {
		b.WriteString(stmt.String())
	}
	return b.String()
}

func (p *Program) ASTDebugString(w io.Writer, depth int) {
	fmt.Fprintf(w, "Program {\n")
	for _, stmt := range p.Statements {
		stmt.ASTDebugString(w, depth)
	}
	fmt.Fprintf(w, "}\n")
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()             {}
func (ls *LetStatement) TokenType() token.TokenType { return ls.Token.Type }
func (ls *LetStatement) TokenLiteral() string       { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var b strings.Builder

	b.WriteString(ls.TokenLiteral() + " ")
	b.WriteString(ls.Name.String())
	b.WriteString(" = ")
	if ls.Value != nil {
		b.WriteString(ls.Value.String())
	}
	b.WriteString(";")
	return b.String()
}

func (ls *LetStatement) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, "LetStatement {\n", depth)
	if ls.Name != nil {
		ls.Name.ASTDebugString(w, depth+1)
	} else {
		debugWrite(w, "Name (nil)\n", depth+1)
	}
	if ls.Value != nil {
		ls.Value.ASTDebugString(w, depth+1)
	} else {
		debugWrite(w, "Value (nil)\n", depth+1)
	}
	debugWrite(w, "}\n", depth)
}

type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

type AssignmentStatement struct {
	Token token.Token // the token.IDENT token
	Name  *Identifier
	Value Expression
}

func (ls *AssignmentStatement) statementNode()             {}
func (ls *AssignmentStatement) TokenType() token.TokenType { return ls.Token.Type }
func (ls *AssignmentStatement) TokenLiteral() string       { return ls.Token.Literal }
func (ls *AssignmentStatement) String() string {
	var b strings.Builder

	b.WriteString(ls.TokenLiteral() + " ")
	b.WriteString(ls.Name.String())
	b.WriteString(" = ")
	if ls.Value != nil {
		b.WriteString(ls.Value.String())
	}
	b.WriteString(";")
	return b.String()
}

func (ls *AssignmentStatement) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, "AssignmentStatement {\n", depth)
	if ls.Name != nil {
		ls.Name.ASTDebugString(w, depth+1)
	} else {
		debugWrite(w, "Name (nil)\n", depth+1)
	}
	if ls.Value != nil {
		ls.Value.ASTDebugString(w, depth+1)
	} else {
		debugWrite(w, "Value (nil)\n", depth+1)
	}
	debugWrite(w, "}\n", depth)
}

func (rs *ReturnStatement) statementNode()             {}
func (rs *ReturnStatement) TokenType() token.TokenType { return rs.Token.Type }
func (rs *ReturnStatement) TokenLiteral() string       { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var b strings.Builder
	b.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		b.WriteString(rs.ReturnValue.String())
	}
	b.WriteString(";")
	return b.String()
}

func (rs *ReturnStatement) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, "ReturnStatement {\n", depth)
	if rs.ReturnValue != nil {
		rs.ReturnValue.ASTDebugString(w, depth+1)
	} else {
		debugWrite(w, "ReturnValue (nil)\n", depth+1)
	}
	debugWrite(w, "}\n", depth)
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()             {}
func (es *ExpressionStatement) TokenType() token.TokenType { return es.Token.Type }
func (es *ExpressionStatement) TokenLiteral() string       { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (es *ExpressionStatement) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, "ExpressionStatement {\n", depth)
	if es.Expression != nil {
		es.Expression.ASTDebugString(w, depth+1)
	} else {
		debugWrite(w, "Expression (nil)\n", depth+1)
	}
	debugWrite(w, "}\n", depth)
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()             {}
func (bs *BlockStatement) TokenType() token.TokenType { return bs.Token.Type }
func (bs *BlockStatement) TokenLiteral() string       { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var b strings.Builder

	for _, stmt := range bs.Statements {
		b.WriteString(stmt.String())
	}
	return b.String()
}

func (bs *BlockStatement) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, "BlockStatement {\n", depth)
	for _, stmt := range bs.Statements {
		stmt.ASTDebugString(w, depth+1)
	}
	debugWrite(w, "}\n", depth)
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()            {}
func (i *Identifier) TokenType() token.TokenType { return i.Token.Type }
func (i *Identifier) TokenLiteral() string       { return i.Token.Literal }
func (i *Identifier) String() string             { return i.Value }
func (i *Identifier) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, fmt.Sprintf("Identifier (%s)\n", i.Value), depth)
}

type IntegerLiteral struct {
	Token token.Token // the token.INT token
	Value int64
}

func (il *IntegerLiteral) expressionNode()            {}
func (il *IntegerLiteral) TokenType() token.TokenType { return il.Token.Type }
func (il *IntegerLiteral) TokenLiteral() string       { return il.Token.Literal }
func (il *IntegerLiteral) String() string             { return il.Token.Literal }
func (il *IntegerLiteral) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, fmt.Sprintf("INT (%d)\n", il.Value), depth)
}

type Boolean struct {
	Token token.Token // the token.TRUE or FALSE token
	Value bool
}

func (b *Boolean) expressionNode()            {}
func (b *Boolean) TokenType() token.TokenType { return b.Token.Type }
func (b *Boolean) TokenLiteral() string       { return b.Token.Literal }
func (b *Boolean) String() string             { return b.Token.Literal }
func (b *Boolean) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, fmt.Sprintf("BOOL (%t)\n", b.Value), depth)
}

type PrefixExpression struct {
	Token    token.Token // the token.BANG token
	Operator string
	Right    Expression
}

func (exp *PrefixExpression) expressionNode()            {}
func (exp *PrefixExpression) TokenType() token.TokenType { return exp.Token.Type }
func (exp *PrefixExpression) TokenLiteral() string       { return exp.Token.Literal }
func (exp *PrefixExpression) String() string {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(exp.Operator)
	if exp.Right != nil {
		b.WriteString(exp.Right.String())
	}
	b.WriteString(")")
	return b.String()
}

func (exp *PrefixExpression) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, "PrefixExpression {\n", depth)
	debugWrite(w, fmt.Sprintf("Operator (%s)\n", exp.Operator), depth+1)
	if exp.Right != nil {
		exp.Right.ASTDebugString(w, depth+1)
	} else {
		debugWrite(w, "Right (nil)\n", depth+1)
	}
	debugWrite(w, "}\n", depth)
}

type InfixExpression struct {
	Token    token.Token // the token.BANG token
	Left     Expression
	Operator string
	Right    Expression
}

func (exp *InfixExpression) expressionNode()            {}
func (exp *InfixExpression) TokenType() token.TokenType { return exp.Token.Type }
func (exp *InfixExpression) TokenLiteral() string       { return exp.Token.Literal }
func (exp *InfixExpression) String() string {
	var b strings.Builder
	b.WriteString("(")
	if exp.Left != nil {
		b.WriteString(exp.Left.String())
	}
	b.WriteString(" " + exp.Operator + " ")
	if exp.Right != nil {
		b.WriteString(exp.Right.String())
	}
	b.WriteString(")")
	return b.String()
}

func (exp *InfixExpression) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, "InfixExpression {\n", depth)
	if exp.Left != nil {
		exp.Left.ASTDebugString(w, depth+1)
	} else {
		debugWrite(w, "Left (nil)\n", depth+1)
	}
	debugWrite(w, fmt.Sprintf("Operator (%s)\n", exp.Operator), depth+1)
	if exp.Right != nil {
		exp.Right.ASTDebugString(w, depth+1)
	} else {
		debugWrite(w, "Right (nil)\n", depth+1)
	}
	debugWrite(w, "}\n", depth)
}

// TODO: Implement PostfixExpressions these when feeling comfortable
// an example of postfix operations are: increment and decrements `i++` `i--`

type PostfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
}

func (exp *PostfixExpression) expressionNode()            {}
func (exp *PostfixExpression) TokenType() token.TokenType { return exp.Token.Type }
func (exp *PostfixExpression) TokenLiteral() string       { return exp.Token.Literal }
func (exp *PostfixExpression) String() string {
	var b strings.Builder
	b.WriteString("(")
	if exp.Left != nil {
		b.WriteString(exp.Left.String())
	}
	b.WriteString(exp.Operator)
	b.WriteString(")")
	return b.String()
}

func (exp *PostfixExpression) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, "PostfixExpression {\n", depth)
	if exp.Left != nil {
		exp.Left.ASTDebugString(w, depth+1)
	} else {
		debugWrite(w, "Left (nil)\n", depth+1)
	}
	debugWrite(w, fmt.Sprintf("Operator (%s)\n", exp.Operator), depth+1)
	debugWrite(w, "}\n", depth)
}

type IfExpression struct {
	Token       token.Token // the token.IF token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (exp *IfExpression) expressionNode()            {}
func (exp *IfExpression) TokenType() token.TokenType { return exp.Token.Type }
func (exp *IfExpression) TokenLiteral() string       { return exp.Token.Literal }
func (exp *IfExpression) String() string {
	var b strings.Builder
	b.WriteString("if")
	b.WriteString(exp.Condition.String())
	b.WriteString(" ")
	b.WriteString(exp.Consequence.String())
	if exp.Alternative != nil {
		b.WriteString(" else ")
		b.WriteString(exp.Alternative.String())
	}
	return b.String()
}

func (exp *IfExpression) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, "IfExpression {\n", depth)
	exp.Condition.ASTDebugString(w, depth+1)
	if exp.Consequence != nil {
		exp.Consequence.ASTDebugString(w, depth+1)
	} else {
		debugWrite(w, "Consequence (nil)\n", depth+1)
	}
	if exp.Alternative != nil {
		exp.Alternative.ASTDebugString(w, depth+1)
	}
	debugWrite(w, "}\n", depth)
}

type FunctionLiteral struct {
	Token      token.Token // the token.IF token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (exp *FunctionLiteral) expressionNode()            {}
func (exp *FunctionLiteral) TokenType() token.TokenType { return exp.Token.Type }
func (exp *FunctionLiteral) TokenLiteral() string       { return exp.Token.Literal }
func (exp *FunctionLiteral) String() string {
	var b strings.Builder
	params := []string{}
	for _, param := range exp.Parameters {
		params = append(params, param.String())
	}
	b.WriteString("fn")
	b.WriteString("(")
	b.WriteString(strings.Join(params, ", "))
	b.WriteString(") ")
	if exp.Body != nil {
		b.WriteString(exp.Body.String())
	}
	return b.String()
}

func (exp *FunctionLiteral) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, "FunctionLiteral {\n", depth)
	if len(exp.Parameters) > 0 {
		for _, param := range exp.Parameters {
			param.ASTDebugString(w, depth+1)
		}
	}
	if exp.Body != nil {
		exp.Body.ASTDebugString(w, depth+1)
	}
	debugWrite(w, "}\n", depth)
}

type CallExpression struct {
	Token     token.Token // the token.IF token
	Function  Expression
	Arguments []Expression
}

func (exp *CallExpression) expressionNode()            {}
func (exp *CallExpression) TokenType() token.TokenType { return exp.Token.Type }
func (exp *CallExpression) TokenLiteral() string       { return exp.Token.Literal }
func (exp *CallExpression) String() string {
	var b strings.Builder
	args := []string{}
	for _, arg := range exp.Arguments {
		args = append(args, arg.String())
	}

	b.WriteString(exp.Function.String())
	b.WriteString("(")
	b.WriteString(strings.Join(args, ", "))
	b.WriteString(")")
	return b.String()
}

func (exp *CallExpression) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, "CallExpression {\n", depth)
	if exp.Function != nil {
		exp.Function.ASTDebugString(w, depth+1)
	}
	if len(exp.Arguments) > 0 {
		for _, arg := range exp.Arguments {
			arg.ASTDebugString(w, depth+1)
		}
	}
	debugWrite(w, "}\n", depth)
}

func debugWrite(w io.Writer, str string, depth int) {
	fmt.Fprintf(w, "%s%s", strings.Repeat(" ", depth*2), str)
}
