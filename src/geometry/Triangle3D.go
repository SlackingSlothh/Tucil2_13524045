package geometry

import (
	"math"

	"github.com/golang/geo/r3"
)

type Triangle3D struct {
	v1, v2, v3 r3.Vector
}

func (tri Triangle3D) projectTriangle(axis r3.Vector) (float64, float64) {
	min := axis.Dot(tri.v1)
	max := min
	projection := axis.Dot(tri.v2)
	if projection < min {
		min = projection
	} else if projection > max {
		max = projection
	}
	projection = axis.Dot(tri.v3)
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

func (tri Triangle3D) isIntersecting(cube Cube) bool {
	center := r3.Vector{X: (cube.v0.X + cube.v1.X) * 0.5, Y: (cube.v0.Y + cube.v1.Y) * 0.5, Z: (cube.v0.Z + cube.v1.Z) * 0.5}
	extents := r3.Vector{X: (cube.v1.X-cube.v0.X)*0.5, Y: (cube.v1.Y-cube.v0.Y)*0.5, Z: (cube.v1.Z-cube.v0.Z)*0.5}.Abs()

	v0 := r3.Vector{X: tri.v1.X - center.X, Y: tri.v1.Y - center.Y, Z: tri.v1.Z - center.Z}
	v1 := r3.Vector{X: tri.v2.X - center.X, Y: tri.v2.Y - center.Y, Z: tri.v2.Z - center.Z}
	v2 := r3.Vector{X: tri.v3.X - center.X, Y: tri.v3.Y - center.Y, Z: tri.v3.Z - center.Z}
	translatedTri := Triangle3D{v0, v1, v2}

	minX, maxX := v0.X, v0.X
	minY, maxY := v0.Y, v0.Y
	minZ, maxZ := v0.Z, v0.Z
	for _, v := range []r3.Vector{v1, v2} {
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

	e0 := r3.Vector{X: v1.X - v0.X, Y: v1.Y - v0.Y, Z: v1.Z - v0.Z}
	e1 := r3.Vector{X: v2.X - v1.X, Y: v2.Y - v1.Y, Z: v2.Z - v1.Z}
	e2 := r3.Vector{X: v0.X - v2.X, Y: v0.Y - v2.Y, Z: v0.Z - v2.Z}
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
