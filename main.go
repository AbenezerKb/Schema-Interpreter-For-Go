package main

import "os"
import "fmt"

func main() {
	lc := newLexingContext(os.Args[1])

	tokens := lc.lex()
	debug := false
	if debug {
		for _, token := range tokens {
			fmt.Println(token.value)
		}
	}

	var parseIndex int
	var a = ast{
		value{
			kind: literalValue,
			literal: &token{
				value: "begin",
				kind:  identifierToken,
			},
		},
	}

	for parseIndex < len(tokens) {
		childAst, nextIndex := parse(tokens, parseIndex)
		a = append(a, value{
			kind: listValue,
			list: &childAst,
		})
		parseIndex = nextIndex
	}

	if parseIndex < len(tokens) {
		panic("Incomplete parse")
	}

	if debug {
		fmt.Println(a.pretty())
	}

	initializeBuiltins()
	ctx := map[string]any{}
	value := astWalk(a, ctx)
	fmt.Println(value)
}
