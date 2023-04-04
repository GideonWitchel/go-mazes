package main

import (
	"fmt"
	"sync"
	"time"
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
	v     *[]bool
	found int
	sync.Mutex
}

type paths struct {
	p *[][]int
	sync.Mutex
}

// dfsMultithreaded finds a value in a graph using a number of simultaneous dfs searchs with a shared visited list.
// exists is an index which specifies which search ended up finding the value in the paths array.
// If exists is -1, there is valid path to the solution from any starting index.
func dfsMultithreaded(g *graph, val int, startIndecies []int) (exists int, p *[][]int) {
	pathsOut := make([][]int, len(startIndecies), len(startIndecies))
	pathsStruct := paths{
		p: &pathsOut,
	}
	visitedArray := make([]bool, len(g.nodes), len(g.nodes))
	visited := visited{
		v:     &visitedArray,
		found: -1,
	}

	for i, start := range startIndecies {
		//println("Starting thread " + strconv.Itoa(i))
		//visited.Lock()
		go dfsRecursiveMultithreaded(g.nodes[start], val, &visited, &pathsStruct, i)
		//time.Sleep(time.Nanosecond)
	}
	// TODO Actually check if subprocesses are done
	time.Sleep(time.Second)
	return visited.found, pathsStruct.p
}

// dfsRecursive returns true if the value is found and false if the value is not
// It writes to pathsOut its solution based on the index passed in from dfsRecursive
func dfsRecursiveMultithreaded(n *node, val int, visited *visited, pathsOut *paths, index int) {
	//visited.Unlock()
	//println(strconv.Itoa(index) + " | " + strconv.Itoa(n.index) + " : " + " is starting")
	//end the search if another path found the target value
	visited.Lock()
	if visited.found != -1 {
		//println(strconv.Itoa(index) + " | " + strconv.Itoa(n.index) + " : " + " is done")
		return
	}
	//otherwise, start this node's search
	(*(*visited).v)[n.index] = true
	//println(strconv.Itoa(index) + " | " + strconv.Itoa(n.index) + " : " + " is starting a search")
	visited.Unlock()

	//append the path as it goes on, not in reverse, to show all searching strands
	pathsOut.Lock()
	(*pathsOut.p)[index] = append((*pathsOut.p)[index], n.index)
	//println(strconv.Itoa(index) + " | " + strconv.Itoa(n.index) + " : " + " appended a path")
	pathsOut.Unlock()

	//end search if found the value
	if n.val == val {
		visited.Lock()
		visited.found = index
		//println(strconv.Itoa(index) + " | " + strconv.Itoa(n.index) + " : " + " is done")
		visited.Unlock()
		return
	}

	for _, currentNeighbor := range n.neighbors {
		visited.Lock()
		if !(*(*visited).v)[currentNeighbor.n.index] {
			visited.Unlock()
			//println(strconv.Itoa(index) + " | " + strconv.Itoa(n.index) + " : " + " found a valid node to move on to")

			dfsRecursiveMultithreaded(currentNeighbor.n, val, visited, pathsOut, index)

			//if found the value, end search
			visited.Lock()
			if visited.found != -1 {
				visited.Unlock()
				//println(strconv.Itoa(index) + " | " + strconv.Itoa(n.index) + " : " + " is done2")
				return
			}
			visited.Unlock()
		} else {
			//println(strconv.Itoa(index) + " | " + strconv.Itoa(n.index) + " : " + " found an invalid node")
			visited.Unlock()
		}
	}
}
