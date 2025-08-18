// Package object
package object

import (
	"bytes"
	"monkey/ast"
	"strconv"
	"strings"
)

type ObjectType string

const (
	IntegerObj  ObjectType = "INTEGER"
	BooleanObj  ObjectType = "BOOLEAN"
	ReturnObj   ObjectType = "RETURN"
	FunctionObj ObjectType = "FUNCTION"
	ErrorObj    ObjectType = "ERROR"
	NullObj     ObjectType = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return IntegerObj }
func (i *Integer) Inspect() string  { return strconv.FormatInt(i.Value, 10) }

type Boolean struct {
	Value bool
}

func (i *Boolean) Type() ObjectType { return BooleanObj }
func (i *Boolean) Inspect() string  { return strconv.FormatBool(i.Value) }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return ReturnObj }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FunctionObj }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ErrorObj }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Null struct{}

func (i *Null) Type() ObjectType { return NullObj }
func (i *Null) Inspect() string  { return "null" }
