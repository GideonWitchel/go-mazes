package mazesrv

import (
	"bytes"
	"errors"
)

type SrvMaze struct {
	sid string
}

type MazeRequest struct {
	Width       uint32
	Height      uint32
	Density     uint32
	TickSpeed   uint32
	Repeats     uint32
	GenerateAlg string
	SolveAlg    string
	StartIndex  uint64
}

type MazeResponse struct {
	Webpage string
}

func mkErr(message string) error {
	return errors.New("MazeSrv: " + message)
}

func GetMaze(req *MazeRequest, rep *MazeResponse) error {
	if req == nil {
		return mkErr("invalid request (empty)")
	}

	in := MazeInputs{
		width:      int(req.Width),
		height:     int(req.Height),
		tickSpeed:  int(req.TickSpeed),
		repeats:    int(req.Repeats),
		density:    int(req.Density),
		solveAlg:   req.SolveAlg,
		genAlg:     req.GenerateAlg,
		startIndex: int(req.StartIndex),
	}

	buf := new(bytes.Buffer)
	err := makeMaze(&in, buf)
	if err != nil {
		return err
	}

	rep.Webpage = buf.String()
	return nil
}
