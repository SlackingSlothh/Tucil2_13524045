package geometry

import (
	"math"

	"github.com/golang/geo/r3"
)

type Triangle3D struct {
	V1, V2, V3 r3.Vector
}

func (tri Triangle3D) projectTriangle(axis r3.Vector) (float64, float64) {
	min := axis.Dot(tri.V1)
	max := min
	projection := axis.Dot(tri.V2)
	if projection < min {
		min = projection
	} else if projection > max {
		max = projection
	}
	projection = axis.Dot(tri.V3)
	if projection < min {
		min = projection
	} else if projection > max {
		max = projection
	}
	return min, max
}

func (tri Triangle3D) isOverlapAxis(axis r3.Vector, extents r3.Vector) bool {
	epsilon := math.Nextafter(1.0, 2.0) - 1.0
	if axis.Norm() < epsilon {
		return true
	}
	min, max := tri.projectTriangle(axis)
	radius := extents.X*math.Abs(axis.X) + extents.Y*math.Abs(axis.Y) + extents.Z*math.Abs(axis.Z)
	return !(min > radius || max < -radius)
}

func (tri Triangle3D) IsIntersecting(cube Cube) bool {
	center := r3.Vector{X: (cube.V0.X + cube.V1.X) * 0.5, Y: (cube.V0.Y + cube.V1.Y) * 0.5, Z: (cube.V0.Z + cube.V1.Z) * 0.5}
	extents := r3.Vector{X: (cube.V1.X-cube.V0.X)*0.5, Y: (cube.V1.Y-cube.V0.Y)*0.5, Z: (cube.V1.Z-cube.V0.Z)*0.5}.Abs()

	V0 := r3.Vector{X: tri.V1.X - center.X, Y: tri.V1.Y - center.Y, Z: tri.V1.Z - center.Z}
	V1 := r3.Vector{X: tri.V2.X - center.X, Y: tri.V2.Y - center.Y, Z: tri.V2.Z - center.Z}
	V2 := r3.Vector{X: tri.V3.X - center.X, Y: tri.V3.Y - center.Y, Z: tri.V3.Z - center.Z}
	translatedTri := Triangle3D{V0, V1, V2}

	minX, maxX := V0.X, V0.X
	minY, maxY := V0.Y, V0.Y
	minZ, maxZ := V0.Z, V0.Z
	for _, v := range []r3.Vector{V1, V2} {
		if v.X < minX { minX = v.X }
		if v.X > maxX { maxX = v.X }
		if v.Y < minY { minY = v.Y }
		if v.Y > maxY { maxY = v.Y }
		if v.Z < minZ { minZ = v.Z }
		if v.Z > maxZ { maxZ = v.Z }
	}
	if maxX < -extents.X || minX > extents.X {
		return false
	}
	if maxY < -extents.Y || minY > extents.Y {
		return false
	}
	if maxZ < -extents.Z || minZ > extents.Z {
		return false
	}

	e0 := r3.Vector{X: V1.X - V0.X, Y: V1.Y - V0.Y, Z: V1.Z - V0.Z}
	e1 := r3.Vector{X: V2.X - V1.X, Y: V2.Y - V1.Y, Z: V2.Z - V1.Z}
	e2 := r3.Vector{X: V0.X - V2.X, Y: V0.Y - V2.Y, Z: V0.Z - V2.Z}
	axes := []r3.Vector{
		{X: 0, Y: -e0.Z, Z: e0.Y}, {X: e0.Z, Y: 0, Z: -e0.X}, {X: -e0.Y, Y: e0.X, Z: 0},
		{X: 0, Y: -e1.Z, Z: e1.Y}, {X: e1.Z, Y: 0, Z: -e1.X}, {X: -e1.Y, Y: e1.X, Z: 0},
		{X: 0, Y: -e2.Z, Z: e2.Y}, {X: e2.Z, Y: 0, Z: -e2.X}, {X: -e2.Y, Y: e2.X, Z: 0},
	}
	for _, axis := range axes {
		if !translatedTri.isOverlapAxis(axis, extents) {
			return false
		}
	}

	normal := e0.Cross(e1)
	if !translatedTri.isOverlapAxis(normal, extents) {
		return false
	}

	return true
}
