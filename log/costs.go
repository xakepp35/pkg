package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Costs представляет информацию о затратах времени выполнения операций.
// @Description Содержит количество операций, названия операций, время выполнения каждой операции и дату начала.
type Costs struct {
	Count uint64          `json:"count" example:"10"`                   // Количество выполненных операций.
	Names []string        `json:"names" example:"dedupe,insert"` // Список наименований операций.
	Costs []time.Duration `json:"costs" example:"1000000,5000000"`     // Список временных затрат для выполнения операций.

	startedAt time.Time `json:"-"`
}

func (s *Costs) At() time.Time {
	return s.startedAt
}

func (s *Costs) Observe(stepName string) {
	observedAt := time.Now()
	cost := observedAt.Sub(s.startedAt)
	s.Names = append(s.Names, stepName)
	s.Costs = append(s.Costs, cost)
	s.startedAt = observedAt
}

func (s *Costs) Outcome() zapcore.Field {
	return zap.Object("costs", s)
}

var _ zapcore.ObjectMarshaler = (*Costs)(nil)

func (s *Costs) MarshalLogObject(e zapcore.ObjectEncoder) error {
	var total time.Duration
	for _, c := range s.Costs {
		total += c
	}
	e.AddUint64("count", s.Count)
	e.AddFloat64("total", Cost(total))
	for i, cost := range s.Costs {
		e.AddFloat64(s.Names[i], Cost(cost))
	}
	return nil
}

func Cost(x time.Duration) float64 {
	return x.Seconds() * 1000
}
