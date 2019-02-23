package gogl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"strings"
)

//GetVersion ...
func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

//MakeShader ...
func MakeShader(shaderSource string, shaderType uint32) uint32 {
	shaderID := gl.CreateShader(shaderType)
	shaderSource = shaderSource + "\x00"
	csource, free := gl.Strs(shaderSource)
	gl.ShaderSource(shaderID, 1, csource, nil)
	free()

	gl.CompileShader(shaderID)

	var status int32
	gl.GetShaderiv(shaderID, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderID, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderID, logLength, nil, gl.Str(log))
		panic("Failed to compile shader: \n" + log)
	}
	return shaderID
}
