package ghf

import (
	"fmt"
	"os"
	"time"
)

type Shader struct {
    id              uint32
    vertexPath      string
    fragmentPath    string
}

type shaderInfo struct {
    path string
    modified time.Time 
}

var loadedShaders []shaderInfo
func NewShader(vertexPath, fragmentPath string) (Shader, error) {
    id, err := CreateProgram(vertexPath, fragmentPath)
    return Shader{id, vertexPath, fragmentPath}, err
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

