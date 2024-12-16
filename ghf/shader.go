package ghf

import (
	"fmt"
	"os"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type ShaderInfo struct {
    id              uint32
    vertexPath      string
    fragmentPath    string
    vertModified        time.Time
    fragModified        time.Time
}

var loadedShaders = make(map[uint32]*ShaderInfo)

func NewShader(vertexPath, fragmentPath string) (ShaderInfo, error) {
    id, err := CreateProgram(vertexPath, fragmentPath)
    return ShaderInfo{id, vertexPath, fragmentPath, GetModifiedTime(vertexPath), GetModifiedTime(fragmentPath)}, err
}

func GetModifiedTime(filePath string) time.Time {
        file, err := os.Stat(filePath)
        if err != nil {
            fmt.Println(err)
        }
        return file.ModTime()
}

func CheckShadersforChanges() {
    for _, ShaderInfo := range loadedShaders {
        vertexModTime := GetModifiedTime(ShaderInfo.vertexPath)
        fragmentModTime := GetModifiedTime(ShaderInfo.fragmentPath)
        if !vertexModTime.Equal(ShaderInfo.vertModified) || 
        !fragmentModTime.Equal(ShaderInfo.fragModified) {
            fmt.Println("Reloading vertex and fragment shader: \n" + ShaderInfo.vertexPath + "\n" + ShaderInfo.fragmentPath)
            id, err := CreateProgram(ShaderInfo.vertexPath, ShaderInfo.fragmentPath)
            if err != nil {
                fmt.Printf("Could not relink shader: \n %s", err)
            } else {
                fmt.Printf("Relinked shader")
                gl.DeleteProgram(ShaderInfo.id)
                ShaderInfo.id = id
            }
        }
    }
}

