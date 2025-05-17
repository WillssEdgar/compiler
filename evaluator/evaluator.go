package evaluator

import (
	"fmt"

	"compiler/ast"
	"compiler/environment"
)

func Eval(node ast.Node, env *environment.Environment) environment.Object {
	fmt.Printf("node: %v", node)
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)

	case *ast.LetStatement:
		val := Eval(node.Assignment.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Assignment.Name.Value, val)
		return val

	case *ast.AssignmentStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val

	case *ast.IntegerLiteral:
		return &environment.Integer{Value: node.Value}

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

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionalLiteral:
		return &environment.Function{Literal: node, Env: env}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		args := []environment.Object{}
		for _, a := range node.Arguments {
			args = append(args, Eval(a, env))
		}
		return applyFunction(function, args...)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		return &environment.ReturnValue{Value: val}
	}

	print("yeah")
	return nil
}

func applyFunction(fn environment.Object, args ...environment.Object) environment.Object {
	function, ok := fn.(*environment.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendedEnv := environment.NewEnclosedEnvironment(function.Env)

	for i, param := range function.Literal.Parameters {
		extendedEnv.Set(param.Value, args[i])
	}

	evaluated := Eval(function.Literal.Body, extendedEnv)

	if returnValue, ok := evaluated.(*environment.ReturnValue); ok {
		return returnValue.Value
	}
	return evaluated
}

func evalProgram(program *ast.Program, env *environment.Environment) environment.Object {
	var result environment.Object
	for _, stmt := range program.Statements {
		result = Eval(stmt, env)
	}
	return result
}

func evalPrefixExpression(operator string, right environment.Object) environment.Object {
	switch operator {
	case "-":
		if right.Type() != environment.INTEGER_OBJ {
			return newError("unknown operator: -%s", right.Type())
		}
		val := right.(*environment.Integer).Value
		return &environment.Integer{Value: -val}
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left, right environment.Object) environment.Object {
	if left.Type() == environment.INTEGER_OBJ && right.Type() == environment.INTEGER_OBJ {
		return evalIntegerInfixExpression(operator, left, right)
	}
	return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
}

func evalIntegerInfixExpression(operator string, left, right environment.Object) environment.Object {
	l := left.(*environment.Integer).Value
	r := right.(*environment.Integer).Value

	switch operator {
	case "+":
		return &environment.Integer{Value: l + r}
	case "-":
		return &environment.Integer{Value: l - r}
	case "*":
		return &environment.Integer{Value: l * r}
	case "/":
		return &environment.Integer{Value: l / r}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIdentifier(node *ast.Identifier, env *environment.Environment) environment.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	return newError("identifier not found: %s", node.Value)
}

func newError(format string, a ...interface{}) *environment.Error {
	return &environment.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj environment.Object) bool {
	if obj != nil {
		return obj.Type() == environment.ERROR_OBJ
	}
	return false
}
