package receiver

import "context"

var _ Receiver = (*LocalReceiver)(nil)

type LocalReceiver struct {
	*eventRouter
}

func NewLocalReceiver(opts ...ReceiverOption) *LocalReceiver {
	router := newEventRouter(opts...)

	return &LocalReceiver{
		eventRouter: router,
	}
}

func (r *LocalReceiver) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}
