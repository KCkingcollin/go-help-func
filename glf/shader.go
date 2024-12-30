// Shader GL Helper Functions
package glf

import (
	"fmt"
	"time"

	"github.com/KCkingcollin/go-help-func/ghf"
	"github.com/go-gl/gl/v4.6-core/gl"
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
    result := &ShaderInfo{id, vertexPath, fragmentPath, ghf.GetModifiedTime(vertexPath), ghf.GetModifiedTime(fragmentPath)}
    loadedShaders[id] = result
    return result, nil
}

// Simply uses the given shader program 
func (shader *ShaderInfo) Use() {
    gl.UseProgram(shader.id)
}

// Checks to see if any of the loaded shaders have been modified, and if so recreates the program for that shader.
func CheckShadersforChanges() {
    for _, shader := range loadedShaders {
        vertexModTime := ghf.GetModifiedTime(shader.vertexPath)
        fragmentModTime := ghf.GetModifiedTime(shader.fragmentPath)
        if !vertexModTime.Equal(shader.vertModified) || 
        !fragmentModTime.Equal(shader.fragModified) {
            fmt.Println("Reloading vertex and fragment shader: \n" + shader.vertexPath + "\n" + shader.fragmentPath)
            id, err := CreateProgram(shader.vertexPath, shader.fragmentPath)
            if err != nil {
                if Verbose {
                    fmt.Printf("Could not relink shader, %s \n", err)
                }
                shader.vertModified = ghf.GetModifiedTime(shader.vertexPath)
                shader.fragModified = ghf.GetModifiedTime(shader.fragmentPath)
            } else {
                if Verbose {
                    fmt.Println("Relinked shader")
                }
                gl.DeleteProgram(shader.id)
                shader.id = id
                shader.vertModified = ghf.GetModifiedTime(shader.vertexPath)
                shader.fragModified = ghf.GetModifiedTime(shader.fragmentPath)
            }
        }
    }
}

