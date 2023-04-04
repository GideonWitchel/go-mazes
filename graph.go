package main

import "fmt"

type neighbor struct {
	n      *node
	weight int
}

type node struct {
	val       int
	index     int
	neighbors []*neighbor
}

type graph struct {
	nodes []*node
}

func makeNode(v int, i int) *node {
	var newNode node
	newNode.val = v
	newNode.index = i
	return &newNode
}

func makeNeighbor(n *node, weight int) *neighbor {
	var e neighbor
	e.n = n
	e.weight = weight
	return &e
}

// AddNode places a node at the end of the graph's slice of nodes and returns the index of the node in the slice.
func (g *graph) AddNode(v int) int {
	g.nodes = append(g.nodes, makeNode(v, len(g.nodes)))
	return len(g.nodes) - 1
}

func (g *graph) SetNode(index int, val int) {
	g.nodes[index].val = val
}

// addEdge is undirected and assumes all weights are 1
func addEdge(n1 *node, n2 *node) {
	n1.neighbors = append(n1.neighbors, makeNeighbor(n2, 1))
	n2.neighbors = append(n2.neighbors, makeNeighbor(n1, 1))
}
func (g *graph) AddEdge(i1 int, i2 int) {
	// no duplicate edges
	// assumes edges are not directed
	for _, adj := range g.nodes[i1].neighbors {
		if adj.n.index == i2 {
			return
		}
	}
	addEdge(g.nodes[i1], g.nodes[i2])
}

// removeEdge is unidirectional
func removeEdge(g *graph, i1 int, i2 int) {
	for i, adj := range g.nodes[i1].neighbors {
		if adj.n.index == i2 {
			g.nodes[i1].neighbors = append(g.nodes[i1].neighbors[:i], g.nodes[i1].neighbors[i+1:]...)
			break
		}
	}
}

// RemoveEdge bidirectional
func (g *graph) RemoveEdge(i1 int, i2 int) {
	removeEdge(g, i1, i2)
	removeEdge(g, i2, i1)
}

// HasEdge assumes bidirectional
func (g *graph) HasEdge(i1 int, i2 int) bool {
	for _, adj := range g.nodes[i1].neighbors {
		if adj.n.index == i2 {
			return true
		}
	}
	return false
}

func getIndex(g *graph, n *node) int {
	for i, node := range g.nodes {
		if node == n {
			return i
		}
	}
	return -1
}

func (g *graph) Print() {
	// TODO cut off excess ", "
	for nodeIndex, currentNode := range g.nodes {
		// index and value of the node
		fmt.Printf("%v (%v) | [", nodeIndex, currentNode.val)
		for _, adj := range currentNode.neighbors {
			// index and value of each neighbor
			fmt.Printf("%v, ", getIndex(g, adj.n))
		}
		fmt.Print("]\n")
	}
}

// DFS finds the first node with a given value and returns:
// - a boolean which is true if the value is accessible
// - a slice of indexes with the order of nodes to get there, starting with the node of the desired value and ending with the starting node
// DFS is not parallelized
func dfs(g *graph, val int, startIndex int) (exists bool, path *[]int) {
	pathOut := make([]int, 0)
	visited := make([]bool, len(g.nodes), len(g.nodes))

	return dfsRecursive(g.nodes[startIndex], val, &visited, &pathOut), &pathOut
}

// dfsRecursive returns true if the value is found and false if the value is not
func dfsRecursive(n *node, val int, visited *[]bool, pathOut *[]int) bool {
	(*visited)[n.index] = true
	if n.val == val {
		*pathOut = append(*pathOut, n.index)
		return true
	}

	for _, currentNeighbor := range n.neighbors {
		if !(*visited)[currentNeighbor.n.index] {
			if dfsRecursive(currentNeighbor.n, val, visited, pathOut) {
				*pathOut = append(*pathOut, n.index)
				return true
			}
		}
	}

	return false
}
