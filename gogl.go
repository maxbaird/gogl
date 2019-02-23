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
func CreateProgram(vertStr string, fragStr string) ProgramID {
	vert := CreateShader(vertStr, gl.VERTEX_SHADER)
	frag := CreateShader(fragStr, gl.FRAGMENT_SHADER)

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

//GenBindBuffer ...
func GenBindBuffer(target uint32) VBOID {
	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(target, VBO)
	return VBOID(VBO)
}

//GenBindVertexArray ...
func GenBindVertexArray() VAOID {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	return VAOID(VAO)
}

//BindVertexArray ...
func BindVertexArray(vaoID VAOID) {
	gl.BindVertexArray(uint32(vaoID))
}

//BufferDataFloat ...
func BufferDataFloat(target uint32, data []float32, usage uint32) {
	gl.BufferData(target, len(data)*4, gl.Ptr(data), usage)
}

//UnbindVertexArray ...
func UnbindVertexArray() {
	gl.BindVertexArray(0)
}

//UseProgram ...
func UseProgram(programID ProgramID) {
	gl.UseProgram(uint32(programID))
}
