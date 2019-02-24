package gogl

import (
	"errors"
	//"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"io/ioutil"
	//"os"
	"strings"
	"time"
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

type programInfo struct {
	vertPath string
	fragPath string
	modified time.Time
}

//GetVersion ...
func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

var loadedShaders []programInfo

//CheckShadersForChanges ...
func CheckShadersForChanges() {
	/*
		for _, shaderInfo := range loadedShaders {
			file, err := os.Stat(shaderInfo.path)

			if err != nil {
				panic(err)
			}

			modTime := file.ModTime()

			if !modTime.Equal(shaderInfo.modified) {
				fmt.Println("Shader modified")
			}
		}
	*/
}

//LoadShader ...
func LoadShader(path string, shaderType uint32) (ShaderID, error) {
	shaderFile, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}

	shaderFileStr := string(shaderFile)
	shaderID, err := CreateShader(shaderFileStr, shaderType)

	if err != nil {
		return 0, err
	}
	return shaderID, nil
}

//CreateShader ...
func CreateShader(shaderSource string, shaderType uint32) (ShaderID, error) {
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
		//fmt.Printf("Failed to compile shader: %s\n", log)
		return 0, errors.New("Failed to compile shader: " + log)
	}
	return ShaderID(shaderID), nil
}

//CreateProgram ...
func CreateProgram(vertPath string, fragPath string) (ProgramID, error) {
	vert, err := LoadShader(vertPath, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	frag, err := LoadShader(fragPath, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

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
		return 0, errors.New("Failed to link program: \n" + log)
	}
	gl.DeleteShader(uint32(vert))
	gl.DeleteShader(uint32(frag))

	//TODO finish hotloading shaders
	/*
		file, err := os.Stat(path)

		if err != nil {
			panic(err)
		}

		modTime := file.ModTime()
		loadedShaders = append(loadedShaders, shaderInfo{path, modTime})
	*/

	return ProgramID(shaderProgram), nil
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
