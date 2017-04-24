package segment

import "context"

type chord struct {
}

func (c *chord) Alive(ctx context.Context, empty *chord.Empty) (*chord.Alive, error) {

}

func (c *chord) FindPredecessor(ctx context.Context, empty *chord.Empty) (*chord.Alive, error) {

}

func (c *chord) Init(ctx context.Context, empty *chord.Empty) (*chord.Alive, error) {

}

func (c *chord) Notify(ctx context.Context, empty *chord.Empty) (*chord.Alive, error) {

}

func (c *chord) Shutdown(ctx context.Context, empty *chord.Empty) (*chord.Alive, error) {

}
