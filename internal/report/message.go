package report

import (
	"encoding/json"
	"fmt"
)

type Severity string

const (
	Info    Severity = "info"
	Warning Severity = "warn"
	Error   Severity = "error"
)

type MessageType string

const (
	GlobalMessage     MessageType = "global"
	IngestorMessage   MessageType = "ingestor"
	PrecomputeMessage MessageType = "precompute"
	ComputeMessage    MessageType = "compute"
)

type Message struct {
	Type     MessageType
	Severity Severity
	Text     string
}

func (m Message) String() string {
	return fmt.Sprintf("%s - %s", m.Severity, m.Text)
}

func (m Message) MarshalJSON() ([]byte, error) {
	out, err := json.Marshal(m.String())
	if err != nil {
		return nil, err
	}
	return out, nil
}
