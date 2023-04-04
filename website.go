package main

import (
	"html/template"
	"net/http"
	"strconv"
)

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
		out += "border-top-color: black; border-top-width: thick; "
	} else {
		out += "border-top-color: LightGray; "
	}

	if node.down {
		out += "border-bottom-color: black; border-bottom-width: thick; "
	} else {
		out += "border-bottom-color: LightGray; "
	}

	if node.right {
		out += "border-right-color: black; border-right-width: thick; "
	} else {
		out += "border-right-color: LightGray; "
	}

	if node.left {
		out += "border-left-color: black; border-left-width: thick; "
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

var tpl = template.Must(template.ParseFiles("templates/index.html"))

func makeMaze(w http.ResponseWriter, r *http.Request, algo int) {
	//algo defines the maze generation algorithm
	//1 = random
	//2 = DFS

	width := 40
	height := 20
	out, err := strconv.Atoi(r.URL.Query().Get("width"))
	if err == nil && out > 2 {
		width = out
	}
	out, err = strconv.Atoi(r.URL.Query().Get("height"))
	if err == nil && out > 2 {
		height = out
	}

	maze := initMaze(height, width)
	maze.SetSquare(height-1, width-1, 3)

	switch algo {
	case 1:
		density := 15
		out, err := strconv.Atoi(r.URL.Query().Get("density"))
		if err == nil && out > 0 {
			density = out
		}
		randomizeMaze(maze, density)
	case 2:
		createDFSMaze(maze)
	}

	ok, path := dfs(&maze.g, 3, 0)
	if ok {
		maze.fillPath(*path)
	} else {
		print("No Valid DFS\n")
	}

	mazeValues := mazeToSlice(maze)
	mazeStyles := mazeSliceToStyle(mazeValues)

	tpl.Execute(w, mazeStyles)
}

func dfsHandler(w http.ResponseWriter, r *http.Request) {
	makeMaze(w, r, 2)
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	makeMaze(w, r, 1)
}

func main() {
	mux := http.NewServeMux()

	//mux.HandleFunc("/", dfsHandler)

	mux.HandleFunc("/dfs", dfsHandler)
	mux.HandleFunc("/random", randomHandler)

	port := "3000"
	http.ListenAndServe(":"+port, mux)
}
