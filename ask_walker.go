package main

import (
	"fmt"
	"os"
	"strconv"
)

var builtins = map[string]func([]value, map[string]any) any{}

func copyContext(in map[string]any) map[string]any {
	out := map[string]any{}
	for key, val := range in {
		out[key] = val
	}

	return out
}

func initializeBuiltins() {
	builtins["if"] = func(args []value, ctx map[string]any) any {
		condition := astWalk2(args[0], ctx)
		then := args[1]
		_else := args[2]

		if condition.(bool) == true {
			return astWalk2(then, ctx)
		}

		return astWalk2(_else, ctx)
	}

	builtins["<"] = func(args []value, ctx map[string]any) any {
		return astWalk2(args[0], ctx).(int64) < astWalk2(args[1], ctx).(int64)
	}

	builtins["+"] = func(args []value, ctx map[string]any) any {
		var i int64
		for _, arg := range args {
			i += astWalk2(arg, ctx).(int64)
		}
		return i
	}

	builtins["-"] = func(args []value, ctx map[string]any) any {
		i := astWalk2(args[0], ctx).(int64)
		for _, arg := range args[1:] {
			i -= astWalk2(arg, ctx).(int64)
		}
		return i
	}

	builtins["begin"] = func(args []value, ctx map[string]any) any {
		var last any
		for _, arg := range args {
			last = astWalk2(arg, ctx)
		}

		return last
	}

	builtins["func"] = func(args []value, ctx map[string]any) any {
		functionName := (*args[0].literal).value

		params := *args[1].list

		body := *args[2].list

		ctx[functionName] = func(args []any, ctx map[string]any) any {
			childCtx := copyContext(ctx)
			if len(params) != len(args) {
				panic(fmt.Sprintf("Expected %d args to `%s`, got %d", len(params), functionName, len(args)))
			}
			for i, param := range params {
				childCtx[(*param.literal).value] = args[i]
			}

			return astWalk(body, childCtx)
		}

		return ctx[functionName]
	}
}

func astWalk(ast []value, ctx map[string]any) any {
	functionName := (*ast[0].literal).value

	if builtinFunction, ok := builtins[functionName]; ok {
		return builtinFunction(ast[1:], ctx)
	}

	maybeFunction, ok := ctx[functionName]
	if !ok {
		(*ast[0].literal).debug(fmt.Sprintf("Expected function, got %s", functionName))
		os.Exit(1)
	}
	userDefinedFunction := maybeFunction.(func([]any, map[string]any) any)

	var args []any
	for _, unevaluatedArg := range ast[1:] {
		args = append(args, astWalk2(unevaluatedArg, ctx))
	}

	return userDefinedFunction(args, ctx)
}

func astWalk2(v value, ctx map[string]any) any {
	if v.kind == literalValue {
		t := *v.literal
		switch t.kind {
		case integerToken:
			i, err := strconv.ParseInt(t.value, 10, 64)
			if err != nil {
				fmt.Println("Expected an integer, got: ", t.value)
				panic(err)
			}

			return i
		case identifierToken:
			return ctx[t.value]
		}
	}

	return astWalk(*v.list, ctx)
}
