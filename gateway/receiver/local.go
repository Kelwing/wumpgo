package receiver

type LocalReceiver struct {
	*eventRouter
}

func NewLocalReceiver(opts ...ReceiverOption) *LocalReceiver {
	router := newEventRouter(opts...)

	return &LocalReceiver{
		eventRouter: router,
	}
}
