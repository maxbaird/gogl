package gogl

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"os"
	"time"
)

//Shader ...
type Shader struct {
	id               ProgramID
	vertexPath       string
	fragmentPath     string
	vertexModified   time.Time
	fragmentModified time.Time
}

//NewShader ...
func NewShader(vertexPath string, fragmentPath string) (*Shader, error) {
	id, err := CreateProgram(vertexPath, fragmentPath)

	if err != nil {
		return nil, err
	}

	vertexModTime, err := getModifiedTime(vertexPath)
	if err != nil {
		return nil, err
	}

	fragmentModTime, err := getModifiedTime(fragmentPath)
	if err != nil {
		return nil, err
	}

	result := &Shader{id,
		vertexPath,
		fragmentPath,
		vertexModTime,
		fragmentModTime}

	return result, nil

}

//Use ...
func (shader *Shader) Use() {
	UseProgram(shader.id)
}

func getModifiedTime(filePath string) (time.Time, error) {
	file, err := os.Stat(filePath)

	if err != nil {
		return time.Time{}, err
	}

	return file.ModTime(), nil
}

//CheckShaderForChanges ...
func (shader *Shader) CheckShaderForChanges() error {

	vertexModTime, err := getModifiedTime(shader.vertexPath)
	if err != nil {
		return err
	}

	fragmentModTime, err := getModifiedTime(shader.fragmentPath)
	if err != nil {
		return err
	}

	if !vertexModTime.Equal(shader.vertexModified) || !fragmentModTime.Equal(shader.fragmentModified) {
		id, err := CreateProgram(shader.vertexPath, shader.fragmentPath)
		if err != nil {
			fmt.Println(err)
		} else {
			gl.DeleteProgram(uint32(shader.id))
			shader.id = id
		}
	}
	return nil
}
