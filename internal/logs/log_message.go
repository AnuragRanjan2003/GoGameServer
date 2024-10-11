package logs

import (
	"time"
)

type LogMessage struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
	Type    string    `json:"type"`
}

func NewLogMessage(message string, err error) *LogMessage {
	if err != nil {
		return &LogMessage{
			Time:    time.Now().UTC(),
			Message: err.Error(),
			Type:    "Error",
		}
	}

	return &LogMessage{
		Time:    time.Now().UTC(),
		Message: message,
		Type:    "Debug",
	}
}

func (m LogMessage) Bytes() []byte {
	msg := m.Time.String() + " : " + m.Type + " : " + m.Message
	return []byte(msg)
}
