// Package ast
package ast

import (
	"monkey/token"
	"strings"
)

type Node interface {
	TokenType() token.TokenType
	TokenLiteral() string
	String() string
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

	b.WriteString(ls.Name.String())
	b.WriteString(" = ")
	if ls.Value != nil {
		b.WriteString(ls.Value.String())
	}
	b.WriteString(";")
	return b.String()
}

type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
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

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()            {}
func (i *Identifier) TokenType() token.TokenType { return i.Token.Type }
func (i *Identifier) TokenLiteral() string       { return i.Token.Literal }
func (i *Identifier) String() string             { return i.Value }

type IntegerLiteral struct {
	Token token.Token // the token.INT token
	Value int64
}

func (il *IntegerLiteral) expressionNode()            {}
func (il *IntegerLiteral) TokenType() token.TokenType { return il.Token.Type }
func (il *IntegerLiteral) TokenLiteral() string       { return il.Token.Literal }
func (il *IntegerLiteral) String() string             { return il.Token.Literal }

type Boolean struct {
	Token token.Token // the token.TRUE or FALSE token
	Value bool
}

func (b *Boolean) expressionNode()            {}
func (b *Boolean) TokenType() token.TokenType { return b.Token.Type }
func (b *Boolean) TokenLiteral() string       { return b.Token.Literal }
func (b *Boolean) String() string             { return b.Token.Literal }

type StringLiteral struct {
	Token token.Token // the token.INT token
	Value string
}

func (il *StringLiteral) expressionNode()            {}
func (il *StringLiteral) TokenType() token.TokenType { return il.Token.Type }
func (il *StringLiteral) TokenLiteral() string       { return il.Token.Literal }
func (il *StringLiteral) String() string             { return il.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token // the token.INT token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()            {}
func (al *ArrayLiteral) TokenType() token.TokenType { return al.Token.Type }
func (al *ArrayLiteral) TokenLiteral() string       { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out strings.Builder
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type IndexExpression struct {
	Token token.Token // the token.INT token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()            {}
func (ie *IndexExpression) TokenType() token.TokenType { return ie.Token.Type }
func (ie *IndexExpression) TokenLiteral() string       { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out strings.Builder
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	out.WriteString(")")
	return out.String()
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
