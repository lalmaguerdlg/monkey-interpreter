package evaluator

import (
	"monkey/object"
	"strconv"
)

func castObjectToString(obj object.Object) *object.String {
	return &object.String{Value: obj.Inspect()}
}

func castObjectToInteger(obj object.Object) object.Object {
	var result object.Integer
	switch obj := obj.(type) {
	case *object.Boolean:
		val := 0
		if obj == TRUE {
			val = 1
		}
		result.Value = int64(val)
	case *object.Integer:
		return obj
	case *object.String:
		val, err := strconv.ParseInt(obj.Value, 0, 64)
		if err != nil {
			return newError("cannot parse %q to int: invalid syntax", obj.Value)
		}
		result.Value = val
	default:
		return newError("cannot cast %s to %s: incompatible types", obj.Type(), object.IntegerType)
	}
	return &result
}

func castObjectToBoolean(obj object.Object) *object.Boolean {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj
	case *object.Integer:
		return castIntegerToBoolean(obj)
	case *object.Null:
		return FALSE
	default:
		return TRUE // Any other value is considered thruty
	}
}

func castIntegerToBoolean(value *object.Integer) *object.Boolean {
	if value.Value == 0 {
		return FALSE
	}
	return TRUE
}
