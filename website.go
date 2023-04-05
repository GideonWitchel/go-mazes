package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"time"
)

var tpl = template.Must(template.ParseFiles("templates/index.html"))

// makeMaze converts a http request into a http response
func makeMazeResponse(w http.ResponseWriter, r *http.Request, generationAlgorithm int) {
	timeStart := time.Now()
	// generationAlgorithm defines the maze generation algorithm
	// 1 = random
	// 2 = DFS

	// Default configuration variables
	width := 180
	height := 90
	tickSpeed := 1
	repeats := 20
	density := 15
	solve := "bfs"

	// Parse values from GET request, if they exist
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
	inputRepeats, err := strconv.Atoi(r.URL.Query().Get("repeats"))
	if err == nil && inputRepeats > 0 {
		repeats = inputRepeats
	}
	inputDensity, err := strconv.Atoi(r.URL.Query().Get("density"))
	if err == nil && inputDensity > 0 {
		density = inputDensity
	}
	inputSolve := r.URL.Query().Get("solve")
	if inputSolve == "dfs" {
		solve = "dfs"
	}

	// TODO there are impossible patterns (closed off areas) on large mazes - not sure if it is a visual bug or a data structure bug

	// Init maze with a given algorithm
	maze := initMaze(height, width)
	maze.SetSquare(height-1, width-1, 3)
	switch generationAlgorithm {
	case 1:
		randomizeMaze(maze, density)
	case 2:
		createDFSMaze(maze)
	}

	// Solve maze with a given algorithm
	var tplData *TemplateData
	switch solve {
	case "bfs":
		tplData = fillTemplateBFS(maze, tickSpeed, repeats)
	case "dfs":
		tplData = fillTemplateDFS(maze, tickSpeed, repeats)
	}

	err = tpl.Execute(w, tplData)
	if err != nil {
		fmt.Println(err)
		return
	}

	timeEnd := time.Now()
	// Manually calculate times with ns to have control over rounding
	timeEndNs := timeEnd.UnixNano() - timeStart.UnixNano()
	timeEndUs := float64(timeEndNs) / 1000.0
	timeEndMs := timeEndUs / 1000.0
	fmt.Printf("Served! in %.0f ms and %.0f us\n", timeEndMs, timeEndUs-(math.Floor(timeEndMs)*1000.0))
}

func dfsHandler(w http.ResponseWriter, r *http.Request) {
	makeMazeResponse(w, r, 2)
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	makeMazeResponse(w, r, 1)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", dfsHandler)
	mux.HandleFunc("/dfs", dfsHandler)
	mux.HandleFunc("/random", randomHandler)

	port := "3000"
	http.ListenAndServe(":"+port, mux)
}
