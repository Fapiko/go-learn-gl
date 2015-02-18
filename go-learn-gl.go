package main

import (
	"log"
	"os"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glh"
)

func main() {
	if !glfw.Init() {
		log.Fatal("Failed to initialize GLFW")

		os.Exit(-1)
	}

	defer glfw.Terminate()

	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)

	window, error := glfw.CreateWindow(1024, 768, "Tutorial 01", nil, nil)

	if error != nil {
		log.Fatal("Failed to open GLFW window.")
		glfw.Terminate()
		os.Exit(-1)
	}

	window.MakeContextCurrent()

	if gl.Init() != 0 {
		log.Fatal("Failed to init GL")
	}

	window.SetInputMode(glfw.StickyKeys, glfw.True)

	// Draw a triangle
	vertexArray := gl.GenVertexArray()
	vertexArray.Bind()

	triangleVertices := []float32{
		-1, -1, 0,
		1, -1, 0,
		0, 1, 0,
	}
	vertexBuffer := gl.GenBuffer()
	vertexBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, int(glh.Sizeof(gl.FLOAT))*len(triangleVertices), triangleVertices, gl.STATIC_DRAW)

	vertexShaderData := `#version 330 core
layout(location = 0) in vec3 vertexPosition_modelspace;
void main(){
	gl_Position.xyz = vertexPosition_modelspace;
	gl_Position.w = 1.0;
 }`

	fragmentShaderData := `#version 330 core
out vec3 color;

void main(){
    color = vec3(1,0,0);
}`

	log.Println(gl.GetError())

	vertexShader := glh.Shader{gl.VERTEX_SHADER, vertexShaderData}
	fragmentShader := glh.Shader{gl.FRAGMENT_SHADER, fragmentShaderData}

	program := glh.NewProgram(vertexShader, fragmentShader)
	program.Use()

	for window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		var vertexAttrib gl.AttribLocation = 0
		vertexAttrib.EnableArray()
		vertexBuffer.Bind(gl.ARRAY_BUFFER)

		vertexAttrib.AttribPointer(
			3,
			gl.FLOAT,
			false,
			0,
			nil)

		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		vertexAttrib.DisableArray()

		window.SwapBuffers()
		glfw.PollEvents()

		//		log.Println(gl.GetError())
	}
}
