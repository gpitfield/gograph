package gograph

import (
	"testing"
)

func TestGraph(t *testing.T) {
	edges := []Edge{
		Edge{"0", "1", 0.99},
		Edge{"0", "2", 0.01},
		Edge{"2", "1", 0.3},
		Edge{"2", "0", 0.1},
	}
	PageRank(&edges, 0.85)
}
