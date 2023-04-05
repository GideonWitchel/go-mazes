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
func makeMazeResponse(w http.ResponseWriter, r *http.Request, algo int) {
	timeStart := time.Now()
	// algo defines the maze generation algorithm
	// 1 = random
	// 2 = DFS

	// Default configuration variables
	width := 180
	height := 90
	tickSpeed := 1
	repeats := 20
	animate := true
	density := 15

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

	// TODO there are impossible patterns (closed off areas) on large mazes - not sure if it is a visual bug or a data structure bug

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
