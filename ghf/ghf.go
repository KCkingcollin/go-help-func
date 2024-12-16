package ghf

import "os"

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

