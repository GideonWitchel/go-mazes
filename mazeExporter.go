package main

import "fmt"

type mazeNode struct {
	val   int
	up    bool
	down  bool
	right bool
	left  bool
}

func makeMazeNode(m *maze, row int, col int) mazeNode {
	var newMazeNode mazeNode
	newMazeNode.val = m.g.nodes[getMazeIndex(m, row, col)].val

	if row == 0 || m.g.HasEdge(getMazeIndex(m, row, col), getMazeIndex(m, row-1, col)) {
		newMazeNode.up = true
	}
	if row == m.height-1 || m.g.HasEdge(getMazeIndex(m, row, col), getMazeIndex(m, row+1, col)) {
		newMazeNode.down = true
	}
	if col == m.width-1 || m.g.HasEdge(getMazeIndex(m, row, col), getMazeIndex(m, row, col+1)) {
		newMazeNode.right = true
	}
	if col == 0 || m.g.HasEdge(getMazeIndex(m, row, col), getMazeIndex(m, row, col-1)) {
		newMazeNode.left = true
	}

	return newMazeNode
}

func mazeToSlice(m *maze) [][]mazeNode {
	//contains the values of every node
	nodes := make([][]mazeNode, m.height)
	for row := 0; row < m.height; row++ {
		newRow := make([]mazeNode, m.width)
		for col := 0; col < m.width; col++ {
			newRow[col] = makeMazeNode(m, row, col)
		}
		nodes[row] = newRow
	}
	return nodes
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
