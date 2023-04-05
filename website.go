package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var tpl = template.Must(template.ParseFiles("templates/index.html"))

// makeMaze converts a http request into a http response
func makeMazeResponse(w http.ResponseWriter, r *http.Request, algo int) {
	timeStart := time.Now()
	// algo defines the maze generation algorithm
	// 1 = random
	// 2 = DFS

	// Parse values from GET request, if they exist
	width := 40
	height := 20
	tickSpeed := 1
	repeats := 1
	animate := true
	density := 15
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
	inputRepeats, err := strconv.Atoi(r.URL.Query().Get("repeats"))
	if err == nil && inputRepeats > 0 {
		repeats = inputRepeats
	}
	inputDensity, err := strconv.Atoi(r.URL.Query().Get("density"))
	if err == nil && inputDensity > 0 {
		density = inputDensity
	}

	// Init maze with a given algorithm
	maze := initMaze(height, width)
	maze.SetSquare(height-1, width-1, 3)
	switch algo {
	case 1:
		randomizeMaze(maze, density)
	case 2:
		createDFSMaze(maze)
	}

	tplData := fillTemplateData(maze, animate, tickSpeed, repeats)

	err = tpl.Execute(w, tplData)
	timeEnd := time.Now()
	fmt.Printf("Served! in %v ms or %v us\n", timeEnd.UnixMilli()-timeStart.UnixMilli(), timeEnd.UnixMicro()-timeStart.UnixMicro())
	if err != nil {
		fmt.Println(err)
		return
	}
}

func dfsHandler(w http.ResponseWriter, r *http.Request) {
	makeMazeResponse(w, r, 2)
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	makeMazeResponse(w, r, 1)
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
