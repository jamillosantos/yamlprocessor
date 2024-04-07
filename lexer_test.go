package yamlprocessor

import (
	"io"
	"testing"

	"github.com/expr-lang/expr/file"
	"github.com/goccy/go-yaml/token"
	"github.com/stretchr/testify/require"
)

func TestScanner_Scan(t *testing.T) {
	t.Run("should tokenize ignore non lexer functions", func(t *testing.T) {
		processor := NewProcessor()
		processor.Env.Add("f", func(ctx *Context) any {
			return func() string {
				return "value"
			}
		})
		scanner := Scanner{
			processor: &processor,
		}
		scanner.Init("d: ${f()}")
		tokens, err := scanner.Scan()
		require.NoError(t, err)
		require.Len(t, tokens, 2)
		require.Equal(t, token.StringType, tokens[0].Type)
		require.Equal(t, "d", tokens[0].Value)
		require.Equal(t, token.MappingValueType, tokens[1].Type)
		require.Equal(t, ":", tokens[1].Value)
		tokens, err = scanner.Scan()
		require.NoError(t, err)
		require.Len(t, tokens, 1)
		require.Equal(t, token.StringType, tokens[0].Type)
		require.Equal(t, "value", tokens[0].Value)
		_, err = scanner.Scan()
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("should return a compilation from express with right line", func(t *testing.T) {
		processor := NewProcessor()
		scanner := Scanner{
			processor: &processor,
		}
		scanner.Init(`a: 1
b: 2
c: ${123 +}`)
		for {
			_, err := scanner.Scan()
			if err == io.EOF {
				require.Fail(t, "should return an error")
			}
			if err == nil {
				continue
			}
			var ferr *file.Error
			require.ErrorAs(t, err, &ferr)
			require.Equal(t, 3, ferr.Line)
			break
		}
	})
}
