# yamlprocessor

Integrates [expr-lang/expr](https://github.com/expr-lang/expr) with [goccy/go-yaml](https://github.com/goccy/go-yaml).

:warning: **This library is in very early stage of development. Use at your own risk.**

So far, all the work that has been done was experimental to check feasibility and there are a lot of flaws, poor error
handling, performance issues and lack tests.

## What does it do?

This library allows you to use expressions in your YAML files. The expressions are evaluated using the `expr` library.

```yaml
name: "John Doe"
age: ${2024 - 2000}
history:
  ${file("history.yaml")}
```

## Usage

```go
package main

import "github.com/jamillosantos/yamlprocessor"

type Entry struct {
	Year int `yaml:"year"`
	Event string `yaml:"event"`
}

type Person struct {
	Name    string `yaml:"name"`
	Age     int   `yaml:"age"`
	History []Entry `yaml:"history"`
}

func main() {
	d := []byte(`
name: "John Doe"
age: ${2024 - 2000}
history:
  ${file("history.yaml")}
`)
	var p Person
	err := yamlprocessor.Unmarshal(d, &p)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", p)
}

// history.yaml
// - year: 2000
//   event: "Something happened"
// - year: 2001
//   event: "Something else happened"


```

## Adding new functions

You can add new capabilities to the processor by adding new functions to the processor.

```go
// ...
processor := yamlprocessor.NewProcessor()
processor.Env.Add("add", func(ctx *yamlprocessor.Context) any {
	return func(a, b int) int {
        return a + b
    }
})

var p Person
err := processor.Unmarshal([]byte(`name: "Name ${add(1, 2)}"`), &p)

```

## Adding new variables

You can add new variables to the processor by adding new variables to the processor.

```go
processor := yamlprocessor.NewProcessor()
processor.Env.Add("var1", yamlprocessor.Value(10))
```