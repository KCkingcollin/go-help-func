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
    shaderProgram := gl.CreateProgram()
    gl.AttachShader(shaderProgram, vertexShader)
    gl.AttachShader(shaderProgram, fragmentShader)
    gl.LinkProgram(shaderProgram)
    var success int32
    gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
    if success == gl.FALSE {
        var logLength int32
        gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
        return shaderProgram, errors.New(log)
    }
    return shaderProgram, nil
}

func CreateVertexShader(ShaderSource string) (uint32, error) {
    return createShader(ShaderSource,  gl.VERTEX_SHADER)
}

func CreateFragmentShader(ShaderSource string) (uint32, error) {
    return createShader(ShaderSource,  gl.FRAGMENT_SHADER)
}

func createShader(ShaderSource string, ShaderType uint32) (uint32, error) {
    Shader := gl.CreateShader(ShaderType)
    csource, free := gl.Strs(ShaderSource)
    gl.ShaderSource(Shader, 1, csource, nil)
    free()
    gl.CompileShader(Shader)
    var status int32
    gl.GetShaderiv(Shader, gl.COMPILE_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetShaderiv(Shader, gl.INFO_LOG_LENGTH, &logLength)
        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetShaderInfoLog(Shader, logLength, nil, gl.Str(log))
        return Shader, errors.New(log)
    }
    return Shader, nil 
}
