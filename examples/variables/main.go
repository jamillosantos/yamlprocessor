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
	err := processor.Env.Add("var1", yamlprocessor.Value(10))
	if err != nil {
		panic(err)
	}

	var p Person
	err = processor.Unmarshal([]byte(`name: "Name ${var1}"`), &p)
	if err != nil {
		panic(err)
	}

	fmt.Println(p.Name)
}
