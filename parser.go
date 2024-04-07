package yamlprocessor

import (
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

type expressionResolverVisitor struct {
	processor *Processor
}

func (v *expressionResolverVisitor) Visit(node ast.Node) ast.Visitor {
	//switch n := node.(type) {
	//case *ast.StringNode:
	//	sb := strings.Builder{}
	//	for i := 0; i < len(n.Value); i++ {
	//		if n.Value[i] != '$' {
	//			sb.WriteByte(n.Value[i])
	//			continue
	//		}
	//		i++
	//		if i < len(n.Value) && n.Value[i] == '{' {
	//			i++
	//			started := i
	//			for i < len(n.Value) && n.Value[i] != '}' {
	//				i++
	//			}
	//			c, err := expr.Compile(n.Value[started:i])
	//			if err != nil {
	//				// TODO(j): How to report this error?
	//				fmt.Println("Compile ERR>>>", err)
	//				return nil // Does interrupt??
	//			}
	//			r, err := expr.Run(c, v.processor.ExprEnv)
	//			if err != nil {
	//				// TODO(j): How to report this error?
	//				fmt.Println("RUN ERR>>>", err)
	//				return nil // Does interrupt??
	//			}
	//			fmt.Println(r)
	//		}
	//	}
	//}
	return v
}

func (p *Processor) ParseBytes(data []byte, mode parser.Mode) (*ast.File, error) {
	tokens, err := p.Tokenize(string(data))
	if err != nil {
		return nil, err
	}
	parsed, err := parser.Parse(tokens, mode)
	if err != nil {
		return nil, err
	}
	visitor := &expressionResolverVisitor{p}
	for _, doc := range parsed.Docs {
		ast.Walk(visitor, doc)
	}
	return parsed, nil
}
