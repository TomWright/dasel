package parser

import "github.com/tomwright/dasel/v3/selector/lexer"

// null denotation tokens are tokens that expect no token to the left of them.
var nullDenotationTokens = []lexer.TokenKind{}

// left denotation tokens are tokens that expect a token to the left of them.
var leftDenotationTokens = []lexer.TokenKind{
	lexer.Plus,
	lexer.Dash,
	lexer.Slash,
	lexer.Star,
	lexer.Percent,
	lexer.Equal,
	lexer.NotEqual,
	lexer.GreaterThan,
	lexer.GreaterThanOrEqual,
	lexer.LessThan,
	lexer.LessThanOrEqual,
}

// right denotation tokens are tokens that expect a token to the right of them.
var rightDenotationTokens = []lexer.TokenKind{
	lexer.Exclamation, // Not operator
}

type bindingPower int

const (
	bpDefault bindingPower = iota
	bpAssignment
	bpLogical
	bpRelational
	bpAdditive
	bpMultiplicative
	bpUnary
	bpCall
	bpProperty
	bpLiteral
)

var tokenBindingPowers = map[lexer.TokenKind]bindingPower{
	lexer.String: bpLiteral,
	lexer.Number: bpLiteral,
	lexer.Bool:   bpLiteral,
	//lexer.Null:             bpLiteral,

	lexer.Variable:    bpProperty,
	lexer.Dot:         bpProperty,
	lexer.OpenBracket: bpProperty,

	lexer.OpenParen: bpCall,

	lexer.Exclamation: bpUnary,

	lexer.Star:    bpMultiplicative,
	lexer.Slash:   bpMultiplicative,
	lexer.Percent: bpMultiplicative,

	lexer.Plus: bpAdditive,
	lexer.Dash: bpAdditive,

	lexer.Equal:              bpRelational,
	lexer.NotEqual:           bpRelational,
	lexer.GreaterThan:        bpRelational,
	lexer.GreaterThanOrEqual: bpRelational,
	lexer.LessThan:           bpRelational,
	lexer.LessThanOrEqual:    bpRelational,

	lexer.And:     bpLogical,
	lexer.Or:      bpLogical,
	lexer.Like:    bpLogical,
	lexer.NotLike: bpLogical,

	lexer.Equals: bpAssignment,
}

func getTokenBindingPower(t lexer.TokenKind) bindingPower {
	if bp, ok := tokenBindingPowers[t]; ok {
		return bp
	}
	return bpDefault
}
