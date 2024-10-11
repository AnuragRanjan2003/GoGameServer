package data

type GameDelta struct {
	Producer  string      `json:"producer"`
	TimeStamp uint        `json:"time"`
	Delta     interface{} `json:"delta"`
}

func (g GameDelta) GetProducer() string {
	return g.Producer
}

func (g GameDelta) GetTimeStamp() uint {
	return g.TimeStamp
}

func (g GameDelta) GetDelta() interface{} {
	return g.Delta
}

func (g GameDelta) GetType() uint8 {
	return 0
}