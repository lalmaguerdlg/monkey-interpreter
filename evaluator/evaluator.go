// Package evaluator
package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
	"strings"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.AssignmentStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		assigned := env.Assign(node.Name.Value, val)
		if isError(assigned) {
			return assigned
		}
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

		// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.DotExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		return evalDotExpression(left, node.Right)
	}
	return NULL
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	if result == nil {
		return NULL
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			switch result.Type() {
			case object.ErrorType:
				fallthrough
			case object.ReturnType:
				return result
			}
		}
	}

	if result == nil {
		return NULL
	}

	return result
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, exp := range expressions {
		evaluated := Eval(exp, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, arguments []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(arguments) < len(fn.Parameters) {
			missing := []string{}
			for _, param := range fn.Parameters[len(arguments):] {
				missing = append(missing, param.Value)
			}
			return newError("function call is missing parameters: %s", strings.Join(missing, ", "))
		}
		extendedEnv := extendFunctionEnv(fn, arguments)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(arguments...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, arguments []object.Object) *object.Environment {
	env := object.NewEnclousedEnvironment(fn.Env)
	for i, param := range fn.Parameters {
		env.Set(param.Value, arguments[i])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	value, ok := obj.(*object.ReturnValue)
	if !ok {
		return obj
	}
	return value.Value
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashable, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashable.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier %s is undefined", node.Value)
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerType && right.Type() == object.IntegerType:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.StringType && right.Type() == object.StringType:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.ArrayType && index.Type() == object.IntegerType:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.StringType && index.Type() == object.IntegerType:
		return evalStringIndexExpression(left, index)
	case left.Type() == object.HashType:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalDotExpression(left object.Object, right *ast.Identifier) object.Object {
	switch {
	case left.Type() == object.HashType:
		return evalHashDotExpression(left, right)
	default:
		return newError("dot operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array object.Object, index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObj.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObj.Elements[idx]
}

func evalStringIndexExpression(str object.Object, index object.Object) object.Object {
	strObj := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObj.Value) - 1)

	if idx < 0 || idx > max {
		return NULL
	}
	return &object.String{Value: string(strObj.Value[idx])}
}

func evalHashIndexExpression(obj object.Object, index object.Object) object.Object {
	hash := obj.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unable to use hash key: %s", index.Type())
	}

	pair, ok := hash.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

func evalHashDotExpression(left object.Object, right *ast.Identifier) object.Object {
	hash := left.(*object.Hash)

	key := &object.String{Value: right.Value}

	pair, ok := hash.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}
	casted := castObjectToBoolean(condition)
	if casted == TRUE {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	}

	return NULL
}

func evalBangOperatorExpression(operand object.Object) object.Object {
	switch operand {
	case FALSE:
		return TRUE
	case TRUE:
		return FALSE
	default:
		casted := castObjectToBoolean(operand)
		return evalBangOperatorExpression(casted)
	}
}

func evalMinusPrefixOperatorExpression(operand object.Object) object.Object {
	if operand.Type() != object.IntegerType {
		return newError("unknown operator: -%s", operand.Type())
	}

	value := operand.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func nativeBoolToBooleanObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ErrorType
}
