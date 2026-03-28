package voxel

import (
	"sync"

	"github.com/SlackingSlothh/OctreeVoxelization/src/geometry"
)

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
	}
	var wg sync.WaitGroup
	results := make(chan struct {
		node      OctreeNode
		depth     int
		nodes     []int
		leaves    []int
	}, len(dividedCube))

	for _, childCube := range dividedCube {
		wg.Add(1)
		go func(childCube geometry.Cube) {
			defer wg.Done()
			childNode := OctreeNode{Boundary: childCube}
			for _, triangle := range node.Triangles {
				if triangle.IsIntersecting(childCube) {
					childNode.Triangles = append(childNode.Triangles, triangle)
				}
			}

			if len(childNode.Triangles) == 0 {
				childNode.Children = nil
				localNodes := []int{}
				localLeaves := []int{}
				if len(stats.NodesPerLevel) <= currentDepth+1 {
					localNodes = make([]int, currentDepth+2)
				} else {
					localNodes = make([]int, currentDepth+2)
				}
				if len(stats.LeavesPerLevel) <= currentDepth+1 {
					localLeaves = make([]int, currentDepth+2)
				} else {
					localLeaves = make([]int, currentDepth+2)
				}
				localNodes[currentDepth+1] = 1
				localLeaves[currentDepth+1] = 1
				results <- struct {
					node   OctreeNode
					depth  int
					nodes  []int
					leaves []int
				}{childNode, currentDepth + 1, localNodes, localLeaves}
				return
			}

			childNodes := &ConstructStats{}
			childDepth := childNode.constructHelper(currentDepth+1, childNodes, maxDepth, minEdge)
			results <- struct {
				node   OctreeNode
				depth  int
				nodes  []int
				leaves []int
			}{childNode, childDepth, childNodes.NodesPerLevel, childNodes.LeavesPerLevel}
		}(childCube)
	}

	wg.Wait()
	close(results)

	for res := range results {
		if len(stats.NodesPerLevel) < len(res.nodes) {
			newNodes := make([]int, len(res.nodes))
			copy(newNodes, stats.NodesPerLevel)
			stats.NodesPerLevel = newNodes
		}
		if len(stats.LeavesPerLevel) < len(res.leaves) {
			newLeaves := make([]int, len(res.leaves))
			copy(newLeaves, stats.LeavesPerLevel)
			stats.LeavesPerLevel = newLeaves
		}
		for i, v := range res.nodes {
			stats.NodesPerLevel[i] += v
		}
		for i, v := range res.leaves {
			stats.LeavesPerLevel[i] += v
		}

		node.Children = append(node.Children, res.node)
		if res.depth > maxDepthReached {
			maxDepthReached = res.depth
		}
	}

	if len(node.Children) == 0 {
		stats.LeavesPerLevel[currentDepth]++
	}

	return maxDepthReached
}
