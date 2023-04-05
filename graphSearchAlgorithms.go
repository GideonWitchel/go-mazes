package main

import "sync"

// DFS finds the first node with a given value and returns:
// - a boolean which is true if the value is accessible
// - a slice of indexes with the order of nodes to get there, starting with the node of the desired value and ending with the starting node
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

func reverseSlice(in *[]int) *[]int {
	out := make([]int, len(*in))

	for i, val := range *in {
		newIndex := len(*in) - i - 1
		out[newIndex] = val
	}

	return &out
}

// BFS finds the first node with a given value and returns:
// - a boolean which is true if the value is accessible
// - a path slice of indexes covering everything the search algorithm covered, in the order they were visited
// - a solution slice of indexes with the order of nodes to efficiently get to the value, starting with the node of the desired value and ending with the starting node
func bfs(g *graph, val int, startIndex int) (exists bool, path *[]int, solution *[]int) {
	pathOut := make([]int, 0)
	solutionOut := make([]int, 0)
	visited := make([]bool, len(g.nodes), len(g.nodes))
	parents := make([]int, len(g.nodes), len(g.nodes))
	queue := make([]int, 0, 0)
	// TODO Swap queue's data structure from a slice to a linked list because a slice has awful performance for popping from the front.

	queue = append(queue, startIndex)
	visited[startIndex] = true
	parents[startIndex] = -1

	success, valIndex := bfsRecursive(g, &queue, val, &visited, &parents, &pathOut)

	// Backtrack through parents to find the shortest path.
	if success {
		// Start at the goal node.
		i := valIndex
		for i != -1 {
			solutionOut = append(solutionOut, i)
			i = parents[i]
		}
	}

	return success, &pathOut, &solutionOut
}

// bfsRecursive returns a boolean success value and the index of the targeted value
func bfsRecursive(g *graph, queue *[]int, val int, visited *[]bool, parents *[]int, pathOut *[]int) (bool, int) {
	if len(*queue) == 0 {
		return false, -1
	}

	currentNode := g.nodes[(*queue)[0]]
	*queue = (*queue)[1:]
	*pathOut = append(*pathOut, currentNode.index)

	if currentNode.val == val {
		return true, currentNode.index
	}

	for _, currentNeighbor := range currentNode.neighbors {
		if !(*visited)[currentNeighbor.n.index] {
			(*visited)[currentNeighbor.n.index] = true
			*queue = append(*queue, currentNeighbor.n.index)
			(*parents)[currentNeighbor.n.index] = currentNode.index
		}
	}

	success, valIndex := bfsRecursive(g, queue, val, visited, parents, pathOut)
	if success {
		return true, valIndex
	}

	return false, -1
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
		// i+1 is the thread index. It is the path index plus one so every ID is greater than 0, because 0 is unclaimed.
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

// dfsRecursive returns true if the value is found and false if the value is not.
// It writes to pathsOut its solution based on the index passed in from dfsRecursive.
func dfsRecursiveMultithreaded(n *node, val int, visited *visited, myPath *[]int, index int) bool {
	visited.Lock()
	// End the search if another path found the target value.
	if visited.found != -1 {
		visited.Unlock()
		return true
	}
	// End the search if this node has already been claimed.
	if (*(*visited).v)[n.index] != 0 {
		visited.Unlock()
		return false
	}
	// End the search if this node contains the target value.
	if n.val == val {
		visited.found = index
		visited.Unlock()
		return true
	}
	// Otherwise, claim the node.
	(*(*visited).v)[n.index] = index - 1
	visited.Unlock()

	// Append the path as it goes on, not in reverse, to show all searching strands.
	*myPath = append(*myPath, n.index)

	for _, currentNeighbor := range n.neighbors {
		if dfsRecursiveMultithreaded(currentNeighbor.n, val, visited, myPath, index) {
			return true
		}
	}

	return false
}
