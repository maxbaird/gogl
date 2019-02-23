package gogl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

//GetVersion ...
func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}
