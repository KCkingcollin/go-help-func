// GL Helper Functions
package glf

import (
	"errors"
	"fmt"
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

func BufferData[T any](target uint32, data []T, usage uint32) {
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

