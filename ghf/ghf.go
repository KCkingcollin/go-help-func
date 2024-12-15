package ghf

import {
	"github.com/go-gl/gl/v4.6-core/gl"
}

func GetVersionGL() {
    return version := gl.GoStr(gl.GetString(gl.VERSION))
}
