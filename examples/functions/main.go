package main

import (
	"fmt"

	"github.com/jamillosantos/yamlprocessor"
)

type Person struct {
	Name string `yaml:"name"`
}

func main() {
	processor := yamlprocessor.NewProcessor()
	err := processor.Env.Add("add", func(ctx *yamlprocessor.Context) any {
		return func(a, b int) int {
			return a + b
		}
	})
	if err != nil {
		panic(err)
	}

	var p Person
	err = processor.Unmarshal([]byte(`name: "Name ${add(1, 2)}"`), &p)
	if err != nil {
		panic(err)
	}

	fmt.Println(p.Name)
}
