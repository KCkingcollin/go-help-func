// GL Helper Functions
package glf

import (
	"errors"
	"fmt"
	"image/png"
	"os"
	"strings"

	"github.com/KCkingcollin/go-help-func/ghf"
	"github.com/go-gl/gl/v4.6-core/gl"
)

var Verbose bool = ghf.Verbose

func PrintVersionGL() {
    version := gl.GoStr(gl.GetString(gl.VERSION))
    if Verbose {
        fmt.Println("OpenGL Version", version)
    }
}

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

func GenBindTexture() uint32 {
    var texID uint32
    gl.GenTextures(1, &texID)
    gl.BindTexture(gl.TEXTURE_2D, texID)
    return texID
}

func BindTexture(texID uint32) {
    gl.BindTexture(gl.TEXTURE_2D, texID)
}

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

func CreateVertexShader(shaderFile string) (uint32, error) {
    ShaderSource := ghf.LoadFile(shaderFile) + "\x00"
    return CreateShader(ShaderSource,  gl.VERTEX_SHADER)
}

func CreateFragmentShader(shaderFile string) (uint32, error) {
    ShaderSource := ghf.LoadFile(shaderFile) + "\x00"
    return CreateShader(ShaderSource,  gl.FRAGMENT_SHADER)
}

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

func GenBindBuffers(target uint32) uint32 {
    var buffer uint32
    gl.GenBuffers(1, &buffer)
    gl.BindBuffer(target, buffer)
    return buffer
}

func GenBindArrays() uint32 {
    var VAO uint32
    gl.GenVertexArrays(1, &VAO)
    gl.BindVertexArray(VAO)
    return VAO
}

