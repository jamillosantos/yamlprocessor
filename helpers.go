package yamlprocessor

var DefaultProcessor = NewProcessor()

func Add(name string, function ExprEntry) error {
	return DefaultProcessor.Env.Add(name, function)
}

var Unmarshal = DefaultProcessor.Unmarshal

func init() {
	DefaultProcessor.Env["file"] = File
	DefaultProcessor.Env["env"] = Env
}

// Value is a helper function that helps to add variables in the ExprEnv map for when expressions are evaluated.
func Value(val any) ExprEntry {
	return func(ctx *Context) any {
		return val
	}
}
