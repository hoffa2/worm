package segment

import "context"

type chord struct {
}

func (c *chord) Alive(ctx context.Context, empty *chord.Empty) (*chord.Alive, error) {
	return nil, nil
}

func (c *chord) FindPredecessor(ctx context.Context, empty *chord.Empty) (*chord.Alive, error) {
	return nil, nil
}

func (c *chord) Init(ctx context.Context, empty *chord.Empty) (*chord.Alive, error) {
	return nil, nil
}

func (c *chord) Notify(ctx context.Context, empty *chord.Empty) (*chord.Alive, error) {
	return nil, nil

}

func (c *chord) Shutdown(ctx context.Context, empty *chord.Empty) (*chord.Alive, error) {
	return nil, nil

}
