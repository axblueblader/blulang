package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var globalScope = NewGlobalScope()
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Println("Error reading file")
			return
		}
		if !strings.HasSuffix(file.Name(), ".blu") {
			fmt.Println("Invalid file format")
			return
		}
		source, err := os.ReadFile(file.Name())
		if err != nil {
			fmt.Println("Error reading file")
			return
		}
		parser := NewParser()
		ast := parser.CreateAST(string(source))
		//fmt.Println("AST:", ast)
		Eval(ast, globalScope)
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("BluLang REPL, type a statement and press enter, type 'exit()' to quit")
		line := ""
		for line != "exit()" {
			fmt.Print("> ")
			scanner.Scan()
			line = scanner.Text()
			parser := NewParser()
			ast := parser.CreateAST(line)
			//fmt.Println("AST:", ast)
			result := Eval(ast, globalScope)
			fmt.Printf("%v :: %v\n", result.Kind(), result.Value())
		}
	}
}
