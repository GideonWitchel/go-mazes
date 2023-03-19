package main

import "fmt"

/*
Each node in a maze has an integer value representing its state
0 is empty
TODO 1 is filled by the user (which overrides any empty node)
TODO 2 is filled by the solution (which overrides any user filled node)
*/
type maze struct {
	g      graph
	height int
	width  int
}

func initMaze(height int, width int) *maze {
	//cannot be smaller than 2x2
	if height < 2 || width < 2 {
		return nil
	}

	totalNodes := height * width
	var m maze
	g := &(m.g)
	for i := 0; i < totalNodes; i++ {
		g.AddNode(0)
	}

	//height is the number of rows, width is the number of columns
	//initialize with no walls, so every node is connected to the node to its North, South, East, and West
	//nodes fill left to right, then top to bottom
	//edges are added on  the bottom and right side, so the bottom row and right column are ignored
	index := 0
	for row := 0; row < height-1; row++ {
		for col := 0; col < width-1; col++ {
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
	for col := 0; col < width-1; col++ {
		g.AddEdge(index, index+1)
		index++
	}

	m.height = height
	m.width = width
	return &m
}

func printNode(val int) {
	switch val {
	case 0:
		fmt.Print("□")
	case 1:
		fmt.Print("■")
	case 2:
		fmt.Print("⧇")
	default:
		fmt.Println("Impossible Value")
	}
}

func (m *maze) Print() {
	index := 0
	for row := 0; row < m.height; row++ {
		for col := 0; col < m.width-1; col++ {
			printNode(m.g.nodes[index].val)

			if m.g.HasEdge(index, index+1) {
				fmt.Print("|")
			} else {
				fmt.Print(" ")
			}
			index++
		}
		//right column, doesn't check for edge
		printNode(m.g.nodes[index].val)
		fmt.Print("\n")

		//print edge values between rows, unless it's the last row
		if row == m.height-1 {
			return
		}
		index -= m.width - 1
		for col := 0; col < m.width; col++ {
			if m.g.HasEdge(index, index+m.width) {
				fmt.Print("— ")
			} else {
				fmt.Print("  ")
			}
			index++
		}
		fmt.Print("\n")
	}
}

func main() {
	maze := initMaze(5, 10)
	maze.g.Print()
	maze.g.RemoveEdge(10, 20)
	maze.g.RemoveEdge(10, 11)
	maze.g.RemoveEdge(10, 12)
	maze.g.RemoveEdge(9, 10)

	maze.Print()

}
