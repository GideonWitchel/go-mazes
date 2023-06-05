package maze

import (
	"context"
)

type BFSReceiver struct {
	val int
	g   *graph
}

type Parent struct {
	Index int
}

type Done struct {
	Val bool
}

func RunBFSReceiver(graph *graph) error {
	_ = &BFSReceiver{
		g:   graph,
		val: 3,
	}
	return nil
}

func (r *BFSReceiver) GetNeighbors(ctx context.Context, req *Parent, res *Done) error {
	currentNode := r.g.nodes[req.Index]
	for _, currentNeighbor := range currentNode.neighbors {
		// XXX DO RPC ChildParentPair{Parent: req.GetIndex(), Child: uint64(currentNeighbor.n.index), ThreadID: Id}
		// Check for termination after sending the value so the parents array knows where the solution is
		if currentNeighbor.n.val == r.val {
			// Terminate
			//XXX DO RPC ChildParentPair{Parent: int64(currentNeighbor.n.index), Child: -1, ThreadID: Id}
		}
	}
	res.Val = true
	return nil
}
