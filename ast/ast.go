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
func (i *Identifier) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, fmt.Sprintf("Identifier (%s)\n", i.Value), depth)
}

type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()             {}
func (rs *ReturnStatement) TokenType() token.TokenType { return rs.Token.Type }
func (rs *ReturnStatement) TokenLiteral() string       { return rs.Token.Literal }
func (rs *ReturnStatement) ASTDebugString(w io.Writer, depth int) {
	debugWrite(w, "ReturnStatement {\n", depth)
	if rs.ReturnValue != nil {
		rs.ReturnValue.ASTDebugString(w, depth+1)
	} else {
		debugWrite(w, "ReturnValue (nil)\n", depth+1)
	}
	debugWrite(w, "}\n", depth)
}

func debugWrite(w io.Writer, str string, depth int) {
	fmt.Fprintf(w, "%s%s", strings.Repeat(" ", depth*2), str)
}
