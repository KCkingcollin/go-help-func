package ghf

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

func PrintVersionGL() {
    version := gl.GoStr(gl.GetString(gl.VERSION))
    fmt.Println("OpenGL Version", version)
}

func LoadFile(file string) string {
    data, err := os.ReadFile(file)
    if err != nil {
        fmt.Println(err)
    }
    return string(data)
}
func CreateProgram(vertexShader, fragmentShader uint32) (uint32, error) {
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
    return ProgramID, nil
}

func CreateVertexShader(ShaderSource string) (uint32, error) {
    return createShader(ShaderSource,  gl.VERTEX_SHADER)
}

func CreateFragmentShader(ShaderSource string) (uint32, error) {
    return createShader(ShaderSource,  gl.FRAGMENT_SHADER)
}

func createShader(ShaderSource string, ShaderType uint32) (uint32, error) {
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
