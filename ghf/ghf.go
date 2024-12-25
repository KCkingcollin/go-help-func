// Main Helper Functions
package ghf

import (
	"fmt"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

var Verbose bool = verbose()

// Returns true if the user adds -v/--verbose option at start up
func verbose() bool {
    if len(os.Args) > 1 {
        if os.Args[1] == "--verbose" || os.Args[1] == "-v" {
            return true
        }
    }
    return false
}

// Load a file via path, and return it as a string
func LoadFile(path string) string {
    data, err := os.ReadFile(path)
    if err != nil && Verbose {
        fmt.Println(err)
    }
    return string(data)
}

// Converts a 64 bit mgl slice to a 32 mgl slice
func Mgl64to32Slice[T mgl64.Mat4 | mgl64.Vec3](data []T) interface{} {
    switch data := any(data).(type) {
    case []mgl64.Mat4:
        var m32 mgl32.Mat4
        ms32 := make([]mgl32.Mat4, len(data))
        for i := range data {
            ms64 := [16]float64(data[i])
            for j := range ms64 {
                m32[j] = float32(ms64[j])
            }
            ms32[i] = m32
        }
        return ms32

    case []mgl64.Vec3:
        var m32 mgl32.Vec3
        ms32 := make([]mgl32.Vec3, len(data))
        for i := range data {
            ms64 := [3]float64(data[i])
            for j := range ms64 {
                m32[j] = float32(ms64[j])
            }
            ms32[i] = m32
        }
        return ms32

    default:
        return nil
    }
}

// Converts a 64Mat4 to a 32Mat4
func Mgl64to32Mat4(m64 mgl64.Mat4) mgl32.Mat4 {
    m32 := make([]float32, 16)
    for i := range [16]float64(m64) {
        m32[i] = float32(m64[i])
    }
    return mgl32.Mat4(m32)
}
