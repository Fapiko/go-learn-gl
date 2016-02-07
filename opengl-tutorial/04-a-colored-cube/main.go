package main

import (
	"io/ioutil"

	log "github.com/Sirupsen/logrus"

	"strings"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw3/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Tutorial 04 - A colored cube ported from
// http://www.opengl-tutorial.org/beginners-tutorials/tutorial-2-the-first-triangle/
func main() {

	if err := glfw.Init(); err != nil {
		panic("Failed to initialize GLFW")
	}

	defer glfw.Terminate()

	glfw.WindowHint(glfw.Samples, 4)

	// Drawing the triangle threw an error with OpenGL 3.3, downgrading to 2.1 seemed to solve it
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	// Open a window and create its OpenGL context
	window, err := glfw.CreateWindow(1024, 768, "Tutorial 04", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()

	// Initialize OpenGL - Go bindings use Glow and now Glew
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// Ensure we can capture the escape key being pressed below
	window.SetInputMode(glfw.StickyKeysMode, gl.TRUE)

	// Dark blue background
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)

	// Enable depth test
	gl.Enable(gl.DEPTH_TEST)

	// Accept fragment if it is closer to the camera than the former one
	gl.DepthFunc(gl.LESS)

	// Create and compile our GLSL program from the shaders
	programId := loadShaders("SimpleTransform.vertexshader", "SingleColor.fragmentshader")
	defer gl.DeleteProgram(programId)

	// Get a handle for our "MVP" uniform
	matrixId := gl.GetUniformLocation(programId, gl.Str("MVP\x00"))

	// Projection matrix : 45° Field of View, 4:3 ratio, display range : 0.1 unit <-> 100 units
	projection := mgl32.Perspective(45.0, 4.0/3.0, 0.1, 100.0)
	// Or, for an ortho camera :
	// projection := mgl32.Ortho(-10.0, 10.0, -10.0, 10.0, 0.0, 100.0) // In world coordinates

	// Camera matrix
	view := mgl32.LookAt(4.0, 3.0, -3.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0)

	// Model matrix : an identity matrix (model will be at the origin)
	model := mgl32.Ident4()

	mvp := projection.Mul4(view).Mul4(model)

	// Get a handle for our buffers
	vertexPositionModelspaceId := uint32(gl.GetAttribLocation(programId, gl.Str("vertexPosition_modelspace\x00")))

	vertexBufferData := []float32{
		-1.0, -1.0, -1.0, // triangle 1 : begin
		-1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0, // triangle 1 : end
		1.0, 1.0, -1.0, // triangle 2 : begin
		-1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0, // triangle 2 : end
		1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
	}

	colorBufferData := []float32{
		0.583, 0.771, 0.014,
		0.609, 0.115, 0.436,
		0.327, 0.483, 0.844,
		0.822, 0.569, 0.201,
		0.435, 0.602, 0.223,
		0.310, 0.747, 0.185,
		0.597, 0.770, 0.761,
		0.559, 0.436, 0.730,
		0.359, 0.583, 0.152,
		0.483, 0.596, 0.789,
		0.559, 0.861, 0.639,
		0.195, 0.548, 0.859,
		0.014, 0.184, 0.576,
		0.771, 0.328, 0.970,
		0.406, 0.615, 0.116,
		0.676, 0.977, 0.133,
		0.971, 0.572, 0.833,
		0.140, 0.616, 0.489,
		0.997, 0.513, 0.064,
		0.945, 0.719, 0.592,
		0.543, 0.021, 0.978,
		0.279, 0.317, 0.505,
		0.167, 0.620, 0.077,
		0.347, 0.857, 0.137,
		0.055, 0.953, 0.042,
		0.714, 0.505, 0.345,
		0.783, 0.290, 0.734,
		0.722, 0.645, 0.174,
		0.302, 0.455, 0.848,
		0.225, 0.587, 0.040,
		0.517, 0.713, 0.338,
		0.053, 0.959, 0.120,
		0.393, 0.621, 0.362,
		0.673, 0.211, 0.457,
		0.820, 0.883, 0.371,
		0.982, 0.099, 0.879,
	}

	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexBufferData)*4, gl.Ptr(vertexBufferData), gl.STATIC_DRAW)
	defer gl.DeleteBuffers(1, &vertexBuffer)

	var colorBuffer uint32
	gl.GenBuffers(1, &colorBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(colorBufferData)*4, gl.Ptr(colorBufferData), gl.STATIC_DRAW)
	defer gl.DeleteBuffers(1, &colorBuffer)

	for window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose() {

		// Clear the screen
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Use our shader
		gl.UseProgram(programId)

		// Send our transformation to the currently bound shader, in the "MVP" uniform
		gl.UniformMatrix4fv(matrixId, 1, false, &mvp[0])

		// 1st attribute buffer : vertices
		gl.EnableVertexAttribArray(vertexPositionModelspaceId)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
		gl.VertexAttribPointer(
			vertexPositionModelspaceId, // The attribute we want to configure
			3,               // size
			gl.FLOAT,        // type
			false,           // normalized?
			0,               // stride
			gl.PtrOffset(0)) // array buffer offset

		// Draw the triangle !
		gl.DrawArrays(gl.TRIANGLES, 0, 12*3) // 12*3 indices starting at 0 -> 12 triangles
		gl.DisableVertexAttribArray(vertexPositionModelspaceId)

		// 2nd attribute buffer : colors
		gl.EnableVertexAttribArray(1)
		gl.BindBuffer(gl.ARRAY_BUFFER, colorBuffer)
		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, nil)

		// Swap buffers
		window.SwapBuffers()
		glfw.PollEvents()

	}

}

func loadShaders(vertexFilePath string, fragmentFilePath string) uint32 {

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

func nullTerminatedString(source []byte) string {
	return string(append(source, make([]byte, 1)...))
}
