package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var text2dTextureId uint32
var text2dVertexBufferId uint32
var text2dUVBufferId uint32
var text2dShaderId uint32
var text2dUniformId int32
var context *glfw.Window

func InitText2d(texturePath string) {

	var err error

	// Initialize the texture
	text2dTextureId, err = LoadDDS(texturePath)
	if err != nil {
		log.Error(err)
		return
	}

	// Initialize VBO
	gl.GenBuffers(1, &text2dVertexBufferId)
	gl.GenBuffers(1, &text2dUVBufferId)

	// Initialize Shader
	text2dShaderId = LoadShaders("TextVertexShader.vertexshader", "TextVertexShader.fragmentshader")

	// Initialize uniforms' IDs
	text2dUniformId = gl.GetUniformLocation(text2dShaderId, gl.Str("myTextureSampler\x00"))

}

func PrintText2D(text string, x int, y int, size int) {

	length := len(text)

	// Fill buffers
	vertices := make([]mgl32.Vec2, 0)
	uvs := make([]mgl32.Vec2, 0)
	for i := 0; i < length; i++ {

		vertexUpLeft := mgl32.Vec2{float32(x + i*size), float32(y + size)}
		vertexUpRight := mgl32.Vec2{float32(x + i*size + size), float32(y + size)}
		vertexDownRight := mgl32.Vec2{float32(x + i*size + size), float32(y)}
		vertexDownLeft := mgl32.Vec2{float32(x + i*size), float32(y)}

		vertices = append(vertices, vertexUpLeft, vertexDownLeft, vertexUpRight, vertexDownRight, vertexUpRight, vertexDownLeft)

		character := int(text[i])

		uvX := float32(character%16) / 16.0
		uvY := float32(character/16) / 16.0

		log.Println(character, uvX, uvY)

		uvUpLeft := mgl32.Vec2{float32(uvX), float32(uvY)}
		uvUpRight := mgl32.Vec2{float32(uvX + float32(1.0/16.0)), float32(uvY)}
		uvDownRight := mgl32.Vec2{float32(uvX + float32(1.0/16.0)), float32(uvY + float32(1.0/16.0))}
		uvDownLeft := mgl32.Vec2{float32(uvX), float32(uvY + float32(1.0/16.0))}

		uvs = append(uvs, uvUpLeft, uvDownLeft, uvUpRight, uvDownRight, uvUpRight, uvDownLeft)

	}

	gl.BindBuffer(gl.ARRAY_BUFFER, text2dVertexBufferId)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4*2, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, text2dUVBufferId)
	gl.BufferData(gl.ARRAY_BUFFER, len(uvs)*4*2, gl.Ptr(uvs), gl.STATIC_DRAW)

	// Bind shader
	gl.UseProgram(text2dShaderId)

	// Bind texture
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, uint32(text2dTextureId))
	// Set our "myTextureSampler" sampler to user Texture Unit 0
	gl.Uniform1i(text2dUniformId, 0)

	// 1st attribute buffer : vertices
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, text2dVertexBufferId)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
	defer gl.DisableVertexAttribArray(0)

	// 2nd attribute buffer : UVs
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, text2dUVBufferId)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, nil)
	defer gl.DisableVertexAttribArray(1)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	defer gl.Disable(gl.BLEND)

	// Draw call
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)))

}

func CleanupText2d() {

	// Delete buffers
	gl.DeleteBuffers(1, &text2dVertexBufferId)
	gl.DeleteBuffers(1, &text2dUVBufferId)

	// Delete texture
	gl.DeleteTextures(1, &text2dTextureId)

	// Delete shader
	gl.DeleteProgram(text2dShaderId)

}
