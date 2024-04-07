package yamlprocessor

import (
	"io"
	"os"

	"github.com/goccy/go-yaml/token"
)

// File will load and paste the content of the YAML file in the place where it is called. The file loaded will be
// processed by the same processor that is calling this function. Also, all the expressions will be evaluated.
func File(ctx *Context) any {
	return func(file string) (token.Tokens, error) {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = f.Close()
		}()
		content, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		return ctx.Processor.Tokenize(string(content))
	}
}

// Env will return the value of the environment variable with the given name.
func Env(ctx *Context) any {
	return func(name string) string {
		return os.Getenv(name)
	}
}
