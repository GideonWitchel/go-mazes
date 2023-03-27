package main

import (
	"fmt"
	"math/rand"
)

/*
Each node in a maze has an integer value representing its state
0 is empty
1 is filled by the user (which overrides any empty node)
2 is filled by the solution (which overrides any user filled node)
3 is filled by the goal (which overrides any user or solution node)
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
	case 3:
		fmt.Print("☆")
	default:
		fmt.Println("Impossible Value")
	}
}

// TODO seems to render walls as empty and empty as walls?
// - for now just manually swapped the symbols
func (m *maze) Print() {
	index := 0
	for row := 0; row < m.height; row++ {
		for col := 0; col < m.width-1; col++ {
			printNode(m.g.nodes[index].val)

			if m.g.HasEdge(index, index+1) {
				fmt.Print(" ")
			} else {
				fmt.Print("|")
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
				fmt.Print("  ")
			} else {
				fmt.Print("— ")
			}
			index++
		}
		fmt.Print("\n")
	}
}

// printSolution assumes the goal is 3
func printSolution(m *maze) {
	ok, path := dfs(&m.g, 3, 0)
	if ok {
		//reverse path to draw from starting location
		//skip the first item which overwrites the solution
		forwardPath := make([]int, 0)
		for i := len(*path) - 1; i >= 1; i-- {
			forwardPath = append(forwardPath, (*path)[i])
		}

		for _, i := range forwardPath {
			row, col := getMazeCoords(m, i)
			m.SetSquare(row, col, 2)
			m.Print()
			print("\n\n")
		}
	} else {
		print("No Valid DFS\n")
	}
}

func getMazeIndex(m *maze, row int, col int) int {
	return row*(m.width) + col
}

// getMazeCoords returns (row, col) from index
func getMazeCoords(m *maze, index int) (int, int) {
	return index / m.width, index % m.width
}

func (m *maze) SetSquare(row int, col int, val int) {
	m.g.SetNode(getMazeIndex(m, row, col), val)
}

func (m *maze) SetWall(row1 int, col1 int, row2 int, col2 int, remove bool) {
	//For direction: 1 = north, 2 = south, 3 = east, 4 = west
	//for removing: true = remove, false = add

	/*
		differential := 0
		switch direction {
		case 1:
			differential = -1 * m.width
		case 2:
			differential = m.width
		case 3:
			differential = 1
		case 4:
			differential = -1
		default:
			//do nothing if invalid direction
			return
		}
	*/

	//TODO error checking to make sure the nodes are next to each other
	index1 := getMazeIndex(m, row1, col1)
	index2 := getMazeIndex(m, row2, col2)

	if remove {
		m.g.RemoveEdge(index1, index2)
	} else {
		m.g.AddEdge(index1, index2)
	}
}

// randomizeMaze randomizes every wall in the maze
// increased sparsity means fewer walls
func randomizeMaze(m *maze, sparsity int) {
	for row := 0; row < m.height-1; row++ {
		for col := 0; col < m.width-1; col++ {
			//randomize edge below
			m.SetWall(row, col, row+1, col, rand.Intn(sparsity) == 1)
			//randomize edge to the right
			m.SetWall(row, col, row, col+1, rand.Intn(sparsity) == 1)
		}
		//randomize just below for the right column
		m.SetWall(row, m.width-1, row+1, m.width-1, rand.Intn(sparsity) == 1)
	}
	//randomize just to the right for the bottom row, and nothing for the bottom right node
	for col := 0; col < m.width-1; col++ {
		m.SetWall(m.height-1, col, m.height-1, col+1, rand.Intn(sparsity) == 1)
	}
}

func main() {
	maze := initMaze(5, 10)
	randomizeMaze(maze, 4)
	maze.SetSquare(4, 9, 3)
	maze.Print()
	print("\n\n")

	printSolution(maze)
}
