package ghf

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type shaderInfo struct {
    path string
    modified time.Time 
}

var loadedShaders []shaderInfo

func ResetLoadedShaders() {
    loadedShaders = nil
}

func CheckShadersforChanges() bool {
    for _, shaderInfo := range loadedShaders {
        file, err := os.Stat(shaderInfo.path)
        if err != nil {
            fmt.Println(err)
            return false
        }
        modTime := file.ModTime()
        if !modTime.Equal(shaderInfo.modified) {
            fmt.Println("Reloading shader: " + shaderInfo.path)
            return true
        }
    }
    return false
}

func PrintVersionGL() {
    version := gl.GoStr(gl.GetString(gl.VERSION))
    fmt.Println("OpenGL Version", version)
}

func LoadFile(path string) string {
    data, err := os.ReadFile(path)
    if err != nil {
        fmt.Println(err)
    }
    file, err := os.Stat(path)
    if err != nil {
        fmt.Println(err)
    }
    modTime := file.ModTime()
    loadedShaders = append(loadedShaders, shaderInfo{path, modTime})
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
    return CreateShader(ShaderSource,  gl.VERTEX_SHADER)
}

func CreateFragmentShader(ShaderSource string) (uint32, error) {
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

func BufferDataFloat(target uint32, data []float32, usage uint32) {
    gl.BufferData(target, len(data)*4, gl.Ptr(data), usage)
}

func BufferDataInt(target uint32, data []uint32, usage uint32) {
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

