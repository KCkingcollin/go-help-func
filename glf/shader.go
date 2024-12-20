// Shader GL Helper Functions
package glf

import (
	"fmt"
	"os"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type ShaderInfo struct {
    id              uint32
    vertexPath      string
    fragmentPath    string
    vertModified    time.Time
    fragModified    time.Time
}

var loadedShaders = make(map[uint32]*ShaderInfo)

func NewShaderProgram(vertexPath, fragmentPath string) (*ShaderInfo, error) {
    id, err := CreateProgram(vertexPath, fragmentPath)
    if err != nil {
        return nil, err
    }
    result := &ShaderInfo{id, vertexPath, fragmentPath, GetModifiedTime(vertexPath), GetModifiedTime(fragmentPath)}
    loadedShaders[id] = result
    return result, nil
}

func (shader *ShaderInfo) SetFloat(name string, f float32) {
    location := gl.GetUniformLocation(shader.id, gl.Str(name + "\x00"))
    gl.Uniform1f(location, f)
}

func BindBufferSubDataMat4(split []mgl32.Mat4, UBOn uint32) {
    for i := range split {
        v := [16]float32(split[i])
        gl.BindBuffer(gl.UNIFORM_BUFFER, UBOn)
        gl.BufferSubData(gl.UNIFORM_BUFFER, i*4*16, 4*16, unsafe.Pointer(&v))
        gl.BindBuffer(gl.UNIFORM_BUFFER, 1)
    }
}

func (shader *ShaderInfo) Use() {
    gl.UseProgram(shader.id)
}

func GetModifiedTime(filePath string) time.Time {
        file, err := os.Stat(filePath)
        if err != nil && Verbose {
            fmt.Println(err)
        }
        return file.ModTime()
}

func CheckShadersforChanges() {
    for _, shader := range loadedShaders {
        vertexModTime := GetModifiedTime(shader.vertexPath)
        fragmentModTime := GetModifiedTime(shader.fragmentPath)
        if !vertexModTime.Equal(shader.vertModified) || 
        !fragmentModTime.Equal(shader.fragModified) {
            fmt.Println("Reloading vertex and fragment shader: \n" + shader.vertexPath + "\n" + shader.fragmentPath)
            id, err := CreateProgram(shader.vertexPath, shader.fragmentPath)
            if err != nil {
                if Verbose {
                    fmt.Printf("Could not relink shader, %s \n", err)
                }
                shader.vertModified = GetModifiedTime(shader.vertexPath)
                shader.fragModified = GetModifiedTime(shader.fragmentPath)
            } else {
                if Verbose {
                    fmt.Println("Relinked shader")
                }
                gl.DeleteProgram(shader.id)
                shader.id = id
                shader.vertModified = GetModifiedTime(shader.vertexPath)
                shader.fragModified = GetModifiedTime(shader.fragmentPath)
            }
        }
    }
}

