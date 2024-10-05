package types

type Delta interface {
	GetProducer() string
	GetTimeStamp() uint
	GetDelta() interface{}
	GetType() uint8
}
