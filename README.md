# JSON Linter

This is a solution to [John Crickett's Coding Challenge #2 - Build a JSON Parser](https://codingchallenges.substack.com/p/coding-challenge-2) written in Golang.

In the source code, you will notice that I call the main application a linter (as opposed to a parser per John's Challenge). It made more sense to me to call this a linter as the task is to create an application to validate whether a passed in JSON file is valid or not.

To achieve this, the validator uses a `lexer` and a `parser` to help determine if the JSON is valid or not.

## Pre-requisites

Requires you to have [golang](https://go.dev/doc/install) installed.

To run the files using `make`, make sure you have [GNU make](https://www.gnu.org/software/make/) installed as well.

The Makefile includes a `lint` command which depends on [golangci-lint](https://github.com/golangci/golangci-lint). If you want to use `make lint`, then install this dependency.

## Installation

1. Clone and navigate into the repository
2. Build the executable binaries using `make build`
3. The binaries can be found in the `bin` directory of the repo

## Usage

```bash
# Build the Binary
make build

# Navigate to the bin directory where the binary is stored
cd ./bin

# Execute the application
./jv <json filepath>
```

## Experimenting w/ a Scratch File

For development purposes, you can create a `./scratch/scratch.go` file and use the `make build/scratch` and `make run/scratch` makefile commands for testing purposes.

Here is what a `scratch.go` file may look like
```go
package main

import "fmt"

func main() {
	// YOUR SCRATCH CODE HERE...
	fmt.Println("Scratch Works!")
}
```

## Notes / Background

### JSON Structure

At the top level, a `JSON` file can either be:
- `object`
- `array`
- These are actually just `values`, so at a high level, the root element of a JSON file is just a `value` (see below for more info on values)

`object`
- An unordered set of `key`-`value` pairs
	- Can also refer to each `key`-`value` pair as a `property`
		- The `key` is always a string
		- Refer to notes on `value` for the types
	- An object con contain 0 or more `properties`
- Begins w/ a left brace `{`
- Ends w/ a right brace `}`
- Each key is followed by a colon `:`
- `properties` (except for the last one) are separated by a comma `,`

`array`
- An ordered set of `value`s
	- An array can contain 0 or more `value`s
- Begins w/ a left bracket `[`
- Ends w/ a right bracket `]`
- Each `value` (except for last one) are separated by a comma `,`
- Refer to notes on `value` for types

`value`
- `object`
- `array`
- `string` (wrapped in double quotes)
- `number`
- boolean (`true` or `false`)
- `null`

More info on the structure of JSON can be found [here](json.org/json-en.html)

### Lexical Analysis / Tokenization / Scanning

During `Lexical Analysis`, a `lexer` (a.k.a. `tokenizer` / `scanner`) takes a sequence of characters as input, and outputs the series of lexical `tokens`

Think of the process of `Lexical Analysis` as taking a sequence of characters and classifying the important parts

```txt
string -> LEXER -> tokens
```

What is a `Lexeme`?
- A `lexeme` is the actual sequence of characters in the source code that matches the pattern for a particular token.
- The smallest unit of meaning in the source code (the raw, unclassified piece of code being analyzed)
- `Lexeme`s are often the input substrings that are recognized by the lexical analyzer based on the language's syntax rules
- Ex. `var x = 10;` => Lexemes: `var`, `x`, `=`, `10`, `;`

What is a `Token`?
- In the context of `Lexical Analysis`, a `token` is a categorized and labeled unit of meaning that results from the `lexical analysis` of `lexemes`.
- The `token` represents the classification of the lexeme based on the language's syntax rules
- Often includes a type (e.g. keyword, identifier, operator, etc.), and an optional attribute value
- A `token` consists of the combination of a `lexeme`, its `token type` / classification, and can include optional info, such as a value for a string or numerical `lexeme`

Example: Tokenize the following string of characters: `var x = 2 + (4 * 10);`
- Each row in the table represents a token

| Lexeme | Classification / Token Type | Value |
| ------ | --------------------------- | ----- |
| `var`  | VAR                         | n/a   |
| `x`    | IDENT                       | n/a   |
| `=`    | ASSIGN                      | n/a   |
| `2`    | INT                         | 2     |
| `+`    | ADD                         | n/a   |
| `(`    | LPAREN                      | n/a   |
| `4`    | INT                         | 4     |
| `*`    | MUL                         | n/a   |
| `10`   | INT                         | 10    |
| `)`    | RPAREN                      | n/a   |
| `;`    | SEMI                        | n/a   |

### Syntactical Analysis / Parsing

Parsing is the process of recognizing the shape of an input string.

During `Syntactical Analysis`, a `parser` takes a series of tokens, and outputs an `Abstract Syntax Tree` (`AST`)

```txt
tokens -> PARSER -> Abstract Syntax Tree
```

What is an `Abstract Syntax Tree`?
- An `AST` is an intermediate representation of source code as a tree structure
- The `AST` represents the hierarchial structure and semantics of the code
- A compiler will use an `AST` to compile a language down to machine code.

Here is what an `AST` could look like for the previous example: `var x = 2 + (4 * 10);`
- The root node is a "VariableDeclaration" representing the declaration of the variable "x."
- The "dataType" property indicates the data type (in this case, "var").
- The "declarations" array contains a "VariableDeclarator" node for the variable "x."
- Inside the "VariableDeclarator," the "init" property contains a "BinaryExpression" representing the initialization of "x."
- The "BinaryExpression" nodes represent the binary operations, such as addition and multiplication.
- The "Literal" nodes represent integer literals.

```JSON
{
	"type": "VariableDeclaration",
	"dataType": "var",
	"declarations": [
		{
			"type": "VariableDeclarator",
			"identifier": "x",
			"init": {
				"type": "BinaryExpression",
				"operator": "=",
				"left": {
					"type": "Literal",
					"value": 2
				},
				"right": {
					"type": "BinaryExpression",
					"operator": "+",
					"left": {
						"type": "Literal",
						"value": 4
					},
					"right": {
						"type": "BinaryExpression",
						"operator": "*",
						"left": {
							"type": "Literal",
							"value": 10
						},
						"right": null  // Placeholder for the missing operand (could be another literal, identifier, etc.)
					}
				}
			}
		}
	]
}
```

## Resources / References

- [John Crickett - Coding Challenge #2 - Build a JSON Parser](https://codingchallenges.substack.com/p/coding-challenge-2)
- [Wikipedia - Lexical Analysis](https://en.wikipedia.org/wiki/Lexical_analysis)
- [Wikipedia - Parsing](https://en.wikipedia.org/wiki/Parsing)
- [Wikipedia - Abstract Syntax Tree](https://en.wikipedia.org/wiki/Abstract_syntax_tree)
- [json.org - Introducing JSON](https://www.json.org/json-en.html)
- [STD 90 - RFC 8259 - The JavaScript Object Notation (JSON) Interchange Format](https://www.rfc-editor.org/info/std90)
- [Computerphile - Parsing](https://www.youtube.com/watch?v=r6vNthpQtSI)
- [TJ DeVries - Writing an interpreter... in OCaml?!?](https://www.youtube.com/watch?v=NjKJ9-ejR6o)
- [TJ DeVries - Writing our own parser in OCaml!](https://www.youtube.com/watch?v=dycsRSOQjho)
- [Tsoding Daily - What is a Lexer (No BS explanation)](https://www.youtube.com/watch?v=BI3K-ME3L74&t=169s)
- [Dmitry Soshnikov - Building a Parser from scratch](https://www.youtube.com/playlist?list=PLGNbPb3dQJ_5FTPfFIg28UxuMpu7k0eT4)
- [Dmitry Soshnikov - Parsing Algorithms. Lecture [5/22] Abstract Syntax Trees](https://www.youtube.com/watch?v=VKM1eLoN-gI)
- [newline - What is An Abstract Syntax Tree, with WealthFont Engineer Spencer Miskoviak](https://www.youtube.com/watch?v=wINY109MG10)
- [newline - Understand Abstract Syntax Trees - ASTs - In Practical and Useful Ways for Frontend Developers](https://www.youtube.com/watch?v=tM_S-pa4xDk)
- [AST Explorer](https://astexplorer.net/)
- [Crafting Interpreters](https://craftinginterpreters.com/contents.html)
- [Rob Pike - Lexical Scanning in Go](https://www.youtube.com/watch?v=HxaD_trXwRE)
- [Aaron Raff - How to Write a Lexer in Go](https://www.aaronraff.dev/blog/how-to-write-a-lexer-in-go)
- [Bradford Lamson-Scribner - Building a JSON Parser and Query Tool with Go](https://medium.com/@bradford_hamilton/building-a-json-parser-and-query-tool-with-go-8790beee239a)
- [Thorsten Ball - Writing an Interpreter in Go](https://interpreterbook.com/)
- [Oguzhan Olguncu - Write Your Own JSON Parsers with Node and Typescript](https://ogzhanolguncu.com/blog/write-your-own-json-parser)