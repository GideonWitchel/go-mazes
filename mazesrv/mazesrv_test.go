package mazesrv_test

import (
	"github.com/stretchr/testify/assert"
	"go-mazes/maze"
	ms "go-mazes/mazesrv"
	"testing"
)

func TestMaze(t *testing.T) {
	// Request maze from server
	arg := ms.MazeRequest{
		Height:      100,
		Width:       100,
		GenerateAlg: maze.GEN_DFS,
		SolveAlg:    maze.SOLVE_BFS_MULTI,
		Repeats:     50,
	}
	res := ms.MazeResponse{}
	err := ms.GetMaze(&arg, &res)
	assert.Nil(t, err, "Maze RPC call failed with arg: %v and err: %v", arg, err)
	print("Maze Output: %v", res.Webpage)
}
