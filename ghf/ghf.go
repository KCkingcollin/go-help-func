// Main Helper Functions
package ghf

import (
	"fmt"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
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

func Mgl64to32Mat4Slice(m64 []mgl64.Mat4) []mgl32.Mat4 {
    var ms32 mgl32.Mat4
    m32 := make([]mgl32.Mat4, len(m64))
    for i := range m64 {
        ms64 := [16]float64(m64[i])
        for j := range ms64 {
            ms32[j] = float32(ms64[j])
        }
        m32[i] = ms32
    }
    return m32
}

func Mgl64to32Mat4(m64 mgl64.Mat4) mgl32.Mat4 {
    m32 := make([]float32, 16)
    for i := range [16]float64(m64) {
        m32[i] = float32(m64[i])
    }
    return mgl32.Mat4(m32)
}
