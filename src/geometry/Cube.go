package geometry

import "github.com/golang/geo/r3"

type Cube struct {
	v0, v1 r3.Vector
}

func (cube Cube) GetEdgeLength() float64 {
	return cube.v1.X - cube.v0.X
}

func (parent Cube) DivideCube() []Cube {
	mid := parent.v0.Add(parent.v1).Mul(0.5)

	return []Cube{
		// 0,0,0
		{parent.v0, mid},
		// 0,0,1
		{r3.Vector{X: parent.v0.X, Y: parent.v0.Y, Z: mid.Z}, r3.Vector{X: mid.X, Y: mid.Y, Z: parent.v1.Z}},
		// 0,1,0
		{r3.Vector{X: parent.v0.X, Y: mid.Y, Z: parent.v0.Z}, r3.Vector{X: mid.X, Y: parent.v1.Y, Z: mid.Z}},
		// 0,1,1
		{r3.Vector{X: parent.v0.X, Y: mid.Y, Z: mid.Z}, r3.Vector{X: mid.X, Y: parent.v1.Y, Z: parent.v1.Z}},
		// 1,0,0
		{r3.Vector{X: mid.X, Y: parent.v0.Y, Z: parent.v0.Z}, r3.Vector{X: parent.v1.X, Y: mid.Y, Z: mid.Z}},
		// 1,0,1
		{r3.Vector{X: mid.X, Y: parent.v0.Y, Z: mid.Z}, r3.Vector{X: parent.v1.X, Y: mid.Y, Z: parent.v1.Z}},
		// 1,1,0
		{r3.Vector{X: mid.X, Y: mid.Y, Z: parent.v0.Z}, r3.Vector{X: parent.v1.X, Y: parent.v1.Y, Z: mid.Z}},
		// 1,1,1
		{r3.Vector{X: mid.X, Y: mid.Y, Z: mid.Z}, parent.v1},
	}
}