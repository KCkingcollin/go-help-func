package ghf

import {
	"github.com/go-gl/gl/v4.6-core/gl"
}

func GetVersionGL() {
    version := gl.GoStr(gl.GetString(gl.VERSION))
    fmt.Println("OpenGL Version", version)
}
