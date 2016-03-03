package common

import (
	"math"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var projectionMatrix mgl32.Mat4
var viewMatrix mgl32.Mat4

func GetViewMatrix() mgl32.Mat4 {
	return viewMatrix
}

func GetProjectionMatrix() mgl32.Mat4 {
	return projectionMatrix
}

// Initial position : on +Z
var position = mgl32.Vec3{0.0, 0.0, 5.0}

// horizontal angle : toward -Z
var horizontalAngle = float32(3.14)

// vertical angle : 0, look at the horizon
var verticalAngle = float32(0.0)

// Initial Field of View
var fov = float32(45.0)

var speed = float32(3.0) // 3 units / second
var mouseSpeed = float32(0.005)

var lastTime float64

var zoomAmount float32

func ComputeMatricesFromInputs() {

	window := glfw.GetCurrentContext()

	// glfwGetTime should only be called once, the first time this function is called
	if lastTime == 0.0 {

		lastTime = glfw.GetTime()

		// Add a callback for the scrollwheel to enable zoom functionality
		window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
			fov += .1 * float32(-yoff)
		})

	}

	// Compute the time difference between current and last frame
	currentTime := glfw.GetTime()
	deltaTime := float32(currentTime - lastTime)

	//// Get mouse position
	xpos, ypos := window.GetCursorPos()

	// Reset mouse position for next frame
	windowWidth, windowHeight := window.GetSize()

	centerX := float32(windowWidth / 2)
	centerY := float32(windowHeight / 2)

	window.SetCursorPos(float64(centerX), float64(centerY))

	// At first startup xpos and ypos register at (0, 0), don't calculate angle until this is adjusted or screen will
	// not start where we want it
	if xpos == 0 && ypos == 0 {
		xpos = float64(centerX)
		ypos = float64(centerY)
	}

	// Compute new orientation
	horizontalAngle += mouseSpeed * float32(centerX-float32(xpos))
	verticalAngle += mouseSpeed * float32(centerY-float32(ypos))

	direction := mgl32.Vec3{
		float32(math.Cos(float64(verticalAngle)) * math.Sin(float64(horizontalAngle))),
		float32(math.Sin(float64(verticalAngle))),
		float32(math.Cos(float64(verticalAngle)) * math.Cos(float64(horizontalAngle))),
	}

	right := &mgl32.Vec3{
		float32(math.Sin(float64(horizontalAngle - 3.14/2.0))),
		0.0,
		float32(math.Cos(float64(horizontalAngle - 3.14/2.0))),
	}

	// Up vector : perpendicular to both direction and right
	up := right.Cross(direction)

	// Move forward
	if window.GetKey(glfw.KeyW) == glfw.Press {
		position = position.Add(direction.Mul(deltaTime * speed))
	}

	// Move backward
	if window.GetKey(glfw.KeyS) == glfw.Press {
		position = position.Sub(direction.Mul(deltaTime).Mul(speed))
	}

	// Strafe right
	if window.GetKey(glfw.KeyD) == glfw.Press {
		position = position.Add(right.Mul(deltaTime).Mul(speed))
	}

	// Strafe left
	if window.GetKey(glfw.KeyA) == glfw.Press {
		position = position.Sub(right.Mul(deltaTime).Mul(speed))
	}

	// Projection matrix : 45&deg; Field of View, 4:3 ratio, display range : 0.1 unit <-> 100 units
	projectionMatrix = mgl32.Perspective(fov, 4.0/3.0, 0.1, 100.0)

	// Camera matrix
	viewMatrix = mgl32.LookAtV(
		position,                // Camera is here
		position.Add(direction), // and looks here : at the same position, plus "direction"
		up) // Head is up (set to 0,-1,0 to look upside-down)

	// For the next frame, the "last time" will be "now"
	lastTime = currentTime

}
