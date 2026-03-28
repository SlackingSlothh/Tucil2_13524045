package geometry

import "github.com/golang/geo/r3"

type Cube struct {
	V0, V1 r3.Vector
}

func (cube Cube) GetEdgeLength() float64 {
	return cube.V1.X - cube.V0.X
}

func (parent Cube) DivideCube() []Cube {
	mid := parent.V0.Add(parent.V1).Mul(0.5)

	return []Cube{
		// 0,0,0
		{parent.V0, mid},
		// 0,0,1
		{r3.Vector{X: parent.V0.X, Y: parent.V0.Y, Z: mid.Z}, r3.Vector{X: mid.X, Y: mid.Y, Z: parent.V1.Z}},
		// 0,1,0
		{r3.Vector{X: parent.V0.X, Y: mid.Y, Z: parent.V0.Z}, r3.Vector{X: mid.X, Y: parent.V1.Y, Z: mid.Z}},
		// 0,1,1
		{r3.Vector{X: parent.V0.X, Y: mid.Y, Z: mid.Z}, r3.Vector{X: mid.X, Y: parent.V1.Y, Z: parent.V1.Z}},
		// 1,0,0
		{r3.Vector{X: mid.X, Y: parent.V0.Y, Z: parent.V0.Z}, r3.Vector{X: parent.V1.X, Y: mid.Y, Z: mid.Z}},
		// 1,0,1
		{r3.Vector{X: mid.X, Y: parent.V0.Y, Z: mid.Z}, r3.Vector{X: parent.V1.X, Y: mid.Y, Z: parent.V1.Z}},
		// 1,1,0
		{r3.Vector{X: mid.X, Y: mid.Y, Z: parent.V0.Z}, r3.Vector{X: parent.V1.X, Y: parent.V1.Y, Z: mid.Z}},
		// 1,1,1
		{r3.Vector{X: mid.X, Y: mid.Y, Z: mid.Z}, parent.V1},
	}
}