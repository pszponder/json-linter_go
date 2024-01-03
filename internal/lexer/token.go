package lexer

// Define type alias for the type of the token (can later change this to a new type like string)
type TokenType int

// Define the a list of TokenType constants (tokenType can only one of these)
const (
	// Default type if the character does not match any of the other defined types
	ILLEGAL TokenType = iota

	// End of File Token
	EOF

	// Structural Tokens
	LBRACE   // Left Curly Brace "{"
	RBRACE   // Right Curly Brace "}"
	LBRACKET // Left Square Bracket "["
	RBRACKET // Right Square Bracket "["
	COMMA    // Comma ","
	COLON    // Colon ":"

	// Values
	STR // string
	NUM // number (integer or float)

	// Identifiers / Keywords (special values)
	TRUE  // true
	FALSE // false
	NULL  // null
)

// Define Position Struct for token positional context
type TokenPosition struct {
	Line     int // Line number Token is found on
	ColStart int // Column start position of Token
	ColEnd   int // Column end position of Token
}

// Define the Token Struct
type Token struct {
	TokType TokenType
	Lexeme  string // The literal which Token represents
	TokPos  TokenPosition
}
