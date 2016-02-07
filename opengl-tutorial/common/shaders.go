package common

import (
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/go-gl/gl/all-core/gl"
)

func LoadShaders(vertexFilePath string, fragmentFilePath string) uint32 {

	vertexShaderId := gl.CreateShader(gl.VERTEX_SHADER)
	defer gl.DeleteShader(vertexShaderId)

	fragmentShaderId := gl.CreateShader(gl.FRAGMENT_SHADER)
	defer gl.DeleteShader(fragmentShaderId)

	// Read the Vertex Shader code from the file
	vertexShaderCode, err := ioutil.ReadFile(vertexFilePath)
	if err != nil {

		log.Errorf("Impossible to open %s. Are you in the right directory ? Don't forget to read the FAQ !",
			vertexFilePath)
		return 0

	}

	// Read the Fragment Shader code from the file
	fragmentShaderCode, err := ioutil.ReadFile(fragmentFilePath)
	if err != nil {

		log.Errorf("Impossible to open %s. Are you in the right directory ? Don't forget to read the FAQ !",
			fragmentFilePath)
		return 0

	}

	// Compile Vertex Shader
	log.Infof("Compiling shader : %s", vertexFilePath)
	vertexSourcePointer := gl.Str(nullTerminatedString(vertexShaderCode))

	gl.ShaderSource(vertexShaderId, 1, &vertexSourcePointer, nil)
	gl.CompileShader(vertexShaderId)

	// Check Vertex Shader
	var infoLogLength int32
	var result int32

	gl.GetShaderiv(vertexShaderId, gl.COMPILE_STATUS, &result)
	gl.GetShaderiv(vertexShaderId, gl.INFO_LOG_LENGTH, &infoLogLength)
	if result != gl.TRUE {

		vertexShaderErrorMessage := strings.Repeat("\x00", int(infoLogLength))

		var messageLength int32
		gl.GetShaderInfoLog(vertexShaderId, infoLogLength, &messageLength, gl.Str(vertexShaderErrorMessage))

		log.Info(vertexShaderErrorMessage)

	}

	// Compile Fragment Shader
	log.Infof("Compiling shader : %s", fragmentFilePath)
	fragmentSourcePointer := gl.Str(nullTerminatedString(fragmentShaderCode))

	gl.ShaderSource(fragmentShaderId, 1, &fragmentSourcePointer, nil)
	gl.CompileShader(fragmentShaderId)

	gl.GetShaderiv(fragmentShaderId, gl.COMPILE_STATUS, &result)
	gl.GetShaderiv(fragmentShaderId, gl.INFO_LOG_LENGTH, &infoLogLength)
	if result != gl.TRUE {

		fragmentShaderErrorMessage := strings.Repeat("\x00", int(infoLogLength))

		var messageLength int32
		gl.GetShaderInfoLog(vertexShaderId, infoLogLength, &messageLength, gl.Str(fragmentShaderErrorMessage))

		log.Info(fragmentShaderErrorMessage)

	}

	// Link the program
	log.Printf("Linking program")
	programId := gl.CreateProgram()
	gl.AttachShader(programId, vertexShaderId)
	gl.AttachShader(programId, fragmentShaderId)
	gl.LinkProgram(programId)

	gl.GetProgramiv(programId, gl.LINK_STATUS, &result)
	gl.GetProgramiv(programId, gl.INFO_LOG_LENGTH, &infoLogLength)
	if result != gl.TRUE {

		programErrorMessage := strings.Repeat("\x00", int(infoLogLength))

		var messageLength int32
		gl.GetProgramInfoLog(programId, infoLogLength, &messageLength, gl.Str(programErrorMessage))

		log.Info(programErrorMessage)

	}

	defer gl.DetachShader(programId, vertexShaderId)
	defer gl.DetachShader(programId, fragmentShaderId)

	return programId

}
