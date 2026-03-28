package voxel

import (
	"math"
	"sort"

	"github.com/golang/geo/r3"
)

type Index3D struct {
	x, y, z int
}

type Voxelizer struct {
	RootOctree OctreeNode
	VoxelSize float64
	VoxelMap map[Index3D]bool
}

func (voxelizer Voxelizer) OctreeToHashMap(octree OctreeNode) {
	if octree.Children == nil {
		if len(octree.Triangles) > 0 {
			x := int(math.Round((octree.Boundary.V0.X - voxelizer.RootOctree.Boundary.V0.X) / voxelizer.VoxelSize))
			y := int(math.Round((octree.Boundary.V0.Y - voxelizer.RootOctree.Boundary.V0.Y) / voxelizer.VoxelSize))
			z := int(math.Round((octree.Boundary.V0.Z - voxelizer.RootOctree.Boundary.V0.Z) / voxelizer.VoxelSize))
			voxelizer.VoxelMap[Index3D{x, y, z}] = true
		}
		return
	}
	for _, child := range octree.Children {
		voxelizer.OctreeToHashMap(child)
	}
}

type Face struct {
	V1, V2, V3 int
}

type VertexIndexPair struct {
    Vertex r3.Vector
    Index int
}

func vertexGridIndexToWorld(base r3.Vector, voxelSize float64, grid Index3D) r3.Vector {
	return r3.Vector{
		X: base.X + float64(grid.x)*voxelSize,
		Y: base.Y + float64(grid.y)*voxelSize,
		Z: base.Z + float64(grid.z)*voxelSize,
	}
}

func getVertexIndex(vertexMap map[Index3D]int, grid Index3D, vertexNum *int) int {
	if idx, ok := vertexMap[grid]; ok {
		return idx
	}
	*vertexNum++
	vertexMap[grid] = *vertexNum
	return *vertexNum
}

func (voxelizer Voxelizer) MakeVerticesFaces() ([]Face, []VertexIndexPair) {
	vertexNum := 0
	vertexMap := make(map[Index3D]int)
	faceList := make([]Face, 0)

	for voxelIndex := range voxelizer.VoxelMap {
		sides := []struct {
			neighbor Index3D
			v0, v1, v2, v3 Index3D
		}{
			{Index3D{voxelIndex.x + 1, voxelIndex.y, voxelIndex.z}, Index3D{voxelIndex.x + 1, voxelIndex.y, voxelIndex.z}, Index3D{voxelIndex.x + 1, voxelIndex.y, voxelIndex.z + 1}, Index3D{voxelIndex.x + 1, voxelIndex.y + 1, voxelIndex.z}, Index3D{voxelIndex.x + 1, voxelIndex.y + 1, voxelIndex.z + 1}},
			{Index3D{voxelIndex.x - 1, voxelIndex.y, voxelIndex.z}, Index3D{voxelIndex.x, voxelIndex.y, voxelIndex.z}, Index3D{voxelIndex.x, voxelIndex.y, voxelIndex.z + 1}, Index3D{voxelIndex.x, voxelIndex.y + 1, voxelIndex.z}, Index3D{voxelIndex.x, voxelIndex.y + 1, voxelIndex.z + 1}},
			{Index3D{voxelIndex.x, voxelIndex.y + 1, voxelIndex.z}, Index3D{voxelIndex.x, voxelIndex.y + 1, voxelIndex.z}, Index3D{voxelIndex.x, voxelIndex.y + 1, voxelIndex.z + 1}, Index3D{voxelIndex.x + 1, voxelIndex.y + 1, voxelIndex.z}, Index3D{voxelIndex.x + 1, voxelIndex.y + 1, voxelIndex.z + 1}},
			{Index3D{voxelIndex.x, voxelIndex.y - 1, voxelIndex.z}, Index3D{voxelIndex.x, voxelIndex.y, voxelIndex.z}, Index3D{voxelIndex.x, voxelIndex.y, voxelIndex.z + 1}, Index3D{voxelIndex.x + 1, voxelIndex.y, voxelIndex.z}, Index3D{voxelIndex.x + 1, voxelIndex.y, voxelIndex.z + 1}},
			{Index3D{voxelIndex.x, voxelIndex.y, voxelIndex.z + 1}, Index3D{voxelIndex.x, voxelIndex.y, voxelIndex.z + 1}, Index3D{voxelIndex.x, voxelIndex.y + 1, voxelIndex.z + 1}, Index3D{voxelIndex.x + 1, voxelIndex.y, voxelIndex.z + 1}, Index3D{voxelIndex.x + 1, voxelIndex.y + 1, voxelIndex.z + 1}},
			{Index3D{voxelIndex.x, voxelIndex.y, voxelIndex.z - 1}, Index3D{voxelIndex.x, voxelIndex.y, voxelIndex.z}, Index3D{voxelIndex.x, voxelIndex.y + 1, voxelIndex.z}, Index3D{voxelIndex.x + 1, voxelIndex.y, voxelIndex.z}, Index3D{voxelIndex.x + 1, voxelIndex.y + 1, voxelIndex.z}},
		}

		for _, s := range sides {
			if _, ok := voxelizer.VoxelMap[s.neighbor]; ok {
				continue
			}

			v1 := getVertexIndex(vertexMap, s.v0, &vertexNum)
			v2 := getVertexIndex(vertexMap, s.v1, &vertexNum)
			v3 := getVertexIndex(vertexMap, s.v2, &vertexNum)
			v4 := getVertexIndex(vertexMap, s.v3, &vertexNum)

			faceList = append(faceList, Face{v1, v2, v3})
			faceList = append(faceList, Face{v4, v2, v3})
		}
	}

	vertices := make([]VertexIndexPair, 0, len(vertexMap))
	base := voxelizer.RootOctree.Boundary.V0
	for grid, idx := range vertexMap {
		v := vertexGridIndexToWorld(base, voxelizer.VoxelSize, grid)
		vertices = append(vertices, VertexIndexPair{Vertex: v, Index: idx})
	}

	sort.Slice(vertices, func(i, j int) bool {
		return vertices[i].Index < vertices[j].Index
	})

	return faceList, vertices
}