package common

import (
	"errors"
	"os"

	"fmt"

	"io/ioutil"

	"bytes"
	"encoding/binary"

	"github.com/go-gl/gl/v3.3-core/gl"
)

const (
	FOURCC_DXT1 = 0x31545844
	FOURCC_DXT3 = 0x33545844
	FOURCC_DXT5 = 0x35545844
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

func LoadDDS(imagepath string) (uint32, error) {

	// try to open the file
	file, err := os.Open(imagepath)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("%s could not be opened. Are you in the right directory ? Don't forget to "+
			"read the FAQ !\n", imagepath))
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, err
	}

	// verify the type of file
	if string(fileData[0:4]) != "DDS " {
		return 0, errors.New("File was not of type DDS")
	}

	// get the surface desc
	height := int32FromByteSlice(fileData[12:16])
	width := int32FromByteSlice(fileData[16:20])
	mipMapCount := int32FromByteSlice(fileData[28:32])
	fourCC := int32FromByteSlice(fileData[84:88])

	var format uint32
	switch fourCC {
	case FOURCC_DXT1:
		format = gl.COMPRESSED_RGBA_S3TC_DXT1_EXT
		break
	case FOURCC_DXT3:
		format = gl.COMPRESSED_RGBA_S3TC_DXT3_EXT
		break
	case FOURCC_DXT5:
		format = gl.COMPRESSED_RGBA_S3TC_DXT5_EXT
		break
	default:
		return 0, errors.New("FourCC not recognized")
	}

	// Create one OpenGL texture
	var textureId uint32
	gl.GenTextures(1, &textureId)

	// "Bind" the newly created texture : all future texture functions will modify this texture
	gl.BindTexture(gl.TEXTURE_2D, textureId)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	var blockSize int32
	if format == gl.COMPRESSED_RGBA_S3TC_DXT1_EXT {
		blockSize = 8
	} else {
		blockSize = 16
	}

	buffer := make([]byte, len(fileData)-128)

	err = binary.Read(bytes.NewReader(fileData[128:]), binary.LittleEndian, buffer)
	if err != nil {
		return 0, err
	}

	var offset int32
	for level := int32(0); level < mipMapCount && (width > 0 || height > 0); level++ {

		size := int32(((width + 3) / 4) * ((height + 3) / 4) * blockSize)

		gl.CompressedTexImage2D(gl.TEXTURE_2D, level, format, width, height, 0, size, gl.Ptr(&buffer[offset]))

		offset += size
		width /= 2
		height /= 2

		//// Deal with Non-Power-Of-Two textures. This code is not included in the webpage to reduce clutter.
		if width < 1 {
			width = 1
		}

		if height < 1 {
			height = 1
		}

	}

	return textureId, nil

}
