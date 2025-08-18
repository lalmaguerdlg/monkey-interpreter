// Package object
package object

import (
	"strconv"
)

type ObjectType string

const (
	IntegerObj ObjectType = "INTEGER"
	BooleanObj ObjectType = "BOOLEAN"
	ReturnObj  ObjectType = "RETURN"
	NullObj    ObjectType = "NULL"
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

type Null struct{}

func (i *Null) Type() ObjectType { return NullObj }
func (i *Null) Inspect() string  { return "null" }
