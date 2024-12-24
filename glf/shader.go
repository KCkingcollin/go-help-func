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

// Creates a new shader program with a path to a glsl vertex shader, and fragment shader source, in that order
//
// This is only meant to be used at start up
//  - Returns a pointer to the ShaderInfo struct or nil if error
//  - Returns a error or nil if none
func NewShaderProgram(vertexPath, fragmentPath string) (*ShaderInfo, error) {
    id, err := CreateProgram(vertexPath, fragmentPath)
    if err != nil {
        return nil, err
    }
    result := &ShaderInfo{id, vertexPath, fragmentPath, GetModifiedTime(vertexPath), GetModifiedTime(fragmentPath)}
    loadedShaders[id] = result
    return result, nil
}

// Unused and will be deleted and deprecated soon.
//
// Originally used to set float32 uniform for the shader, but this is NOT forward compatible with vulkan, and is not recommended. 
// Use BindBufferSubData instead with a uniform block in the shader.
func (shader *ShaderInfo) SetFloat(name string, f float32) {
    location := gl.GetUniformLocation(shader.id, gl.Str(name + "\x00"))
    gl.Uniform1f(location, f)
}

// Sends a generic slice to a uniform buffer (UBO) binding number "n" for use in a shader block.
func BindBufferSubData[T mgl32.Mat4 | mgl32.Vec3](data []T, UBO, n uint32) {
    switch data := any(data).(type) {
    case []mgl32.Mat4:
        for i := range data {
            v := [16]float32(data[i])
            gl.BindBuffer(gl.UNIFORM_BUFFER, UBO)
            gl.BufferSubData(gl.UNIFORM_BUFFER, i*4*16, 4*16, unsafe.Pointer(&v))
            gl.BindBuffer(gl.UNIFORM_BUFFER, n)
        }
    case []mgl32.Vec3:
        for i := range data {
            v := [3]float32(data[i])
            gl.BindBuffer(gl.UNIFORM_BUFFER, UBO)
            gl.BufferSubData(gl.UNIFORM_BUFFER, i*4*3, 4*3, unsafe.Pointer(&v))
            gl.BindBuffer(gl.UNIFORM_BUFFER, n)
        }
    default:
        panic("unsupported type for BufferData")
    }
}

// Simply uses the given shader program 
func (shader *ShaderInfo) Use() {
    gl.UseProgram(shader.id)
}

// Gets the modified time of a file path and returns the time as a time.Time
func GetModifiedTime(filePath string) time.Time {
        file, err := os.Stat(filePath)
        if err != nil && Verbose {
            fmt.Println(err)
        }
        return file.ModTime()
}

// Checks to see if any of the loaded shaders have been modified, and if so recreates the program for that shader.
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

