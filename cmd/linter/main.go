package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pszponder/json-linter_go/internal/args"
	"github.com/pszponder/json-linter_go/internal/lexer"
	"github.com/pszponder/json-linter_go/internal/parser"
)

func main() {
	// Retrieve filepath to the file to validate
	filePath := args.GetFilePath()
	fmt.Println(filePath)

	// Pass in the file to a lexer in order to generate a token representation of the file (tokenize it)
	tokens := lexer.Lex(filePath)

	// Parse the tokens and determine if the JSON is valid
	_, err := parser.ParseJSON(tokens)
	if err != nil {
		log.Print("Error: ", err)
		os.Exit(1)
	}

	log.Printf("JSON file located in %v is valid", filePath)
	os.Exit(0)
}
