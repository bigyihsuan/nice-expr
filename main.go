package main

import (
	"fmt"
	"nice-expr/lexer"
)

func main() {
	tok := lexer.Token{Tt: lexer.Identifier, Lexeme: "print", Line: 0, Start: 11, End: 222}
	fmt.Println(tok)
}
