package main

import (
	"log"

	"fmt"

	"runtime"

	"sync"

	"github.com/fapiko/go-learn-gl/opengl-tutorial/common"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw3/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Tutorial 11 - 2d text ported from
// http://www.opengl-tutorial.org/intermediate-tutorials/tutorial-11-2d-text/
func main() {

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	// OpenGL needs to be locked to a thread, so make a goroutine that calls runtime.LockOSThread()
	go renderThread(waitGroup)

	waitGroup.Wait()

}

func renderThread(waitGroup *sync.WaitGroup) {

	defer waitGroup.Done()
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		panic("Failed to initialize GLFW")
	}

	defer glfw.Terminate()

	glfw.WindowHint(glfw.Samples, 4)

	// Drawing the triangle threw an error with OpenGL 3.3, downgrading to 2.1 seemed to solve it
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	// Open a window and create its OpenGL context
	window, err := glfw.CreateWindow(1024, 768, "Tutorial 11", nil, nil)
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

	// Hide the mouse and enable unlimited movement
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// Dark blue background
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)

	// Enable depth test
	gl.Enable(gl.DEPTH_TEST)

	// Accept fragment if it is closer to the camera than the former one
	gl.DepthFunc(gl.LESS)

	// Disable backface culling
	gl.Disable(gl.CULL_FACE)

	// Create and compile our GLSL program from the shaders
	programId := common.LoadShaders("StandardShading.vertexshader", "StandardShading.fragmentshader")
	defer gl.DeleteProgram(programId)

	// Get a handle for our "MVP" uniform
	matrixId := gl.GetUniformLocation(programId, gl.Str("MVP\x00"))
	viewMatrixId := gl.GetUniformLocation(programId, gl.Str("V\x00"))
	modelMatrixId := gl.GetUniformLocation(programId, gl.Str("M\x00"))

	// Get a handle for our buffers
	vertexPositionModelspaceId := uint32(gl.GetAttribLocation(programId, gl.Str("vertexPosition_modelspace\x00")))

	vertices, uvs, normals, err := common.LoadObj("suzanne.obj")
	if err != nil {
		log.Panic(err)
	}

	indices, indexedVertices, indexedUvs, indexedNormals := common.IndexVBO(vertices, uvs, normals)

	var vertexBuffer uint32
	gl.GenBuffers(1, &vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(indexedVertices)*4*3, gl.Ptr(indexedVertices), gl.STATIC_DRAW)
	defer gl.DeleteBuffers(1, &vertexBuffer)

	var normalBuffer uint32
	gl.GenBuffers(1, &normalBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, normalBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(indexedNormals)*4*3, gl.Ptr(indexedNormals), gl.STATIC_DRAW)
	defer gl.DeleteBuffers(1, &normalBuffer)

	textureId, err := common.LoadDDS("uvmap.DDS")
	if err != nil {
		panic(err)
	}
	defer gl.DeleteTextures(1, &textureId)

	var uvBuffer uint32
	gl.GenBuffers(1, &uvBuffer)
	gl.BindBuffer(gl.TEXTURE_BUFFER, uvBuffer)
	gl.BufferData(gl.TEXTURE_BUFFER, len(indexedUvs)*4*2, gl.Ptr(indexedUvs), gl.STATIC_DRAW)
	defer gl.DeleteBuffers(1, &uvBuffer)

	// Generate a buffer for the indices as well
	var elementBuffer uint32
	gl.GenBuffers(1, &elementBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, elementBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
	defer gl.DeleteBuffers(1, &elementBuffer)

	// Set the mouse at the center of the screen
	glfw.PollEvents()
	windowWidth, windowHeight := window.GetSize()
	window.SetCursorPos(float64(windowWidth/2), float64(windowHeight/2))

	lightId := gl.GetUniformLocation(programId, gl.Str("LightPosition_worldspace\x00"))

	// Initialize our little text library with the Holstein font
	common.InitText2d("Holstein.DDS")
	defer common.CleanupText2d()

	// For speed computation
	lastTime := glfw.GetTime()
	var nbFrames int

	for window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose() {

		// Measure speed
		currentTime := glfw.GetTime()
		nbFrames++
		if currentTime-lastTime >= 1 {
			fmt.Printf("\r%f ms/frame", 1000.0/float64(nbFrames))
			nbFrames = 0
			lastTime = currentTime
		}

		// Clear the screen
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Use our shader
		gl.UseProgram(programId)

		// Compute the MVP matrix from keyboard and mouse input
		common.ComputeMatricesFromInputs()
		projection := common.GetProjectionMatrix()
		view := common.GetViewMatrix()
		model := mgl32.Ident4()
		mvp := projection.Mul4(view).Mul4(model)

		// Send our transformation to the currently bound shader, in the "MVP" uniform
		gl.UniformMatrix4fv(matrixId, 1, false, &mvp[0])
		gl.UniformMatrix4fv(modelMatrixId, 1, false, &model[0])
		gl.UniformMatrix4fv(viewMatrixId, 1, false, &view[0])

		lightPos := &mgl32.Vec3{4, 4, 4}
		gl.Uniform3f(lightId, lightPos.X(), lightPos.Y(), lightPos.Z())

		// Bind our texture in Texture Unit 0
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, uint32(textureId))

		// Set our "myTextureSampler" sampler to user Texture Unit 0
		gl.Uniform1i(int32(textureId), 0)

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

		// 2nd attribute buffer : uvs
		gl.EnableVertexAttribArray(1)
		gl.BindBuffer(gl.ARRAY_BUFFER, uvBuffer)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, nil)

		// 3rd attribute buffer : normals
		gl.EnableVertexAttribArray(2)
		gl.BindBuffer(gl.ARRAY_BUFFER, normalBuffer)
		gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 0, nil)

		// Index buffer
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, elementBuffer)

		// Draw the triangle !
		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil)

		text := fmt.Sprintf("%.2f sec", glfw.GetTime())
		common.PrintText2D(text, 10, 500, 60)

		gl.DisableVertexAttribArray(vertexPositionModelspaceId)
		gl.DisableVertexAttribArray(1)
		gl.DisableVertexAttribArray(2)

		// Swap buffers
		window.SwapBuffers()
		glfw.PollEvents()

	}

}
