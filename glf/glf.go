// GL Helper Functions
package glf

import (
	"errors"
	"fmt"
	"image/png"
	"math"
	"os"
	"strings"
	"unsafe"

	"github.com/KCkingcollin/go-help-func/ghf"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

var Verbose bool = ghf.Verbose

// Prints OpenGL version information
func PrintVersionGL() {
    version := gl.GoStr(gl.GetString(gl.VERSION))
    if Verbose {
        fmt.Println("OpenGL Version", version)
    }
}

// Load a RGBA texture file via path, and returns a uint32 as texture ID
func LoadTexture(filename string) uint32 {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}

    texture := GenBindTexture()

    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
    
    gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

    gl.GenerateMipmap(gl.TEXTURE_2D)

	return texture
}

// Creates and binds a texture ID via a uint32 ID
func GenBindTexture() uint32 {
    var texID uint32
    gl.GenTextures(1, &texID)
    BindTexture(texID)
    return texID
}

// Binds a 2D texture ID via a uint32 ID
func BindTexture(texID uint32) {
    gl.BindTexture(gl.TEXTURE_2D, texID)
}

// Create and bind shader program via a vertex and fragment glsl source file path in that order.
//
//  - Returns shader ID as a uint32 if no errors
//  - Returns error if shader creation fails
func CreateProgram(vertPath, fragPath string) (uint32, error) {
    vertexShader, err := CreateVertexShader(vertPath)
    if err != nil && Verbose {
        fmt.Printf("Failed to compile vertex shader: %s \n", err)
    } else if Verbose {
        println("Vertex shader compiled successfully")
    }
    fragmentShader, err := CreateFragmentShader(fragPath)
    if err != nil && Verbose {
        fmt.Printf("Failed to compile fragment shader: %s \n", err)
    } else if Verbose {
        println("Fragment shader compiled successfully")
    }
    ProgramID := gl.CreateProgram()
    gl.AttachShader(ProgramID, vertexShader)
    gl.AttachShader(ProgramID, fragmentShader)
    gl.LinkProgram(ProgramID)
    var success int32
    gl.GetProgramiv(ProgramID, gl.LINK_STATUS, &success)
    if success == gl.FALSE {
        var logLength int32
        gl.GetProgramiv(ProgramID, gl.INFO_LOG_LENGTH, &logLength)
        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetProgramInfoLog(ProgramID, logLength, nil, gl.Str(log))
        return ProgramID, errors.New(log)
    }

    gl.DeleteShader(vertexShader)
    gl.DeleteShader(fragmentShader)

    return ProgramID, err
}

// Create vertex shader from path to glsl vertex shader source
//
//  - Returns shader ID as a uint32 if no errors
//  - Returns error if shader creation fails
func CreateVertexShader(shaderFile string) (uint32, error) {
    ShaderSource := ghf.LoadFile(shaderFile) + "\x00"
    return CreateShader(ShaderSource,  gl.VERTEX_SHADER)
}

// Create fragment shader from path to glsl fragment shader source
//
//  - Returns shader ID as a uint32 if no errors
//  - Returns error if shader creation fails
func CreateFragmentShader(shaderFile string) (uint32, error) {
    ShaderSource := ghf.LoadFile(shaderFile) + "\x00"
    return CreateShader(ShaderSource,  gl.FRAGMENT_SHADER)
}

// Create shader via shader source file and shader type.
//
//  - Returns shader ID as a uint32 if no errors
//  - Returns error if shader creation fails
func CreateShader(ShaderSource string, ShaderType uint32) (uint32, error) {
    ShaderID:= gl.CreateShader(ShaderType)
    csource, free := gl.Strs(ShaderSource)
    gl.ShaderSource(ShaderID, 1, csource, nil)
    free()
    gl.CompileShader(ShaderID)
    var status int32
    gl.GetShaderiv(ShaderID, gl.COMPILE_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetShaderiv(ShaderID, gl.INFO_LOG_LENGTH, &logLength)
        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetShaderInfoLog(ShaderID, logLength, nil, gl.Str(log))
        return ShaderID, errors.New(log)
    }
    return ShaderID, nil 
}

// Create and initialize buffer via generic slice (only takes float32 for now)
func BufferData[T float32](target uint32, data []T, usage uint32) {
    // switch any(data).(type) {
    // case []float32:
    //     byteSize = 4
    // case []float64:
    //     byteSize = 8
    // default:
    //     panic("unsupported type for BufferData")
    // }
    gl.BufferData(target, len(data)*4, gl.Ptr(data), usage)
}

// Generate and bind buffers, and return the ID as a uint32
func GenBindBuffers(target uint32) uint32 {
    var buffer uint32
    gl.GenBuffers(1, &buffer)
    gl.BindBuffer(target, buffer)
    return buffer
}

// Generate and bind vertex arrays, and returns the ID as a uint32.
func GenBindVertexArrays() uint32 {
    var VAO uint32
    gl.GenVertexArrays(1, &VAO)
    gl.BindVertexArray(VAO)
    return VAO
}

// Calculates the normal direction via 3 points in order as Vec3s, and returns it as a Vec3.
func TriangleNormalCalc(p1, p2, p3 mgl32.Vec3) mgl32.Vec3 {
    vecU := p2.Sub(p1)
    vecV := p3.Sub(p1)
    return vecU.Cross(vecV).Normalize()
}

// Makes a new uniform buffer on binding (n) set to the type and size of the (data)
//
// Returns the uint32 ID of the created buffer
func CreateNewUniformBuffer[T mgl64.Mat4 | mgl64.Vec3](data []T, n int) uint32 {
    UBO := GenBindBuffers(gl.UNIFORM_BUFFER)
    BufferData(gl.UNIFORM_BUFFER, make([]float32, int(math.Round(float64(len(data)*len(data[0])) * 2))), gl.DYNAMIC_DRAW)
    gl.BindBufferBase(gl.UNIFORM_BUFFER, uint32(n), UBO)
    return UBO
}

// Sends a generic slice to a uniform buffer (UBO) for use in a shader block.
func BindBufferSubData[T mgl32.Mat4 | mgl32.Vec3](data []T, buffer uint32) {
    switch data := any(data).(type) {
    case []mgl32.Mat4:
        for i := range data {
            v := [16]float32(data[i])
            gl.BindBuffer(gl.UNIFORM_BUFFER, buffer)
            gl.BufferSubData(gl.UNIFORM_BUFFER, i*4*16, 4*16, unsafe.Pointer(&v))
            gl.BindBuffer(gl.UNIFORM_BUFFER, 0)
        }
    case []mgl32.Vec3:
        for i := range data {
            var v []float32
            for j := range data[i] {
                v = append(v, float32(data[i][j]))
            }
            v = append(v, 0.0) // glsl only takes in even values 12 bits needs to padded to 16
            gl.BindBuffer(gl.UNIFORM_BUFFER, buffer)
            gl.BufferSubData(gl.UNIFORM_BUFFER, i*4*4, 4*4, gl.Ptr(v))
            gl.BindBuffer(gl.UNIFORM_BUFFER, 0)
        }

    default:
        panic("BufferData: unsupported type")
    }
}

// Sets the (data) for a given uniform buffer (UBOn)
func SetUBO[T mgl64.Mat4 | mgl64.Vec3](data []T, UBOn uint32) {
    switch data := any(data).(type) {
    case []mgl64.Mat4:
        BindBufferSubData(ghf.Mgl64to32Slice(data).([]mgl32.Mat4), UBOn)
    case []mgl64.Vec3:
        BindBufferSubData(ghf.Mgl64to32Slice(data).([]mgl32.Vec3), UBOn)
    default:
        panic("SetUBO: unsupported type")
    }
}
