package gogl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"strings"
)

type (
	//ShaderID ...
	ShaderID uint32
	//ProgramID ...
	ProgramID uint32
	//VAOID ...
	VAOID uint32
	//VBOID ...
	VBOID uint32
)

//GetVersion ...
func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

//CreateShader ...
func CreateShader(shaderSource string, shaderType uint32) ShaderID {
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
	return ShaderID(shaderID)
}

//CreateProgram ...
func CreateProgram(vert ShaderID, frag ShaderID) ProgramID {
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, uint32(vert))
	gl.AttachShader(shaderProgram, uint32(frag))
	gl.LinkProgram(shaderProgram)

	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)

	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
		panic("Failed to link program: \n" + log)
	}
	gl.DeleteShader(uint32(vert))
	gl.DeleteShader(uint32(frag))

	return ProgramID(shaderProgram)
}
