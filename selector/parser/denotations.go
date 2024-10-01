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
