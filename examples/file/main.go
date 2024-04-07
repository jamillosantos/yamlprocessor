package main

import (
	"fmt"

	"github.com/jamillosantos/yamlprocessor"
)

type Entry struct {
	Year  int    `yaml:"year"`
	Event string `yaml:"event"`
}

type Person struct {
	Name    string  `yaml:"name"`
	Age     int     `yaml:"age"`
	History []Entry `yaml:"history"`
}

func main() {
	d := []byte(`
name: "John Doe"
age: ${2024-1985}
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
