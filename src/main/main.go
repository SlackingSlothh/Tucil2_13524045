package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/SlackingSlothh/OctreeVoxelization/src/inputoutput"
	"github.com/SlackingSlothh/OctreeVoxelization/src/voxel"
)

func readLine(prompt string, reader *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n--- Octree Voxelizer Menu ---")
		fmt.Println("1. Start")
		fmt.Println("2. Quit")
		choice, err := readLine("Choose [1-2]: ", reader)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}

		if choice == "2" || strings.EqualFold(choice, "quit") {
			fmt.Println("Exiting")
			return
		}
		if choice != "1" && !strings.EqualFold(choice, "start") {
			fmt.Println("Invalid option")
			continue
		}

		cwd, _ := os.Getwd()
		fmt.Println("Current directory:", cwd)

		inputPath, err := readLine("Enter OBJ input path (relative): ", reader)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading path:", err)
			continue
		}

		inputPath = filepath.Clean(inputPath)
		fullInputPath := inputPath
		if !filepath.IsAbs(inputPath) {
			fullInputPath = filepath.Join(cwd, inputPath)
		}
		if _, err := os.Stat(fullInputPath); err != nil {
			fmt.Println("Input file not found. Returning to menu.")
			continue
		}

		maxDepthStr, err := readLine("Enter max octree depth (negative for no max, default 6): ", reader)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading max depth:", err)
			continue
		}
		maxDepth64, err := strconv.ParseInt(maxDepthStr, 10, 64)
		if err != nil {
			fmt.Println("Invalid max depth. Defaulting to 6.")
			maxDepth64 = 6
		}
		maxDepth := int(maxDepth64)

		minVoxelStr, err := readLine("Enter minimum voxel size (positive float, default 0.01): ", reader)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading min voxel size:", err)
			continue
		}
		minVoxel, err := strconv.ParseFloat(minVoxelStr, 64)
		if err != nil || minVoxel <= 0 {
			fmt.Println("Invalid minimum voxel size. Defaulting to 0.01.")
			minVoxel = 0.01
		}

		start := time.Now()
		cube, triangles := inputoutput.ReadObj(fullInputPath)
		rootOctree := voxel.OctreeNode{Boundary: cube, Triangles: triangles, Children: nil}
		maxDepthReached, nodesPerLevel, leavesPerLevel := rootOctree.Construct(maxDepth, minVoxel)

		var voxelSize float64
		if maxDepth >= 0 {
			voxelSize = cube.GetEdgeLength() * math.Pow(0.5, float64(maxDepthReached))
		} else {
			voxelSize = minVoxel
		}
		voxelizer := voxel.Voxelizer{RootOctree: rootOctree, VoxelSize: voxelSize, VoxelMap: make(map[voxel.Index3D]bool)}
		voxelizer.OctreeToHashMap(rootOctree)
		faces, vertices := voxelizer.MakeVerticesFaces()
		elapsed := time.Since(start)

		fmt.Printf("Voxel count: %d\n", len(voxelizer.VoxelMap))
		fmt.Printf("Vertex count: %d\n", len(vertices))
		fmt.Printf("Face count: %d\n", len(faces))
		fmt.Println("Nodes per level")
		for i, nodeCount := range nodesPerLevel {
			fmt.Printf("%d: %d\n", i + 1, nodeCount)
		}
		fmt.Println("Empty nodes per level")
		for i, nodeCount := range leavesPerLevel {
			fmt.Printf("%d: %d\n", i + 1, nodeCount)
		}
		fmt.Printf("Max depth reached: %d\n", maxDepthReached)
		fmt.Printf("Processing time: %v\n", elapsed)

		outputPath, err := readLine("Enter output OBJ path (relative): ", reader)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading output path:", err)
			continue
		}
		outputPath = filepath.Clean(outputPath)
		if !filepath.IsAbs(outputPath) {
			outputPath = filepath.Join(cwd, outputPath)
		}
		inputoutput.WriteObj(outputPath, faces, vertices)
		fmt.Println("Results written to", outputPath)
	}
}