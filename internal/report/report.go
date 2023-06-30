package report

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Report struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	mutex *sync.Mutex

	Logs        map[MessageType][]Message `json:"logs"`
	Status      jobStatus                 `json:"status"`
	Performance PerformanceStats          `json:"performance"`
}

// NewReport creates and initializes a new Report.
func NewReport() *Report {
	return &Report{
		mutex:  &sync.Mutex{},
		Status: Pending,
		Logs:   make(map[MessageType][]Message),
	}
}

func logMessage(msg Message) {
	level, err := zerolog.ParseLevel(string(msg.Severity))
	if err != nil {
		log.Error().Err(err).Msgf("severity %s is not valid", msg.Severity)
		level = zerolog.ErrorLevel
	}
	log.WithLevel(level).Msg(msg.Text)
}

// Watch watches for log messages coming from the different workers.
// It ends when the channel is closed.
// This function is concurrent-safe.
func (r *Report) Watch(messageChan <-chan Message) {
	log.Info().Msg("Starting report dispatcher")
	r.StartTime = time.Now()
	for msg := range messageChan {
		logMessage(msg)
		r.mutex.Lock()
		r.Logs[msg.Type] = append(r.Logs[msg.Type], msg)
		r.mutex.Unlock()
	}
	r.EndTime = time.Now()
	log.Info().Msg("Stopping report dispatcher")
}

func (r *Report) ToJSON() ([]byte, error) {
	if r == nil {
		return []byte("{}"), nil
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return json.MarshalIndent(r, "", "  ")
}
