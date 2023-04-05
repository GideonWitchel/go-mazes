package main

import (
	"fmt"
	"sync"
)

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

type visited struct {
	// v is a visited array that also keeps track of the thread ID of whoever claimed the node.
	// int default is 0, so all thread IDs must be greater than 0, and 0 is unclaimed.
	v *[]int
	// found represents the path ID (thread ID - 1) of whoever found the solution node
	found int
	sync.Mutex
	sync.WaitGroup
}

// dfsMultithreaded finds a value in a graph using a number of simultaneous dfs searches with a shared visited list.
// exists is an index which specifies which search ended up finding the value in the paths array.
// If exists is -1, there is valid path to the solution from any starting index.
func dfsMultithreaded(g *graph, val int, startIndecies []int) (exists int, p *[][]int) {
	pathsOut := make([][]int, len(startIndecies), len(startIndecies))

	visitedArray := make([]int, len(g.nodes), len(g.nodes))
	visited := visited{
		v:     &visitedArray,
		found: -1,
	}

	for i, start := range startIndecies {
		pathsOut[i] = make([]int, 0)
		//i+1 is the thread index. It is the path index plus one so every ID is greater than 0, because 0 is unclaimed.
		visited.Add(1)
		go dfsRecursiveSynchronizer(g.nodes[start], val, &visited, &pathsOut[i], i+1)
	}
	visited.Wait()
	return visited.found, &pathsOut
}

func dfsRecursiveSynchronizer(n *node, val int, visited *visited, myPath *[]int, index int) {
	defer visited.Done()
	dfsRecursiveMultithreaded(n, val, visited, myPath, index)
}

// dfsRecursive returns true if the value is found and false if the value is not
// It writes to pathsOut its solution based on the index passed in from dfsRecursive
func dfsRecursiveMultithreaded(n *node, val int, visited *visited, myPath *[]int, index int) bool {
	visited.Lock()
	//end the search if another path found the target value
	if visited.found != -1 {
		visited.Unlock()
		return true
	}
	//end the search if this node has already been claimed
	if (*(*visited).v)[n.index] != 0 {
		visited.Unlock()
		return false
	}
	//end the search if this node contains the target value
	if n.val == val {
		visited.found = index
		visited.Unlock()
		return true
	}
	//otherwise, claim the node
	(*(*visited).v)[n.index] = index - 1
	visited.Unlock()

	//append the path as it goes on, not in reverse, to show all searching strands
	*myPath = append(*myPath, n.index)

	for _, currentNeighbor := range n.neighbors {
		if dfsRecursiveMultithreaded(currentNeighbor.n, val, visited, myPath, index) {
			return true
		}
	}

	return false
}
