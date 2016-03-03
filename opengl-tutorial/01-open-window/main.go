package main

import (
	"runtime"
	"sync"

	"github.com/go-gl/glfw/v3.1/glfw"
)

// Tutorial 01 - Creating a Window ported from
// http://www.opengl-tutorial.org/beginners-tutorials/tutorial-1-opening-a-window/
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

	err := glfw.Init()
	if err != nil {
		panic("Failed to initialize GLFW")
	}

	defer glfw.Terminate()

	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(1024, 768, "Tutorial 01", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	for window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose() {

		window.SwapBuffers()
		glfw.PollEvents()

	}

}
