package ghf

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

func GetVersionGL() string {
    return gl.GoStr(gl.GetString(gl.VERSION))
}
