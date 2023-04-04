package main

import (
	"fmt"
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
		out += "border-top-style: solid; border-top-width: thick; "
	} else {
		out += "border-top-style: none; "
	}

	if node.down {
		out += "border-bottom-style: solid; border-bottom-width: thick; "
	} else {
		out += "border-bottom-style: none; "
	}

	if node.right {
		out += "border-right-style: solid; border-right-width: thick; "
	} else {
		out += "border-right-style: none; "
	}

	if node.left {
		out += "border-left-style: solid; border-left-width: thick; "
	} else {
		out += "border-left-style: none; "
	}

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

func pathToJs(m *maze, path *[]int) template.JS {
	if len(*path) == 0 {
		return template.JS("[]")
	}

	out := "["

	for i := 0; i < len(*path); i++ {
		row, col := getMazeCoords(m, (*path)[i])
		out += "[" + strconv.Itoa(row) + ", " + strconv.Itoa(col) + "], "
	}

	//cut off ending comma
	out = out[:len(out)-2]
	out += "]"

	return template.JS(out)
}

func pathsToJs(m *maze, paths *[][]int) template.JS {
	out := template.JS("[")
	for _, path := range *paths {
		out += pathToJs(m, &path)
		out += ", "
	}
	//cut off ending comma
	out = out[:len(out)-2]
	out += template.JS("]")
	return out
}

var tpl = template.Must(template.ParseFiles("templates/index.html"))

type TemplateData struct {
	// Template data structs must have exported names so the template Executer can read them.

	MStyles   [][]template.CSS
	MPath     template.JS
	MBestPath template.JS
	TickSpeed template.JS
}

// fillTemplateData executes the search algorithms and processes their results.
func fillTemplateData(m *maze, animate bool, tickSpeed int) *TemplateData {
	//TODO sometimes DFS cheats on the right side of a large maze - not sure if it is a visual bug or a data structure bug

	// Run DFS to find the best solution
	dfsOk, dfsPath := dfs(&m.g, 3, 0)
	var bestPath template.JS
	if dfsOk {
		if !animate {
			m.fillPath(*dfsPath)
			bestPath = template.JS("[]")
		} else {
			bestPath = pathToJs(m, dfsPath)
		}
	} else {
		print("No Valid DFS\n")
		bestPath = template.JS("[]")
	}

	startI := getSeekerLocations(m, 4)

	// Run Multithreaded DFS to find a solution
	multithreadedOk, paths := dfsMultithreaded(&m.g, 3, startI)
	var mazePath template.JS
	if multithreadedOk != -1 {
		if !animate {
			mazePath = template.JS("[]")
		} else {
			mazePath = pathsToJs(m, paths)
		}
	} else {
		print("No Valid DFS\n")
		mazePath = template.JS("[]")
	}

	mazeValues := mazeToSlice(m)
	mazeStyles := mazeSliceToStyle(mazeValues)

	// Convert data into the types that a template can take
	tplData := TemplateData{
		MStyles:   mazeStyles,
		MPath:     mazePath,
		MBestPath: bestPath,
		TickSpeed: template.JS(strconv.Itoa(tickSpeed)),
	}
	return &tplData
}

func makeMaze(w http.ResponseWriter, r *http.Request, algo int) {
	// algo defines the maze generation algorithm
	// 1 = random
	// 2 = DFS

	// Parse values from GET request, if they exist
	width := 40
	height := 20
	tickSpeed := 1
	animate := true
	inputWidth, err := strconv.Atoi(r.URL.Query().Get("width"))
	if err == nil && inputWidth > 2 {
		width = inputWidth
	}
	inputHeight, err := strconv.Atoi(r.URL.Query().Get("height"))
	if err == nil && inputHeight > 2 {
		height = inputHeight
	}
	inputTickSpeed, err := strconv.Atoi(r.URL.Query().Get("tickSpeed"))
	if err == nil && inputTickSpeed > 0 {
		tickSpeed = inputTickSpeed
	}
	inputAnimate := r.URL.Query().Get("animate")
	if inputAnimate == "f" {
		animate = false
	}

	// Init maze with a given algorithm
	maze := initMaze(height, width)
	maze.SetSquare(height-1, width-1, 3)
	switch algo {
	case 1:
		density := 15
		inputDensity, err := strconv.Atoi(r.URL.Query().Get("density"))
		if err == nil && inputDensity > 0 {
			density = inputDensity
		}
		randomizeMaze(maze, density)
	case 2:
		createDFSMaze(maze)
	}

	tplData := fillTemplateData(maze, animate, tickSpeed)

	err = tpl.Execute(w, tplData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func dfsHandler(w http.ResponseWriter, r *http.Request) {
	makeMaze(w, r, 2)
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	makeMaze(w, r, 1)
}

func main() {
	/*maze := initMaze(100, 300)
	maze.SetSquare(300-1, 300-1, 3)
	createDFSMaze(maze)

	// Run DFS to find the solution
	startIndecies := []int{0, 30149, 60299}
	ok, paths := dfsMultithreaded(&maze.g, 3, startIndecies)
	if ok != -1 {
		fmt.Println(paths)
	} else {
		print("No Valid DFS\n")
	}*/

	mux := http.NewServeMux()

	//mux.HandleFunc("/", dfsHandler)

	mux.HandleFunc("/dfs", dfsHandler)
	mux.HandleFunc("/random", randomHandler)

	port := "3000"
	http.ListenAndServe(":"+port, mux)
}
