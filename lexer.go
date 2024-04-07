package yamlprocessor

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/file"
	"github.com/goccy/go-yaml/scanner"
	"github.com/goccy/go-yaml/token"
)

// ExprEntry is a type alias for a function. The function receives a Context and returns should return a func() that
// receives and return any number of params.
type ExprEntry func(ctx *Context) any

// Scanner is a scanner for YAML. It is a wrapper around the scanner.Scanner of the goccy/go-yaml library adding the
// ability to scan expressions inside the YAML string values.
//
// This Scanner uses the Processor LexerFunctions as the expr.Run env.
type Scanner struct {
	processor *Processor
	s         scanner.Scanner
	line      int
}

func (s *Scanner) Init(src string) {
	s.s.Init(src)
	s.line = 0
}

// Scan scans the YAML string and returns the tokens found in the string, using the goccy/go-yaml Scanner. If an
// expression (${...}) is found in the string, it will be compiled and executed (using the expr-lang/expr library) using
// the Processor LexerFunctions.
func (s *Scanner) Scan() (token.Tokens, error) {
	ctx := Context{
		Processor: s.processor,
	}
	env := make(map[string]any, len(s.processor.Env))
	for k, v := range s.processor.Env {
		env[k] = v(&ctx)
	}
	var tokens token.Tokens
	subTokens, err := s.s.Scan()
	if errors.Is(err, io.EOF) {
		return nil, err
	}
loopSubTokens:
	for _, n := range subTokens {
		switch n.Type {
		case token.StringType, token.SingleQuoteType, token.DoubleQuoteType:
			break
		default:
			tokens.Add(n)
			s.line += strings.Count(n.Origin, "\n")
			continue
		}

		var (
			before string
			found  bool
		)
		remain := n.Value

		for {
			before, remain, found = strings.Cut(n.Value, "${")
			if !found { // No start of expression found.
				s.line += strings.Count(n.Origin, "\n")
				tokens.Add(n)
				break
			}
			s.line += strings.Count(before, "\n")

			var expression string
			for {
				expression, remain, found = findCloseExpression(remain)
				if !found {
					break
				}
				s.line += strings.Count(expression, "\n")

				c, err := expr.Compile(expression)
				if err != nil {
					var f *file.Error
					if errors.As(err, &f) {
						f.Location.Line += s.line
						return nil, f
					}
					return nil, err
				}

				r, err := expr.Run(c, env)
				if err != nil {
					var f *file.Error
					if errors.As(err, &f) {
						f.Location.Line += s.line
						return nil, f
					}
					return nil, err
				}

				switch v := r.(type) {
				case token.Tokens:
					for _, n := range v {
						s.line += strings.Count(n.Origin, "\n")
					}
					tokens.Add(v...)
					continue loopSubTokens
				default:
					n.Value = before + fmt.Sprint(r) + remain
				}
				break
			}
		}
	}
	return tokens, nil
}

func findCloseExpression(value string) (string, string, bool) {
	for i := 0; i < len(value); i++ {
		switch value[i] {
		case '"':
		FindDQuote:
			for i++; i < len(value); i++ {
				switch value[i] {
				case '\\':
					i++
				case '"':
					break FindDQuote
				}
			}
		case '\'':
		FindSQuote:
			for i++; i < len(value); i++ {
				switch value[i] {
				case '\\':
					i++
				case '\'':
					break FindSQuote
				}
			}
		case '}':
			return value[0:i], value[i+1:], true
		}
	}
	return "", "", false
}

// Tokenize split to token instances from string
func (p *Processor) Tokenize(src string) (token.Tokens, error) {
	s := Scanner{
		processor: p,
	}
	s.Init(src)
	var tokens token.Tokens
	for {
		subTokens, err := s.Scan()
		switch {
		case errors.Is(err, io.EOF):
			if len(subTokens) > 0 {
				tokens.Add(subTokens...)
			}
			return tokens, nil
		case err != nil:
			return nil, err
		}
		tokens.Add(subTokens...)
	}
}
