package peer

type Message[T any] struct {
	Type   string
	Status int32
	Body   T
}
