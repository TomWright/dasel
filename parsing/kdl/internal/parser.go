package internal

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Parser parses KDL tokens into a Document.
type Parser struct {
	tok *Tokenizer
}

// Parse parses the given KDL input and returns a Document.
func Parse(input string) (*Document, error) {
	tok := NewTokenizer(input)

	// Check for version marker: /- kdl-version N
	p := &Parser{tok: tok}
	if err := p.detectVersionMarker(input); err != nil {
		return nil, err
	}

	return p.parseDocument()
}

func (p *Parser) detectVersionMarker(input string) error {
	trimmed := strings.TrimSpace(input)
	if strings.HasPrefix(trimmed, "/- kdl-version ") {
		rest := strings.TrimPrefix(trimmed, "/- kdl-version ")
		if len(rest) > 0 {
			// Extract the version number (up to next whitespace or newline)
			verStr := rest
			if idx := strings.IndexAny(verStr, " \t\n\r"); idx >= 0 {
				verStr = verStr[:idx]
			}
			switch verStr {
			case "1":
				p.tok.Version = Version1
			case "2":
				p.tok.Version = Version2
			default:
				return fmt.Errorf("kdl: unsupported version %q in version marker", verStr)
			}
		}
	}
	return nil
}

func (p *Parser) parseDocument() (*Document, error) {
	doc := &Document{}

	for {
		p.skipNewlines()

		tok, err := p.tok.PeekToken()
		if err != nil {
			return nil, err
		}

		if tok.Type == TokenEOF {
			break
		}

		if tok.Type == TokenCloseBrace {
			break
		}

		// Check for slashdash before node
		if tok.Type == TokenSlashDash {
			if _, err := p.tok.NextToken(); err != nil {
				return nil, err
			}
			// Skip the next node
			_, err := p.parseNode()
			if err != nil {
				return nil, err
			}
			continue
		}

		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if node != nil {
			doc.Nodes = append(doc.Nodes, node)
		}
	}

	return doc, nil
}

func (p *Parser) parseNode() (*Node, error) {
	p.skipNewlines()

	node := &Node{}

	// Optional type annotation
	tok, err := p.tok.PeekToken()
	if err != nil {
		return nil, err
	}
	if tok.Type == TokenOpenParen {
		typeAnn, err := p.parseTypeAnnotation()
		if err != nil {
			return nil, err
		}
		node.Type = typeAnn
	}

	// Node name
	nameTok, err := p.tok.NextToken()
	if err != nil {
		return nil, err
	}

	switch nameTok.Type {
	case TokenIdentifier, TokenQuotedString, TokenRawString:
		node.Name = nameTok.Value
	default:
		return nil, fmt.Errorf("kdl: expected node name, got %v at line %d, col %d", nameTok.Type, nameTok.Line, nameTok.Col)
	}

	// Parse entries (arguments and properties) and optional children
	if err := p.parseNodeEntries(node); err != nil {
		return nil, err
	}

	return node, nil
}

func (p *Parser) parseNodeEntries(node *Node) error {
	for {
		tok, err := p.tok.PeekToken()
		if err != nil {
			return err
		}

		// Node terminators
		if tok.Type == TokenNewline || tok.Type == TokenSemicolon || tok.Type == TokenEOF || tok.Type == TokenCloseBrace {
			if tok.Type == TokenSemicolon {
				if _, err := p.tok.NextToken(); err != nil {
					return err
				}
			}
			return nil
		}

		// Children block
		if tok.Type == TokenOpenBrace {
			if _, err := p.tok.NextToken(); err != nil {
				return err
			}
			children, err := p.parseDocument()
			if err != nil {
				return err
			}
			node.Children = children.Nodes

			// Expect closing brace
			p.skipNewlines()
			closeTok, err := p.tok.NextToken()
			if err != nil {
				return err
			}
			if closeTok.Type != TokenCloseBrace {
				return fmt.Errorf("kdl: expected '}', got %v at line %d, col %d", closeTok.Type, closeTok.Line, closeTok.Col)
			}

			// After children block, node must end
			return p.expectNodeEnd()
		}

		// Slashdash before entry
		if tok.Type == TokenSlashDash {
			if _, err := p.tok.NextToken(); err != nil {
				return err
			}

			// Peek what follows to know if we're slashdashing an arg, prop, or children
			next, err := p.tok.PeekToken()
			if err != nil {
				return err
			}

			if next.Type == TokenOpenBrace {
				// Slashdash children
				if _, err := p.tok.NextToken(); err != nil {
					return err
				}
				if _, err := p.parseDocument(); err != nil {
					return err
				}
				p.skipNewlines()
				closeTok, err := p.tok.NextToken()
				if err != nil {
					return err
				}
				if closeTok.Type != TokenCloseBrace {
					return fmt.Errorf("kdl: expected '}' after slashdashed children")
				}
				continue
			}

			// Slashdash an entry (arg or prop) — parse and discard
			if err := p.skipEntry(); err != nil {
				return err
			}
			continue
		}

		// Try to parse as property or argument
		if err := p.parseEntry(node); err != nil {
			return err
		}
	}
}

func (p *Parser) parseEntry(node *Node) error {
	// Check for type annotation
	tok, err := p.tok.PeekToken()
	if err != nil {
		return err
	}

	typeAnn := ""
	if tok.Type == TokenOpenParen {
		typeAnn, err = p.parseTypeAnnotation()
		if err != nil {
			return err
		}
	}

	// Read the value/name token
	valTok, err := p.tok.NextToken()
	if err != nil {
		return err
	}

	// Check if this is a property (identifier/string followed by =)
	if (valTok.Type == TokenIdentifier || valTok.Type == TokenQuotedString || valTok.Type == TokenRawString) && typeAnn == "" {
		next, err := p.tok.PeekToken()
		if err != nil {
			return err
		}
		if next.Type == TokenEquals {
			if _, err := p.tok.NextToken(); err != nil {
				return err
			}
			val, err := p.parseValue()
			if err != nil {
				return err
			}
			node.Properties = append(node.Properties, &Property{
				Key:   valTok.Value,
				Value: val,
			})
			return nil
		}
	}

	// It's an argument
	val, err := tokenToValue(valTok, typeAnn)
	if err != nil {
		return err
	}
	node.Arguments = append(node.Arguments, val)
	return nil
}

func (p *Parser) skipEntry() error {
	// Check for type annotation
	tok, err := p.tok.PeekToken()
	if err != nil {
		return err
	}

	if tok.Type == TokenOpenParen {
		if _, err := p.parseTypeAnnotation(); err != nil {
			return err
		}
	}

	// Read the value/name token
	valTok, err := p.tok.NextToken()
	if err != nil {
		return err
	}

	// Check for property (name = value)
	if valTok.Type == TokenIdentifier || valTok.Type == TokenQuotedString || valTok.Type == TokenRawString {
		next, err := p.tok.PeekToken()
		if err != nil {
			return err
		}
		if next.Type == TokenEquals {
			if _, err := p.tok.NextToken(); err != nil {
				return err
			}
			if _, err := p.parseValue(); err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func (p *Parser) parseValue() (*Value, error) {
	// Check for type annotation
	tok, err := p.tok.PeekToken()
	if err != nil {
		return nil, err
	}

	typeAnn := ""
	if tok.Type == TokenOpenParen {
		typeAnn, err = p.parseTypeAnnotation()
		if err != nil {
			return nil, err
		}
	}

	valTok, err := p.tok.NextToken()
	if err != nil {
		return nil, err
	}

	return tokenToValue(valTok, typeAnn)
}

func (p *Parser) parseTypeAnnotation() (string, error) {
	// Consume (
	if _, err := p.tok.NextToken(); err != nil {
		return "", err
	}

	// Read the type name
	typeTok, err := p.tok.NextToken()
	if err != nil {
		return "", err
	}
	if typeTok.Type != TokenIdentifier && typeTok.Type != TokenQuotedString {
		return "", fmt.Errorf("kdl: expected type name, got %v at line %d, col %d", typeTok.Type, typeTok.Line, typeTok.Col)
	}

	// Consume )
	closeTok, err := p.tok.NextToken()
	if err != nil {
		return "", err
	}
	if closeTok.Type != TokenCloseParen {
		return "", fmt.Errorf("kdl: expected ')', got %v at line %d, col %d", closeTok.Type, closeTok.Line, closeTok.Col)
	}

	return typeTok.Value, nil
}

func (p *Parser) expectNodeEnd() error {
	tok, err := p.tok.PeekToken()
	if err != nil {
		return err
	}
	if tok.Type == TokenNewline || tok.Type == TokenSemicolon || tok.Type == TokenEOF || tok.Type == TokenCloseBrace {
		if tok.Type == TokenSemicolon {
			if _, err := p.tok.NextToken(); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("kdl: expected node terminator, got %v at line %d, col %d", tok.Type, tok.Line, tok.Col)
}

func (p *Parser) skipNewlines() {
	for {
		tok, err := p.tok.PeekToken()
		if err != nil || tok.Type != TokenNewline {
			return
		}
		_, _ = p.tok.NextToken()
	}
}

func tokenToValue(tok Token, typeAnn string) (*Value, error) {
	v := &Value{Type: typeAnn}

	switch tok.Type {
	case TokenQuotedString, TokenRawString, TokenMultiLineString:
		v.Value = tok.Value
	case TokenInteger:
		n, err := strconv.ParseInt(tok.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("kdl: invalid integer %q: %w", tok.Value, err)
		}
		v.Value = n
	case TokenFloat:
		n, err := strconv.ParseFloat(tok.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("kdl: invalid float %q: %w", tok.Value, err)
		}
		v.Value = n
	case TokenHexInt:
		s := tok.Value
		neg := false
		if strings.HasPrefix(s, "-") {
			neg = true
			s = s[1:]
		} else if strings.HasPrefix(s, "+") {
			s = s[1:]
		}
		s = strings.TrimPrefix(s, "0x")
		s = strings.TrimPrefix(s, "0X")
		n, err := strconv.ParseInt(s, 16, 64)
		if err != nil {
			return nil, fmt.Errorf("kdl: invalid hex integer %q: %w", tok.Value, err)
		}
		if neg {
			n = -n
		}
		v.Value = n
	case TokenOctalInt:
		s := tok.Value
		neg := false
		if strings.HasPrefix(s, "-") {
			neg = true
			s = s[1:]
		} else if strings.HasPrefix(s, "+") {
			s = s[1:]
		}
		s = strings.TrimPrefix(s, "0o")
		s = strings.TrimPrefix(s, "0O")
		n, err := strconv.ParseInt(s, 8, 64)
		if err != nil {
			return nil, fmt.Errorf("kdl: invalid octal integer %q: %w", tok.Value, err)
		}
		if neg {
			n = -n
		}
		v.Value = n
	case TokenBinaryInt:
		s := tok.Value
		neg := false
		if strings.HasPrefix(s, "-") {
			neg = true
			s = s[1:]
		} else if strings.HasPrefix(s, "+") {
			s = s[1:]
		}
		s = strings.TrimPrefix(s, "0b")
		s = strings.TrimPrefix(s, "0B")
		n, err := strconv.ParseInt(s, 2, 64)
		if err != nil {
			return nil, fmt.Errorf("kdl: invalid binary integer %q: %w", tok.Value, err)
		}
		if neg {
			n = -n
		}
		v.Value = n
	case TokenTrue:
		v.Value = true
	case TokenFalse:
		v.Value = false
	case TokenNull:
		v.Value = nil
	case TokenInf:
		v.Value = math.Inf(1)
	case TokenNegInf:
		v.Value = math.Inf(-1)
	case TokenNaN:
		v.Value = math.NaN()
	case TokenIdentifier:
		// Bare identifier used as a value — treat as string
		v.Value = tok.Value
	default:
		return nil, fmt.Errorf("kdl: unexpected token %v for value at line %d, col %d", tok.Type, tok.Line, tok.Col)
	}

	return v, nil
}
