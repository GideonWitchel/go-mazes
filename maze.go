package main

/*
Each node in a maze has an integer value representing its state
0 is empty
TODO 1 is filled by the user (which overrides any empty node)
TODO 2 is filled by the solution (which overrides any user filled node)
*/

func initMaze(height int, width int) *graph {
	//cannot be smaller than 2x2
	if height < 2 || width < 2 {
		return nil
	}

	totalNodes := height * width
	var g graph
	for i := 0; i < totalNodes; i++ {
		g.AddNode(0)
	}

	//height is the number of rows, width is the number of columns
	//initialize with no walls, so every node is connected to the node to its North, South, East, and West
	//nodes fill left to right, then top to bottom
	//edges are added on  the bottom and right side, so the bottom row and right column are ignored
	index := 0
	for row := 0; row < width-1; row++ {
		for col := 0; col < height-1; col++ {
			//edge to below
			g.AddEdge(index, index+width)
			//edge to the right
			g.AddEdge(index, index+1)
			index++
		}
		//add just below for the right column
		g.AddEdge(index, index+width)
		index++
	}
	//add just to the right for the bottom row, and nothing for the bottom right node
	for col := 0; col < height-1; col++ {
		g.AddEdge(index, index+1)
		index++
	}

	return &g
}

func main() {
	maze := initMaze(3, 3)
	maze.Print()
}
