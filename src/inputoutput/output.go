package inputoutput

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/SlackingSlothh/OctreeVoxelization/src/voxel"
)

func WriteObj(path string, faces []voxel.Face, vertices []voxel.VertexIndexPair) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, pair := range vertices {
		_, err = fmt.Fprintf(writer, "v %f %f %f\n", pair.Vertex.X, pair.Vertex.Y, pair.Vertex.Z)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, face := range faces {
		_, err = fmt.Fprintf(writer, "f %d %d %d\n", face.V1, face.V2, face.V3)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := writer.Flush(); err != nil {
		log.Fatal(err)
	}
}