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

// func (p *Program) TokenLiteral() string {
// 	if len(p.Statements) > 0 {
// 		return p.Statements[0].TokenLiteral()
// 	} else {
// 		return ""
// 	}
// }
//
// func (p *Program) TokenType() token.TokenType {
// 	if len(p.Statements) > 0 {
// 		return p.Statements[0].TokenType()
// 	} else {
// 		return ""
// 	}
// }

func (p *Program) String() string {
	var b strings.Builder

	for _, stmt := range p.Statements {
		b.WriteString(stmt.String())
	}
	return b.String()
}

func (p *Program) ASTDebugString() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Program {\n")
	for _, stmt := range p.Statements {
		stmt.ASTDebugString(&b, 1)
	}
	fmt.Fprintf(&b, "}\n")
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

type PrefixExpression struct {
	Token    token.Token // the token.BANG token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()            {}
func (pe *PrefixExpression) TokenType() token.TokenType { return pe.Token.Type }
func (pe *PrefixExpression) TokenLiteral() string       { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(pe.Operator)
	if pe.Right != nil {
		b.WriteString(pe.Right.String())
	}
	b.WriteString(")")
	return b.String()
}

func (pe *PrefixExpression) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, "PrefixExpression {\n", depth)
	debugWrite(w, fmt.Sprintf("Operator (%s)\n", pe.Operator), depth+1)
	if pe.Right != nil {
		debugWrite(w, "Right:\n", depth+1)
		pe.Right.ASTDebugString(w, depth+1)
	} else {
		debugWrite(w, "Right (nil)\n", depth+1)
	}
	debugWrite(w, "}\n", depth)
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

func debugWrite(w io.Writer, str string, depth int) {
	fmt.Fprintf(w, "%s%s", strings.Repeat(" ", depth*2), str)
}
