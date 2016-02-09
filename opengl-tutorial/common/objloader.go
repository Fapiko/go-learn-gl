package common

import (
	"errors"
	"os"

	"bufio"

	"bytes"

	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

func LoadObj(path string) ([]mgl32.Vec3, []mgl32.Vec2, []mgl32.Vec3, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, errors.New("Impossible to open the file!")
	}

	// Used in this method for parsing the normalized obj file
	tmpVertices := make([]mgl32.Vec3, 0)
	tmpUvs := make([]mgl32.Vec2, 0)
	tmpNormals := make([]mgl32.Vec3, 0)

	// Denormalized vectors to be returned
	vertices := make([]mgl32.Vec3, 0)
	uvs := make([]mgl32.Vec2, 0)
	normals := make([]mgl32.Vec3, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		buffer := bytes.NewBuffer(scanner.Bytes())

		lineHeader, err := buffer.ReadString(' ')
		if err != nil {
			// EOF, continue
			continue
		}

		// Parse vertex
		if lineHeader == "v " {

			var vertex mgl32.Vec3
			fmt.Fscanf(buffer, "%g %g %g\n", &vertex[0], &vertex[1], &vertex[2])
			tmpVertices = append(tmpVertices, vertex)

		} else if lineHeader == "vt " { // Parse UV

			var uv mgl32.Vec2
			fmt.Fscanf(buffer, "%g %g\n", &uv[0], &uv[1])
			tmpUvs = append(tmpUvs, uv)

		} else if lineHeader == "vn " { // Parse vector normal

			var normal mgl32.Vec3
			fmt.Fscanf(buffer, "%g %g %g\n", &normal[0], &normal[1], &normal[2])
			tmpNormals = append(tmpNormals, normal)

		} else if lineHeader == "f " { // Parse faces

			vertexIndices := make([]int, 3)
			uvIndices := make([]int, 3)
			normalIndices := make([]int, 3)

			matches, err := fmt.Fscanf(buffer, "%d/%d/%d %d/%d/%d %d/%d/%d\n",
				&vertexIndices[0], &uvIndices[0], &normalIndices[0],
				&vertexIndices[1], &uvIndices[1], &normalIndices[1],
				&vertexIndices[2], &uvIndices[2], &normalIndices[2])

			if matches != 9 || err != nil {
				return nil, nil, nil, errors.New("File can't be read by our simple parser : ( Try exporting with other options\n")
			}

			vertices = append(vertices,
				tmpVertices[vertexIndices[0]-1],
				tmpVertices[vertexIndices[1]-1],
				tmpVertices[vertexIndices[2]-1])

			uvs = append(uvs,
				tmpUvs[uvIndices[0]-1],
				tmpUvs[uvIndices[1]-1],
				tmpUvs[uvIndices[2]-1])

			normals = append(normals,
				tmpNormals[normalIndices[0]-1],
				tmpNormals[normalIndices[1]-1],
				tmpNormals[normalIndices[2]-1])

		}

	}

	return vertices, uvs, normals, nil

}
