package main

import (
	"log"
	"os"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	//	"time"
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

	gl.Init()

	window.SetInputMode(glfw.StickyKeys, glfw.True)

	for window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose() {
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
