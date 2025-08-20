// Package object
package object

import (
	"bytes"
	"hash/fnv"
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
	HashType     ObjectType = "HASH"
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

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return IntegerType }
func (i *Integer) Inspect() string  { return strconv.FormatInt(i.Value, 10) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BooleanType }
func (b *Boolean) Inspect() string  { return strconv.FormatBool(b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return StringType }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	key := HashKey{Type: s.Type(), Value: h.Sum64()}
	return key
}

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

type HashPair struct {
	Key   Object
	Value Object
}
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (hm *Hash) Type() ObjectType { return HashType }
func (hm *Hash) Inspect() string {
	var out strings.Builder
	pairs := []string{}
	for _, pair := range hm.Pairs {
		pairs = append(pairs, pair.Key.Inspect()+":"+pair.Value.Inspect())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
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
