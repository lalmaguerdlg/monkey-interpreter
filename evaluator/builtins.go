package evaluator

import (
	"fmt"
	"monkey/object"
	"strings"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=%d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", arg.Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=%d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.String:
				if len(arg.Value) == 0 {
					return &object.String{Value: ""}
				}
				return &object.String{Value: string(arg.Value[0])}
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[0]
			default:
				return newError("argument to `first` not supported, got %s", arg.Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=%d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.String:
				length := len(arg.Value)
				if length == 0 {
					return &object.String{Value: ""}
				}
				return &object.String{Value: string(arg.Value[length-1])}
			case *object.Array:
				length := len(arg.Elements)
				if length == 0 {
					return NULL
				}
				return arg.Elements[length-1]
			default:
				return newError("argument to `last` not supported, got %s", arg.Type())
			}
		},
	},
	"tail": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=%d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.String:
				length := len(arg.Value)
				if length == 0 {
					return &object.String{Value: ""}
				}

				return &object.String{Value: strings.Clone(arg.Value[1:length])}
			case *object.Array:
				length := len(arg.Elements)
				if length == 0 {
					return NULL
				}
				newElements := make([]object.Object, length-1)
				copy(newElements, arg.Elements[1:length])
				return &object.Array{Elements: newElements}
			default:
				return newError("argument to `tail` not supported, got %s", arg.Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments, got=%d, want=%d", len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				newElements := make([]object.Object, length+1)
				copy(newElements, arg.Elements)
				newElements[length] = args[1]
				return &object.Array{Elements: newElements}
			default:
				return newError("argument to `push` not supported, got %s", arg.Type())
			}
		},
	},
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Printf("%s\n", arg.Inspect())
			}
			return NULL
		},
	},

	// Type conversion:
	"string": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=%d", len(args), 1)
			}
			return castObjectToString(args[0])
		},
	},
	"int": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=%d", len(args), 1)
			}
			return castObjectToInteger(args[0])
		},
	},
	"bool": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, got=%d, want=%d", len(args), 1)
			}
			return castObjectToBoolean(args[0])
		},
	},
}
