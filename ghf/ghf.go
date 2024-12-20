// Main Helper Functions
package ghf

import (
	"fmt"
	"os"

	"github.com/go-gl/mathgl/mgl32"
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

func Mat4ToFloat32(mats []mgl32.Mat4) []float32 {
	float32s := make([]float32, 0)

	for _, mat := range mats {
		for i := 0; i < 16; i++ {
            float32s = append(float32s, mat[i])
		}
	}

	return float32s
}
