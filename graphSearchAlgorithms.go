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

type dfsShared struct {
	visited []bool
	// found represents the path ID (thread ID - 1) of whoever found the solution node.
	// -1 means the solution has not been found
	found int
	sync.Mutex
	sync.WaitGroup
}

// dfsMultithreaded finds a value in a graph using a number of simultaneous dfs searches with a shared visited list.
// exists is an index which specifies which search ended up finding the value in the paths array.
// If exists is -1, there is valid path to the solution from any starting index.
func dfsMultithreaded(g *graph, val int, startIndecies []int) (exists int, p *[][]int) {
	pathsOut := make([][]int, len(startIndecies), len(startIndecies))

	visitedArray := make([]bool, len(g.nodes), len(g.nodes))
	dfsData := dfsShared{
		visited: visitedArray,
		found:   -1,
	}

	for i, start := range startIndecies {
		pathsOut[i] = make([]int, 0)
		dfsData.Add(1)
		go dfsRecursiveSynchronizer(g.nodes[start], val, &dfsData, &pathsOut[i], i)
	}

	dfsData.Wait()
	return dfsData.found, &pathsOut
}

func dfsRecursiveSynchronizer(n *node, val int, dfsData *dfsShared, myPath *[]int, index int) {
	defer dfsData.Done()
	dfsRecursiveMultithreaded(n, val, dfsData, myPath, index)
}

// dfsRecursive returns true if the value is found and false if the value is not.
// It writes to pathsOut its solution based on the index passed in from dfsRecursive.
func dfsRecursiveMultithreaded(n *node, val int, dfsData *dfsShared, myPath *[]int, index int) bool {
	dfsData.Lock()
	// End the search if another path found the target value.
	if dfsData.found != -1 {
		dfsData.Unlock()
		return true
	}
	// End the search if this node has already been claimed.
	if (*dfsData).visited[n.index] == true {
		dfsData.Unlock()
		return false
	}
	// End the search if this node contains the target value.
	if n.val == val {
		dfsData.found = index
		dfsData.Unlock()
		return true
	}
	// Otherwise, claim the node.
	(*dfsData).visited[n.index] = true
	dfsData.Unlock()

	// Append the path as it goes on, not in reverse, to show all searching strands.
	*myPath = append(*myPath, n.index)

	for _, currentNeighbor := range n.neighbors {
		if dfsRecursiveMultithreaded(currentNeighbor.n, val, dfsData, myPath, index) {
			return true
		}
	}

	return false
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

func bfsRecursive(g *graph, queue *[]int, val int, visited *[]bool, parents *[]int, pathOut *[]int) (success bool, valIndex int) {
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

	ok, index := bfsRecursive(g, queue, val, visited, parents, pathOut)
	if ok {
		return true, index
	}

	return false, -1
}

func bfsIterative(g *graph, val int, startIndex int) (exists bool, path *[]int, solution *[]int) {
	pathOut := make([]int, 0)
	solutionOut := make([]int, 0)
	visited := make([]bool, len(g.nodes), len(g.nodes))
	parents := make([]int, len(g.nodes), len(g.nodes))
	queue := make([]int, 0, 0)
	// TODO Swap queue's data structure from a slice to a linked list because a slice has awful performance for popping from the front.

	queue = append(queue, startIndex)
	visited[startIndex] = true
	parents[startIndex] = -1

	success := false
	valIndex := -1

	for len(queue) != 0 {
		currentNode := g.nodes[queue[0]]
		queue = queue[1:]
		pathOut = append(pathOut, currentNode.index)

		if currentNode.val == val {
			success = true
			valIndex = currentNode.index
			break
		}

		for _, currentNeighbor := range currentNode.neighbors {
			if !visited[currentNeighbor.n.index] {
				visited[currentNeighbor.n.index] = true
				queue = append(queue, currentNeighbor.n.index)
				parents[currentNeighbor.n.index] = currentNode.index
			}
		}
	}

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

type lockedPaths struct {
	p [][]int
	sync.Mutex
}

type bfsShared struct {
	visited []bool
	// parents is an array of nodes' parents for finding the optimal path
	parents []int
	// found stores the index of the solution node. If it is -1, the solution has not been found.
	solution int
	// queue is a queue of the next values to test
	queue []int
	// nextID is the next thread ID to use when making a new thread
	nextID int
	// numThreads is the number of active threads
	numThreads int
	sync.Mutex
	sync.WaitGroup
}

func bfsMultithreaded(g *graph, val int, startIndex int, maxThreads int) (exists bool, paths *[][]int, solution *[]int) {
	if maxThreads < 1 || startIndex < 0 {
		return false, nil, nil
	}

	pathsOut := lockedPaths{
		p: make([][]int, 0),
	}

	solutionOut := make([]int, 0)
	v := make([]bool, len(g.nodes), len(g.nodes))
	p := make([]int, len(g.nodes), len(g.nodes))
	q := make([]int, 0, 0)
	// TODO Swap queue's data structure from a slice to a linked list because a slice has awful performance for popping from the front.

	q = append(q, startIndex)
	v[startIndex] = true
	p[startIndex] = -1

	bfsData := bfsShared{
		visited:    v,
		parents:    p,
		queue:      q,
		solution:   -1,
		nextID:     0,
		numThreads: 0,
	}

	bfsData.Add(1)
	go bfsRecursiveSynchronizer(g, val, &bfsData, &pathsOut, maxThreads)
	bfsData.Wait()

	// Backtrack through parents to find the shortest path.
	if bfsData.solution != -1 {
		// Start at the goal node.
		i := bfsData.solution
		for i != -1 {
			solutionOut = append(solutionOut, i)
			i = bfsData.parents[i]
		}
	}

	return bfsData.solution != -1, &pathsOut.p, &solutionOut
}

func bfsRecursiveSynchronizer(g *graph, val int, data *bfsShared, paths *lockedPaths, maxThreads int) {
	defer data.Done()
	defer print("Done!")

	data.queue = append(data.queue, 0)
	data.parents[0] = -1

	for {
		//Always lock data before paths
		data.Lock()
		if data.numThreads < maxThreads {
			//check for termination
			if data.solution != -1 {
				data.Unlock()
				return
			}

			if len(data.queue) != 0 {
				paths.Lock()

				paths.p = append(paths.p, make([]int, 0))
				currID := data.nextID
				data.nextID += 1
				pathPointer := &paths.p[currID]
				data.numThreads++

				paths.Unlock()
				data.Unlock()

				go bfsRecursiveMultithreadedContainer(g, val, data, pathPointer, currID)
			} else {
				data.Unlock()
			}
		}
	}
}

func bfsRecursiveMultithreadedContainer(g *graph, val int, data *bfsShared, myPath *[]int, id int) {
	bfsRecursiveMultithreaded(g, val, data, myPath, id)
	data.Lock()
	data.numThreads--
	data.Unlock()
}

func bfsRecursiveMultithreaded(g *graph, val int, data *bfsShared, myPath *[]int, id int) {
	// Any thread can grab from the queue in any order, adding all the neighbors from their grabbed node

	if id == 1 {
		print("")
	}
	data.Lock()
	// End the search if the search has failed
	if len(data.queue) == 0 {
		data.Unlock()
		return
	}
	// End the search if another path found the target value.
	if data.solution != -1 {
		data.Unlock()
		return
	}

	currentNode := g.nodes[data.queue[0]]
	// End the search if this node has already been claimed.
	//if data.visited[currentNode.index] == true {
	//	data.Unlock()
	//	return
	//}

	// Otherwise, claim the node.
	data.queue = data.queue[1:]
	// End the search if this node contains the target value.
	if currentNode.val == val {
		data.solution = currentNode.index
		data.Unlock()
		return
	}
	data.Unlock()

	// Append the path as it goes on, not in reverse, to show all searching strands.
	*myPath = append(*myPath, currentNode.index)
	for _, currentNeighbor := range currentNode.neighbors {
		data.Lock()
		if !data.visited[currentNeighbor.n.index] {
			data.queue = append(data.queue, currentNeighbor.n.index)
			data.parents[currentNeighbor.n.index] = currentNode.index
		}
		data.Unlock()
	}

	bfsRecursiveMultithreaded(g, val, data, myPath, id)
}
