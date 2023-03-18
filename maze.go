package main

import "fmt"

/*
Each node in a maze has an integer value representing its state
0 is empty
TODO 1 is filled by the user (which overrides any empty node)
TODO 2 is filled by the solution (which overrides any user filled node)
*/

func initMaze(height int, width int) *graph {
	//height is the number of rows, width is the number of columns
	totalNodes := height * width
	var g graph
	for i := 0; i < totalNodes; i++ {
		g.addNode(0)
	}

	//initialize with no walls, so every node is connected to the node to its North, South, Eeast, and Rest
}

func main() {
	fmt.Println("hi")
}
