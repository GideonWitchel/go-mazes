package main

import (
	"sync"
)

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
	// TODO Swap queue to a channel.

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
	// TODO Swap queue to a channel

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

type bfsShared struct {
	visited []bool
	// parents is an array of nodes' parents for finding the optimal path
	parents []int
	// queue is a queue of the next values to test
	queue []int
	// found stores the index of the solution node. If it is -1, the solution has not been found.
	solution int
	sync.Mutex
}

type safePaths struct {
	p []*[]int
	sync.Mutex
}

type threadManager struct {
	maxThreads int
	// WaitGroup for bfsMultithreaded to know all threads are done
	sync.WaitGroup
}

type threadDone struct {
	//separate data structure for the manager to know when it needs to start new threads.
	numThreads int
	sync.Mutex
}

func bfsMultithreaded(g *graph, val int, startIndex int, maxThreads int) (exists bool, paths *[]*[]int, solution *[]int) {
	if maxThreads < 1 || startIndex < 0 {
		return false, nil, nil
	}

	pathsOut := safePaths{
		p: make([]*[]int, 0),
	}
	solutionOut := make([]int, 0)
	v := make([]bool, len(g.nodes), len(g.nodes))
	p := make([]int, len(g.nodes), len(g.nodes))
	q := make([]int, 0)
	// TODO Swap queue to a channel

	q = append(q, startIndex)
	v[startIndex] = true
	p[startIndex] = -1

	bfsData := bfsShared{
		visited:  v,
		parents:  p,
		queue:    q,
		solution: -1,
	}

	tManager := threadManager{
		maxThreads: maxThreads,
	}

	tManager.Add(1)
	go bfsMultithreadedManager(&tManager, g, &bfsData, val, &pathsOut)
	tManager.Wait()

	// Backtrack through parents to find the shortest path.
	if bfsData.solution != -1 {
		// Start at the goal node.
		i := bfsData.solution
		for i != -1 {
			solutionOut = append(solutionOut, i)
			i = bfsData.parents[i]
		}
	}

	// Ensure paths are caught up
	pathsOut.Lock()
	pathsOut.Unlock()

	return bfsData.solution != -1, &pathsOut.p, &solutionOut
}

func bfsMultithreadedManager(threadManager *threadManager, g *graph, bfsData *bfsShared, val int, paths *safePaths) {
	defer threadManager.Done()

	threadData := threadDone{
		numThreads: 0,
	}

	for {
		threadData.Lock()
		if threadData.numThreads < threadManager.maxThreads {
			threadData.Unlock()

			bfsData.Lock()
			// The second condition is needed to deap with graphs without solution but breaks real mazes.
			// I'm pretty sure the best fix is using channels.
			if bfsData.solution == -1 && len(bfsData.queue) > 0 {
				bfsData.Unlock()

				newPath := make([]int, 0)
				paths.Lock()
				paths.p = append(paths.p, &newPath)
				paths.Unlock()

				threadData.Lock()
				threadData.numThreads++
				threadData.Unlock()

				go bfsMultithreadedSubprocess(&threadData, g, bfsData, &newPath, val)

			} else {
				bfsData.Unlock()
				return
			}

		} else {
			threadData.Unlock()
		}
	}
}

func killThread(threadData *threadDone) {
	threadData.Lock()
	threadData.numThreads--
	threadData.Unlock()
}

func bfsMultithreadedSubprocess(threadData *threadDone, g *graph, data *bfsShared, pathOut *[]int, val int) {
	defer killThread(threadData)

	for {
		data.Lock()
		//end if the search is over or the solution has been found
		if len(data.queue) == 0 || data.solution != -1 {
			data.Unlock()
			return
		}

		currentNode := g.nodes[data.queue[0]]
		data.queue = data.queue[1:]
		data.Unlock()

		*pathOut = append(*pathOut, currentNode.index)

		if currentNode.val == val {
			data.Lock()
			data.solution = currentNode.index
			data.Unlock()
			// Do not overwrite the solution.
			*pathOut = (*pathOut)[:len(*pathOut)-1]
			return
		}

		for _, currentNeighbor := range currentNode.neighbors {
			data.Lock()
			if !data.visited[currentNeighbor.n.index] {
				data.visited[currentNeighbor.n.index] = true
				data.queue = append(data.queue, currentNeighbor.n.index)
				data.parents[currentNeighbor.n.index] = currentNode.index
			}
			data.Unlock()
		}
	}
}
