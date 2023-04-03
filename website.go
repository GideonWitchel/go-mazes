package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var tpl = template.Must(template.ParseFiles("templates/index.html"))

var mazeTemplate struct {
	values [][]int
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func main() {
	maze := initMaze(5, 10)
	randomizeMaze(maze, 4)
	maze.SetSquare(4, 9, 3)
	maze.Print()
	print("\n\n")

	printSolution(maze)
	nodes, edges := mazeToSlice(maze)
	fmt.Println(nodes)
	fmt.Println(edges)

	/*mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	port := "3000"
	http.ListenAndServe(":"+port, mux)*/
}
