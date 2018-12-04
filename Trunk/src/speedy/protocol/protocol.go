package protocol

type Protocol interface {
	IOLoop() error
}
