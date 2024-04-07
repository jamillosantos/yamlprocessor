package yamlprocessor

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
)

// Processor is a struct that holds lexer and parser functions.
type Processor struct {
	Env ExprEnv
}

// NewProcessor returns a new Processor.
func NewProcessor() Processor {
	return Processor{
		Env: make(ExprEnv),
	}
}

// Unmarshal will unmarshal the given data into the given interface. It will use the ParseBytes method to parse the data
// with the expressions and then unmarshal the parsed data into the given interface.
//
// This methods resource wasteful as it will parse the data twice, once to resolve the expressions (In the future
// releaseWe could try to parse the data directly from the AST built from the ParseBytes).
func (p *Processor) Unmarshal(data []byte, v interface{}) error {
	doc, err := p.ParseBytes(data, parser.ParseComments)
	if err != nil {
		return err
	}
	buf2 := bytes.Buffer{}
	for i, n := range doc.Docs {
		if i > 0 {
			_, err := fmt.Fprintln(&buf2, "---")
			if err != nil {
				return err
			}
		}
		_, err := fmt.Fprintln(&buf2, n.String())
		if err != nil {
			return err
		}
	}
	return yaml.Unmarshal(buf2.Bytes(), v)
}
