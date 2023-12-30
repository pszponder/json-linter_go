package args

import (
	"fmt"
	"os"
)

// GetFilePath parses and returns passed in filepath
//
// Returns:
// string containing the file path
func GetFilePath() string {
	if len(os.Args) != 2 {
		fmt.Println("Usage: parser <filepath>")
		os.Exit(1) // Exit the app w/ a non-zero status code to indicate an error
	}

	return os.Args[1] // 1st arg is the app binary, 2nd arg should be the filepath
}
