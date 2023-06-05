package main

import (
	ms "go-mazes/mazesrv"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", ms.MakeMazeResponse)

	port := "3000"
	http.ListenAndServe(":"+port, mux)
}
