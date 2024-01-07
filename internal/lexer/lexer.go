package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// Define Position Struct to track the current position of lexer's reader
type LexerPosition struct {
	Line   int // Current line Lexer's reader is scanning
	Column int // Current column position of Lexer's reader
}

// lexer struct is responsible for tokenizing input
type Lexer struct {
	Reader *bufio.Reader // Reader object of file to be tokenized
	Pos    LexerPosition
}

// CreateLexer creates & returns a new lexer instance for lexical analysis of the input from the given reader.
//
// The lexer is initialized with a buffered reader for efficient reading and the initial position set to the beginning (line 1, column 0).
//
// Parameters:
//   - reader: An io.Reader representing the source of input for lexical analysis.
//
// Returns:
//   - A pointer to the created lexer.
func CreateLexer(reader io.Reader) *Lexer {
	lxrPtr := &Lexer{
		Reader: bufio.NewReader(reader),
		Pos:    LexerPosition{Line: 1, Column: 0},
	}

	return lxrPtr
}

// Lex is responsible for opening the JSON file specified at the filePath.
// Returns a slice of Tokens representing the JSON file.
func Lex(filePath string) []Token {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return []Token{}
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	lxr := CreateLexer(reader)

	var tokens []Token
	for {
		tok := lxr.GetNextToken()

		// Break the loop if EOF is reached
		if tok.TokType == EOF {
			break
		}

		tokens = append(tokens, tok)
	}
	return tokens
}

// GetNextToken scans the Lexer's input to return the next token
func (lxr *Lexer) GetNextToken() Token {
	var token Token

	// Keep scanning until a token is found or EOF is reached
	for {
		r, err := lxr.advanceReader()
		if err != nil {
			if err == io.EOF {
				token = createToken(EOF, lxr.Pos, '0')
				return token
			}
			panic(err)
		}

		// Skip whitespace / tabs before proceeding
		if r == ' ' || r == '\t' {
			continue
		}

		// Evaluate the rune (r) at the current scan position
		switch r {
		case '0':
			token = createToken(EOF, lxr.Pos, r)
			return token
		case '\n':
			// Reset lexer's position at each newline
			lxr.resetPosition()
		case '{':
			token = createToken(LBRACE, lxr.Pos, r)
			return token
		case '}':
			token = createToken(RBRACE, lxr.Pos, r)
			return token
		case '[':
			token = createToken(LBRACKET, lxr.Pos, r)
			return token
		case ']':
			token = createToken(RBRACKET, lxr.Pos, r)
			return token
		case ',':
			token = createToken(COMMA, lxr.Pos, r)
			return token
		case ':':
			token = createToken(COLON, lxr.Pos, r)
			return token
		case '"':
			return handleStringToken(lxr, r)
		default:
			if isNumberMaybe(r) {
				return handleNumberToken(lxr, r)
			} else if unicode.IsLetter(r) {
				return handleIdentifierToken(lxr, r)
			} else {
				// Handle Unknown Tokens
				token = createToken(ILLEGAL, lxr.Pos, r)
				return token
			}
		}
	}
}

// resetPosition is a helper func to reset the pos of the lexer to the next line and 0th column position
func (lxr *Lexer) resetPosition() {
	lxr.Pos.Line++
	lxr.Pos.Column = 0
}

// advanceReader moves the reader position forwarder by 1 rune & updates the Lexer's position
func (lxr *Lexer) advanceReader() (rune, error) {
	r, _, err := lxr.Reader.ReadRune()
	if err != nil {
		return 0, err // Return error
	}

	lxr.Pos.Column++ // Advance position of lexer

	return r, nil // Return rune and no error
}

// backupReader backs up the reader by 1 position
func (lxr *Lexer) backupReader() {
	// Only backup if we are not at column position 0
	if lxr.Pos.Column > 0 {
		err := lxr.Reader.UnreadRune()
		if err != nil {
			panic(err)
		}

		lxr.Pos.Column-- // Backup column position
	}
}

// peekForward peeks forward by specified number of steps without advancing the reader's position.
// Defaults to one step if steps is not provided.
func (lxr *Lexer) peekForward(steps ...int) (rune, error) {
	numSteps := 1
	if len(steps) > 0 {
		numSteps = steps[0]
	}

	var runePeeked rune

	// Save the current position in case an error occurs
	savePos := lxr.Pos

	// Advance the reader by the specified number of steps
	for i := 0; i < numSteps; i++ {
		r, _, err := lxr.Reader.ReadRune()
		if err != nil {
			// Restore the position and return the error
			lxr.Pos = savePos
			return 0, err
		}

		runePeeked = r
	}

	// Backup the reader by the specified number of steps
	for i := 0; i < numSteps; i++ {
		err := lxr.Reader.UnreadRune()
		if err != nil {
			// Restore the position and return the error
			lxr.Pos = savePos
			return 0, err
		}
	}

	return runePeeked, nil
}

// createToken creates & returns a new token based on the specified TokenType,
// lexer position, and a variable number of runes representing the lexeme.
//
// Parameters:
//   - tokenType: TokenType - The type of the token, such as Identifier, Number, etc.
//   - pos: lexerPosition - The position information (line and column) where the lexeme starts.
//   - lexemeChars ...rune: A variadic parameter allowing the passing of one or more runes
//     representing the characters of the lexeme.
//
// Returns:
//   - Token: A newly created token containing information about the token type, lexeme,
//     and position in the source code.
func createToken(tokType TokenType, pos LexerPosition, lexemeChars ...rune) Token {
	// Convert lexemeChar to a string
	lexemeStr := string(lexemeChars)

	// Handle special characters & update colEnd accordingly
	visualWidth := len([]rune(lexemeStr)) // Consideration for characters with escape sequences
	colEnd := pos.Column + visualWidth - 1

	// Generate a new struct to store data for token position
	tokenPos := TokenPosition{
		Line:     pos.Line,
		ColStart: pos.Column,
		ColEnd:   colEnd,
	}

	// Generate a new token struct
	token := Token{
		TokType: tokType,
		Lexeme:  lexemeStr,
		TokPos:  tokenPos,
	}

	// Update token for EOF condition
	if tokType == EOF {
		token.Lexeme = "EOF"
		token.TokPos.ColEnd = token.TokPos.ColStart
	}

	return token
}

// isNumberMaybe checks if the rune at the current position could be a number.
func isNumberMaybe(r rune) bool {
	return (r >= '0' && r <= '9') || r == '-' || r == '.' || r == 'e' || r == 'E'
}

// isValidJSONNumber checks if the given runes form a valid JSON number.
func isValidJSONNumber(runes []rune) bool {
	input := string(runes)

	// TODO: Fix to also accept -.# pattern, also e#, E#, e-#, or E-#
	jsonNumberPattern := `^-?(?:0|[1-9]\d*)(?:\.\d+)?(?:[eE][+-]?\d+)?$`

	return regexp.MustCompile(jsonNumberPattern).MatchString(input)
}

// handleNumberToken returns NUM or ILLEGAL token
func handleNumberToken(lxr *Lexer, r rune) Token {

	var token Token
	lxr.backupReader()
	numRune, startPos, err := lxr.readNumber()
	if err != nil {
		if strings.Contains(err.Error(), "invalid JSON number") {
			token = createToken(ILLEGAL, startPos, numRune...)
			return token
		}
		// Invalid number, return Unknown Token
		token = createToken(ILLEGAL, startPos, r)
	} else {
		token = createToken(NUM, startPos, numRune...)
	}
	return token
}

// readNumber reads attempts to read in a number and return the read in value
func (lxr *Lexer) readNumber() ([]rune, LexerPosition, error) {
	var num []rune

	// Store starting position
	startPos := LexerPosition{
		Line:   lxr.Pos.Line,
		Column: lxr.Pos.Column + 1,
	}

	// Keep reading until hit a non-numeric condition
	for {
		r, err := lxr.advanceReader()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, startPos, err
		}

		if unicode.IsSpace(r) || !isNumberMaybe(r) {
			lxr.backupReader()
			break
		}

		num = append(num, r)
	}

	if !isValidJSONNumber(num) {
		return num, startPos, errors.New("invalid JSON number")
	}

	return num, startPos, nil
}

// handleStringToken returns STR or ILLEGAL token
func handleStringToken(lxr *Lexer, r rune) Token {
	var token Token
	strRune, startPos, err := lxr.readString()
	if err != nil || len(strRune) == 0 {
		// Invalid string, return Unknown Token
		token = createToken(ILLEGAL, startPos, r)
	} else {
		token = createToken(STR, startPos, strRune...)
	}
	return token
}

// readString reads the string from the current position of the Lexer's reader
func (lxr *Lexer) readString() ([]rune, LexerPosition, error) {
	var str []rune

	// Store starting position
	startPos := LexerPosition{
		Line:   lxr.Pos.Line,
		Column: lxr.Pos.Column + 1,
	}

	for {
		r, err := lxr.advanceReader()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, startPos, err
		}

		// Break if we hit the next " and it is not escaped
		if r == '"' && str[len(str)-1] != '\\' {
			break
		}

		str = append(str, r)
	}

	return str, startPos, nil
}

// handleIdentifierToken returns TRUE, FALSE, NULL or ILLEGAL token
func handleIdentifierToken(lxr *Lexer, r rune) Token {
	var token Token
	lxr.backupReader()
	identRune, startPos, err := lxr.readIdentifier()
	if err != nil {
		// Invalid string, return Unknown Token
		token = createToken(ILLEGAL, startPos, r)
	} else if string(identRune) == "true" {
		token = createToken(TRUE, startPos, identRune...)
	} else if string(identRune) == "false" {
		token = createToken(FALSE, startPos, identRune...)

	} else if string(identRune) == "null" {
		token = createToken(NULL, startPos, identRune...)
	} else {
		token = createToken(ILLEGAL, startPos, identRune...)
	}
	return token
}

// readIdentifier attempts to read an identifier
func (lxr *Lexer) readIdentifier() ([]rune, LexerPosition, error) {
	var ident []rune

	// Store starting position
	startPos := LexerPosition{
		Line:   lxr.Pos.Line,
		Column: lxr.Pos.Column + 1,
	}

	for {
		r, err := lxr.advanceReader()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, startPos, err
		}

		if unicode.IsSpace(r) || !unicode.IsLetter(r) {
			lxr.backupReader()
			break
		}

		ident = append(ident, r)
	}

	return ident, startPos, nil
}
