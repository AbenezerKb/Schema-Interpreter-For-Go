package main

import (
	"fmt"
	"os"
	"unicode"
)

type lexingContext struct {
	source         []rune
	sourceFileName string
}

type tokenKind uint

const (
	syntaxToken tokenKind = iota
	integerToken
	identifierToken
)

type token struct {
	value    string
	kind     tokenKind
	location int
	lc       lexingContext
}

func (t token) debug(description string) {
	var tokenLine []rune
	var tokenLineNumber int
	var tokenColumn int
	var inTokenLine bool
	var i int

	for i < len(t.lc.source) {
		r := t.lc.source[i]

		if i < t.location {
			tokenColumn++
		}

		tokenLine = append(tokenLine, r)

		if r == '\n' {
			tokenLineNumber++
			if inTokenLine {
				break
			}

			tokenColumn = 1
			tokenLine = nil
		}

		if i == t.location {
			inTokenLine = true
		}

		i++
	}

	fmt.Printf("%s [at line %d, column %d in file %s]\n",
		description, tokenLineNumber, tokenColumn, t.lc.sourceFileName)
	fmt.Println(string(tokenLine))

	for tokenColumn >= 1 {
		fmt.Printf(" ")
		tokenColumn--
	}
	fmt.Println("^ near here")
}

func eatWhitespace(source []rune, cursor int) int {
	for cursor < len(source) {
		if unicode.IsSpace(source[cursor]) {
			cursor++
			continue
		}

		break
	}

	return cursor
}

func (lc lexingContext) lexSyntaxToken(cursor int) (int, *token) {
	if lc.source[cursor] == '(' || lc.source[cursor] == ')' {
		return cursor + 1, &token{
			value:    string([]rune{lc.source[cursor]}),
			kind:     syntaxToken,
			location: cursor,
			lc:       lc,
		}
	}

	return cursor, nil
}

func (lc lexingContext) lexIntegerToken(cursor int) (int, *token) {
	originalCursor := cursor

	var value []rune
	for cursor < len(lc.source) {
		r := lc.source[cursor]
		if r >= '0' && r <= '9' {
			value = append(value, r)
			cursor++
			continue
		}

		break
	}

	if len(value) == 0 {
		return originalCursor, nil
	}

	return cursor, &token{
		value:    string(value),
		kind:     integerToken,
		location: originalCursor,
		lc:       lc,
	}
}

func (lc lexingContext) lexIdentifierToken(cursor int) (int, *token) {
	originalCursor := cursor
	var value []rune

	for cursor < len(lc.source) {
		r := lc.source[cursor]
		if !(unicode.IsSpace(r) || r == ')') {
			value = append(value, r)
			cursor++
			continue
		}

		break
	}

	if len(value) == 0 {
		return originalCursor, nil
	}

	return cursor, &token{
		value:    string(value),
		kind:     identifierToken,
		location: originalCursor,
		lc:       lc,
	}
}

func (lc lexingContext) lex() []token {
	var tokens []token
	var t *token

	cursor := 0
	for cursor < len(lc.source) {
		cursor = eatWhitespace(lc.source, cursor)
		if cursor == len(lc.source) {
			break
		}

		cursor, t = lc.lexSyntaxToken(cursor)
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}

		cursor, t = lc.lexIntegerToken(cursor)
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}

		cursor, t = lc.lexIdentifierToken(cursor)
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}

		panic("Could not lex")
	}

	return tokens
}

func newLexingContext(file string) lexingContext {
	program, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	return lexingContext{
		sourceFileName: file,
		source:         []rune(string(program)),
	}
}
