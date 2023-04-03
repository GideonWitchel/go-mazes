package main

import (
	"html/template"
	"net/http"
)

var tpl = template.Must(template.ParseFiles("templates/index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	maze := initMaze(5, 10)
	randomizeMaze(maze, 4)
	maze.SetSquare(4, 9, 3)

	ok, path := dfs(&maze.g, 3, 0)
	if ok {
		maze.fillPath(*path)
	}

	mazeVals := mazeToSlice(maze)
	mazeStyles := mazeSliceToStyle(mazeVals)

	tpl.Execute(w, mazeStyles)
}

func toStyle(node mazeNode) template.CSS {
	out := ""

	switch node.val {
	case 0:
		out += "background-color: white; "
	case 1:
		out += "background-color: blue; "
	case 2:
		out += "background-color: green; "
	case 3:
		out += "background-color: yellow; "
	}

	if node.up {
		out += "border-top-color: black; "
	} else {
		out += "border-top-color: LightGray; "
	}

	if node.down {
		out += "border-bottom-color: black; "
	} else {
		out += "border-bottom-color: LightGray; "
	}

	if node.right {
		out += "border-right-color: black; "
	} else {
		out += "border-right-color: LightGray; "
	}

	if node.left {
		out += "border-left-color: black; "
	} else {
		out += "border-left-color: LightGray; "
	}

	out += "border-style: solid;"
	return template.CSS(out)
}

func mazeSliceToStyle(mazeVals [][]mazeNode) [][]template.CSS {
	mazeStyles := make([][]template.CSS, len(mazeVals))
	for row := range mazeVals {
		newRow := make([]template.CSS, len(mazeVals[0]))
		for col := range mazeVals[row] {
			newRow[col] = toStyle(mazeVals[row][col])
		}
		mazeStyles[row] = newRow
	}
	return mazeStyles
}

func main() {
	/*maze := initMaze(5, 10)
	randomizeMaze(maze, 2)
	maze.SetSquare(4, 9, 3)
	maze.Print()
	print("\n\n")

	printSolution(maze)*/

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	port := "3000"
	http.ListenAndServe(":"+port, mux)
}
