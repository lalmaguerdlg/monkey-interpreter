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
	IntegerType  ObjectType = "INTEGER"
	BooleanType  ObjectType = "BOOLEAN"
	StringType   ObjectType = "STRING"
	ArrayType    ObjectType = "ARRAY"
	ReturnType   ObjectType = "RETURN"
	FunctionType ObjectType = "FUNCTION"
	BuiltinType  ObjectType = "BUILTIN"
	ErrorType    ObjectType = "ERROR"
	NullType     ObjectType = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return IntegerType }
func (i *Integer) Inspect() string  { return strconv.FormatInt(i.Value, 10) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BooleanType }
func (b *Boolean) Inspect() string  { return strconv.FormatBool(b.Value) }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return StringType }
func (s *String) Inspect() string  { return s.Value }

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ArrayType }
func (a *Array) Inspect() string {
	var out strings.Builder
	list := []string{}
	for _, el := range a.Elements {
		list = append(list, el.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(list, ", "))
	out.WriteString("]")
	return out.String()
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return ReturnType }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FunctionType }
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

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BuiltinType }
func (b *Builtin) Inspect() string  { return "builting function" }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ErrorType }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Null struct{}

func (i *Null) Type() ObjectType { return NullType }
func (i *Null) Inspect() string  { return "null" }
