package maze

import (
	"context"
)

type BFSSender struct {
}

type ChildParentPair struct {
	Parent   int
	Child    int
	ThreadID int
}

func RunBFSSender(public bool) error {
	_ = &BFSSender{}
	return nil

}

func (s *BFSSender) CheckNode(ctx context.Context, req *ChildParentPair, res *Done) error {
	return nil
}
