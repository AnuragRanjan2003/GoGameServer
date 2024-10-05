package data

type NPDelta struct {
	Producer  string      `json:"producer"`
	TimeStamp uint        `json:"time_stamp"`
	Delta     interface{} `json:"delta"`
}

func (d NPDelta) GetProducer() string {
	return d.Producer
}

func (d NPDelta) GetTimeStamp() uint {
	return d.TimeStamp
}

func (d NPDelta) GetDelta() interface{} {
	return d.Delta
}
func (d NPDelta) GetType() uint8 {
	return 1
}
