package common

import "github.com/go-gl/glfw3/v3.1/glfw"

// horizontal angle : toward -Z
var horizontalAngle = 3.14

// vertical angle : 0, look at the horizon
var verticalAngle = 0.0

var mouseSpeed = 0.005

var lastTime float64

func ComputeMatricesFromInputs() {

	// glfwGetTime should only be called once, the first time this function is called
	if lastTime == 0.0 {
		lastTime = glfw.GetTime()
	}

	// Compute the time difference between current and last frame
	currentTime := glfw.GetTime()
	deltaTime := currentTime - lastTime

	window := glfw.GetCurrentContext()

	// Get mouse position
	xpos, ypos := window.GetCursorPos()

	// Reset mouse position for next frame
	windowWidth, windowHeight := window.GetSize()
	window.SetCursorPos(windowWidth/2, windowHeight/2)

	// Compute new orientation
	horizontalAngle += mouseSpeed * deltaTime * float32(windowWidth/2-xpos)
	verticalAngle += mouseSpeed * deltaTime * float32(windowHeight/2-ypos)
}
