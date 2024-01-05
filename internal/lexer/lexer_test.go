package lexer

import (
	"strings"
	"testing"
)

func TestGetNextToken(t *testing.T) {
	// Define tests cases
	testCases := []struct {
		input            string
		expectedTokens   []Token
		expectedErrorMsg string
	}{
		// Testing braces
		{
			input: `{}`,
			expectedTokens: []Token{
				{LBRACE, "{", TokenPosition{1, 1, 1}},
				{RBRACE, "}", TokenPosition{1, 2, 2}},
			},
		},
		// Testing empty string
		{
			input:          "",
			expectedTokens: []Token{},
		},
		// Testing brackets
		{
			input: `[]`,
			expectedTokens: []Token{
				{LBRACKET, "[", TokenPosition{1, 1, 1}},
				{RBRACKET, "]", TokenPosition{1, 2, 2}},
			},
		},
		// Testing brackets and string
		{
			input: `["hello"]`,
			expectedTokens: []Token{
				{LBRACKET, "[", TokenPosition{1, 1, 1}},
				{STR, "hello", TokenPosition{1, 3, 7}},
				{RBRACKET, "]", TokenPosition{1, 9, 9}},
			},
		},
		// Testing strings
		{
			input: `"a", "bc", "def","ghij" whaat "`,
			expectedTokens: []Token{
				{STR, "a", TokenPosition{1, 2, 2}},
				{COMMA, ",", TokenPosition{1, 4, 4}},
				{STR, "bc", TokenPosition{1, 7, 8}},
				{COMMA, ",", TokenPosition{1, 10, 10}},
				{STR, "def", TokenPosition{1, 13, 15}},
				{COMMA, ",", TokenPosition{1, 17, 17}},
				{STR, "ghij", TokenPosition{1, 19, 22}},
				{ILLEGAL, "whaat", TokenPosition{1, 25, 29}},
				{ILLEGAL, "\"", TokenPosition{1, 32, 32}},
			},
		},
		// Testing identifiers
		{
			input: `invalid true false null`,
			expectedTokens: []Token{
				{ILLEGAL, "invalid", TokenPosition{1, 1, 7}},
				{TRUE, "true", TokenPosition{1, 9, 12}},
				{FALSE, "false", TokenPosition{1, 14, 18}},
				{NULL, "null", TokenPosition{1, 20, 23}},
			},
		},
		// Testing numbers
		{
			input: `123 1.23 -1.23 1.23e10 -1.23e10 1.23e-10 -1.23e-10 1.23E10 -1.23E10 1.23E-10 -1.23E-10 e10 e-10 E10 E-10 -1.2.3 --1.2.3`,
			expectedTokens: []Token{
				{NUM, "123", TokenPosition{1, 1, 3}},
				{NUM, "1.23", TokenPosition{1, 5, 8}},
				{NUM, "-1.23", TokenPosition{1, 10, 14}},
				{NUM, "1.23e10", TokenPosition{1, 16, 22}},
				{NUM, "-1.23e10", TokenPosition{1, 24, 31}},
				{NUM, "1.23e-10", TokenPosition{1, 33, 40}},
				{NUM, "-1.23e-10", TokenPosition{1, 42, 50}},
				{NUM, "1.23E10", TokenPosition{1, 52, 58}},
				{NUM, "-1.23E10", TokenPosition{1, 60, 67}},
				{NUM, "1.23E-10", TokenPosition{1, 69, 76}},
				{NUM, "-1.23E-10", TokenPosition{1, 78, 86}},
				{ILLEGAL, "e10", TokenPosition{1, 88, 90}},
				{ILLEGAL, "e-10", TokenPosition{1, 92, 95}},
				{ILLEGAL, "E10", TokenPosition{1, 97, 99}},
				{ILLEGAL, "E-10", TokenPosition{1, 101, 104}},
				{ILLEGAL, "-1.2.3", TokenPosition{1, 106, 111}},
				{ILLEGAL, "--1.2.3", TokenPosition{1, 113, 119}},
			},
		},
	}

	for _, testCase := range testCases {
		// Run a subtest for each test case
		// The callback function passed to t.Run is executed for each testCase in the parent for loop
		t.Run(testCase.input, func(t *testing.T) {
			// Create a reader from the input string
			reader := strings.NewReader(testCase.input)

			// Create a lexer for testing
			lexer := CreateLexer(reader)

			// Iterate through expected tokens and compare with actual tokens
			for _, expectedToken := range testCase.expectedTokens {
				actualToken := lexer.GetNextToken()
				assertTokenEquality(t, expectedToken, actualToken)
			}

			// This if statement is an example of "statement initialization" where we declare a variable within the if statement (actualToken) and use it in the if statement
			// This variable is only in-scope for the if statement
			if actualToken := lexer.GetNextToken(); actualToken.TokType != EOF {
				t.Errorf("Expected EOF, got %v", actualToken.TokType)
			}
		})
	}
}

func assertTokenEquality(t *testing.T, expected Token, actual Token) {
	if expected.TokType != actual.TokType {
		t.Errorf("Expected token type %v, got %v", expected.TokType, actual.TokType)
	}
	if expected.Lexeme != actual.Lexeme {
		t.Errorf("Expected lexeme %v, got %v", expected.Lexeme, actual.Lexeme)
	}
	if expected.TokPos.Line != actual.TokPos.Line ||
		expected.TokPos.ColStart != actual.TokPos.ColStart ||
		expected.TokPos.ColEnd != actual.TokPos.ColEnd {
		t.Errorf("Expected token position %v, got %v", expected.TokPos, actual.TokPos)
	}
}
