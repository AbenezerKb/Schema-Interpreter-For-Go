package main

import "os"

type valueKind uint

const (
	literalValue valueKind = iota
	listValue
)

type value struct {
	kind    valueKind
	literal *token
	list    *ast
}

func (v value) pretty() string {
	if v.kind == literalValue {
		return v.literal.value
	}

	return v.list.pretty()
}

type ast []value

func (a ast) pretty() string {
	p := "("
	for _, value := range a {
		p += value.pretty()
		p += " "
	}

	return p + ")"
}

func parse(tokens []token, index int) (ast, int) {
	var a ast

	token := tokens[index]
	if !(token.kind == syntaxToken &&
		token.value == "(") {
		panic("Should be an open parenthesis")
	}
	index++

	for index < len(tokens) {
		token := tokens[index]
		if token.kind == syntaxToken &&
			token.value == "(" {
			child, nextIndex := parse(tokens, index)
			a = append(a, value{
				kind: listValue,
				list: &child,
			})
			index = nextIndex
			continue
		}

		if token.kind == syntaxToken &&
			token.value == ")" {
			return a, index + 1
		}

		a = append(a, value{
			kind:    literalValue,
			literal: &token,
		})
		index++
	}

	if tokens[index-1].kind == syntaxToken &&
		tokens[index-1].value != ")" {
		tokens[index-1].debug("Expected closing paren")
		os.Exit(1)
	}

	return a, index
}
