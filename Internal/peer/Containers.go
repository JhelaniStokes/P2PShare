package peer

type Message[T any] struct {
	Type string
	Body T
}
