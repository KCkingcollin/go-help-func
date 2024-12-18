// Main Helper Functions
package ghf

import (
	"fmt"
	"os"
)

var Verbose bool = verbose()

func verbose() bool {
    verbose := false
    if len(os.Args) > 1 {
        if os.Args[1] == "--verbose" || os.Args[1] == "-v" {
            verbose = true
        }
    }
    return verbose
}

func LoadFile(path string) string {
    data, err := os.ReadFile(path)
    if err != nil && Verbose {
        fmt.Println(err)
    }
    return string(data)
}

