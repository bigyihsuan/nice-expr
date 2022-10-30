package parser

import (
	"fmt"
	"nice-expr/lexer/token"
)

type ParseError struct {
	Reason   string
	Token    *token.Token
	LastRule string
}

func NewParseError(reason string, token *token.Token, lastRule string) *ParseError {
	err := new(ParseError)
	err.Reason = reason
	err.Token = token
	err.LastRule = lastRule
	return err
}

func (pe *ParseError) addRule(rule string) *ParseError {
	pe.LastRule += fmt.Sprintf(" %s", rule)
	return pe
}

// for interface error.Error()
func (pe *ParseError) Error() string {
	return fmt.Sprintf("%v %v because of token `%v` within rule trace `%v`", "PARSE ERROR:", pe.Reason, pe.Token, pe.LastRule)
}
