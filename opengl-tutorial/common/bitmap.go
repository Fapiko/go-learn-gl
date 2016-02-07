package common

import (
	"errors"
	"os"

	"github.com/go-gl/gl/all-core/gl"
)

func LoadBmpCustom(filepath string) (int32, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return 0, errors.New("Image could not be opened")
	}
	defer file.Close()

	//bufferedReader := bufio.NewReader(file)
	header := make([]byte, 54)

	_, err = file.Read(header)
	if err != nil {
		return 0, errors.New("Not a correct bmp file")
	}

	if header[0] != 'B' || header[1] != 'M' {
		return 0, errors.New("Not a correct bmp file")
	}

	// Read ints from the byte array
	dataPos := intFromByteSlice(header[0x0A : 0x0A+4])
	imageSize := intFromByteSlice(header[0x22 : 0x22+4])
	width := intFromByteSlice(header[0x12 : 0x12+4])
	height := intFromByteSlice(header[0x16 : 0x16+4])

	// Some BMP files are misformatted, guess missing information
	if imageSize == 0 {
		imageSize = width * height * 3 // 3 : one byte for each Red, Green and Blue component
	}

	if dataPos == 0 {
		dataPos = 54 // The BMP header is done that way
	}

	// Create a buffer
	data := make([]byte, 786432)

	// Read the actual data from the file into the buffer
	_, err = file.Read(data)
	if err != nil {
		return 0, errors.New("bmp smaller than expected")
	}

	// Create one OpenGL texture
	var textureId uint32
	gl.GenTextures(1, &textureId)

	// "Bind" the newly created texture : all future texture functions will modify this texture
	gl.BindTexture(gl.TEXTURE_2D, textureId)

	// Give the image to OpenGL
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(width), int32(height), 0, gl.BGR, gl.UNSIGNED_BYTE, gl.Ptr(data))

	// Poor filter, or ...
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	// ... nice trilinear filtering
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	// When MAGnifying the image (no bigger mipmap available), use LINEAR filtering
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// When MINifying the image, use a LINEAR blend of two mipmaps, each filtered LINEARLY too
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	// Generate mipmaps, by the way.
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return int32(textureId), nil

}
