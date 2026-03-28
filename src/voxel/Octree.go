package voxel

import "github.com/SlackingSlothh/OctreeVoxelization/src/geometry"

type OctreeNode struct {
	Boundary geometry.Cube
	Triangles []geometry.Triangle3D
	Children []OctreeNode
}

const MAX_DEPTH = 6
const MIN_EDGE = 0.005

type ConstructStats struct {
	NodesPerLevel  []int
	LeavesPerLevel []int
}

func (node *OctreeNode) Construct(maxDepth int, minEdge float64) (int, []int, []int) {
	stats := &ConstructStats{}
	actualMaxDepth := maxDepth
	if actualMaxDepth < 0 {
		actualMaxDepth = 99999
	}
	actualMinEdge := minEdge
	if actualMinEdge <= 0 {
		actualMinEdge = 0.0
	}
	maxDepthReached := node.constructHelper(0, stats, actualMaxDepth, actualMinEdge)
	return maxDepthReached, stats.NodesPerLevel, stats.LeavesPerLevel
}

func (node *OctreeNode) constructHelper(currentDepth int, stats *ConstructStats, maxDepth int, minEdge float64) int {
	for len(stats.NodesPerLevel) <= currentDepth {
		stats.NodesPerLevel = append(stats.NodesPerLevel, 0)
	}
	for len(stats.LeavesPerLevel) <= currentDepth {
		stats.LeavesPerLevel = append(stats.LeavesPerLevel, 0)
	}

	stats.NodesPerLevel[currentDepth]++

	if len(node.Triangles) == 0 {
		node.Children = nil
		stats.LeavesPerLevel[currentDepth]++
		return currentDepth
	}

	if (maxDepth >= 0 && currentDepth >= maxDepth) || node.Boundary.GetEdgeLength() <= minEdge {
		node.Children = nil
		stats.LeavesPerLevel[currentDepth]++
		return currentDepth
	}

	node.Children = nil
	dividedCube := node.Boundary.DivideCube()
	maxDepthReached := currentDepth

	for _, childCube := range dividedCube {
		childNode := OctreeNode{Boundary: childCube}
		for _, triangle := range node.Triangles {
			if triangle.IsIntersecting(childCube) {
				childNode.Triangles = append(childNode.Triangles, triangle)
			}
		}

		if len(childNode.Triangles) == 0 {
			childNode.Children = nil
			for len(stats.NodesPerLevel) <= currentDepth+1 {
				stats.NodesPerLevel = append(stats.NodesPerLevel, 0)
			}
			for len(stats.LeavesPerLevel) <= currentDepth+1 {
				stats.LeavesPerLevel = append(stats.LeavesPerLevel, 0)
			}
			stats.NodesPerLevel[currentDepth+1]++
			stats.LeavesPerLevel[currentDepth+1]++
			node.Children = append(node.Children, childNode)
			continue
		}

		childDepth := childNode.constructHelper(currentDepth+1, stats, maxDepth, minEdge)
		if childDepth > maxDepthReached {
			maxDepthReached = childDepth
		}
		node.Children = append(node.Children, childNode)
	}

	if len(node.Children) == 0 {
		stats.LeavesPerLevel[currentDepth]++
	}

	return maxDepthReached
}
