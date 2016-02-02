// Go wrapper for iGraph
package gograph

/*
#cgo CFLAGS: -I/usr/local/include/igraph
#cgo LDFLAGS: -L/usr/local/lib -ligraph
#include <igraph.h>
#include <stdio.h>
*/
import "C"

import (
// "fmt"
)

// Go type wrapping an igraph graph
type Graph struct {
	graph    *C.igraph_t
	vertices []Vertex // map string vertex IDs to integers
	edges    *Vector
	weights  *Vector
}

func NewGraph() *Graph {
	var cgraph C.igraph_t
	var vertices []Vertex
	return &Graph{
		graph:    &cgraph,
		vertices: vertices,
		edges:    NewVector(),
		weights:  NewVector(),
	}
}

// release memory allocated to the graph
func (g *Graph) Cleanup() {
	C.igraph_destroy(g.graph)
	g.edges.Cleanup()
	g.weights.Cleanup()
}

type Vertex struct {
	Id     string
	Int_id int
}

type Edge struct {
	FromId string
	ToId   string
	Weight float64
}

// Go type wrapping an igraph vector (*not* the same as a C array)
type Vector struct {
	vec         *C.igraph_vector_t
	initialized bool
}

// Vector constructor
func NewVector() *Vector {
	var vec C.igraph_vector_t
	return &Vector{
		vec:         &vec,
		initialized: false,
	}
}

// Initialize the vector to the given length
func (v *Vector) Initialize(length int) {
	if v.initialized {
		v.Cleanup()
	}
	C.igraph_vector_init(v.vec, C.long(length))
	v.initialized = true
}

// Set the Vector's value at index i to value (int)
// (unclear what happens if index is out of range)
func (v *Vector) SetInt(value int, index int) {
	C.igraph_vector_set(v.vec, C.long(index), C.igraph_real_t(value))
}

// Set the Vector's value at index i to value (float)
// (unclear what happens if index is out of range)
func (v *Vector) SetFloat(value float64, index int) {
	C.igraph_vector_set(v.vec, C.long(index), C.igraph_real_t(value))
}

// Clean up memory allocated to Vector
func (v *Vector) Cleanup() {
	if v.initialized {
		C.igraph_vector_destroy(v.vec)
		v.initialized = false
	}
}

func (g *Graph) PopulateFromEdges(edges *[]Edge) {
	g.edges.Initialize(2 * len(*edges))
	g.weights.Initialize(len(*edges))
	vertices := make(map[string]int)
	var from_id, to_id int
	i := 0
	for c, el := range *edges {
		if id, ok := vertices[el.FromId]; ok {
			from_id = id
		} else {
			vertices[el.FromId] = i
			g.vertices = append(g.vertices, Vertex{el.FromId, i})
			i += 1
		}
		if id, ok := vertices[el.ToId]; ok {
			to_id = id
		} else {
			vertices[el.ToId] = i
			g.vertices = append(g.vertices, Vertex{el.ToId, i})
			i += 1
		}
		g.edges.SetInt(from_id, c*2)
		g.edges.SetInt(to_id, c*2+1)
		g.weights.SetFloat(el.Weight, i)
	}
	C.igraph_create(g.graph, g.edges.vec, 0, C.IGRAPH_DIRECTED)
}

func PageRank(edges *[]Edge, damping float64, multiplier float64) *map[string]float64 {
	g := NewGraph()
	ranks := make(map[string]float64)
	defer g.Cleanup()
	g.PopulateFromEdges(edges)
	result := NewVector()
	result.Initialize(0)
	defer result.Cleanup()
	realOne := C.igraph_real_t(1)
	C.igraph_pagerank(g.graph, C.IGRAPH_PAGERANK_ALGO_PRPACK, result.vec, &realOne,
		C.igraph_vss_all(), C.IGRAPH_DIRECTED, C.igraph_real_t(damping), g.weights.vec, nil)

	for i := 0; i < len(g.vertices); i++ {
		ranks[g.vertices[i].Id] = float64(C.igraph_vector_e(result.vec, C.long(i))) * multiplier
	}
	return &ranks
}
