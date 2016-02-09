package common

import "github.com/go-gl/mathgl/mgl32"

type PackedVertex struct {
	position mgl32.Vec3
	uv       mgl32.Vec2
	normal   mgl32.Vec3
}

func IndexVBO(vertices []mgl32.Vec3, uvs []mgl32.Vec2, normals []mgl32.Vec3) ([]uint32, []mgl32.Vec3, []mgl32.Vec2, []mgl32.Vec3) {

	var indices []uint32
	var indexedVertices []mgl32.Vec3
	var indexedUvs []mgl32.Vec2
	var indexedNormals []mgl32.Vec3

	vertexToOutIndex := make(map[PackedVertex]uint32, 0)

	for i := 0; i < len(vertices); i++ {

		packed := &PackedVertex{
			position: vertices[i],
			uv:       uvs[i],
			normal:   normals[i],
		}

		if index, ok := vertexToOutIndex[*packed]; ok {

			indices = append(indices, index)

		} else {

			indexedVertices = append(indexedVertices, vertices[i])
			indexedUvs = append(indexedUvs, uvs[i])
			indexedNormals = append(indexedNormals, normals[i])

			var newIndex uint32
			if len(indices) > 0 {
				newIndex = uint32(len(indexedVertices) - 1)
			} else {
				newIndex = 0
			}
			indices = append(indices, newIndex)
			vertexToOutIndex[*packed] = newIndex

		}

	}

	return indices, indexedVertices, indexedUvs, indexedNormals

}
