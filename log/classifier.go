package log

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Classes = map[string]Level

// log by string classes in addition to levels
type Classifier struct {
	logger         zap.Logger
	allowedClasses Classes
	funcCounts     map[string]uint64
	funcCosts      map[string][]time.Duration
	l              zap.Logger
	mu             sync.Mutex
}

func (s *Classifier) Try(className string) (*zap.Logger, *Costs) {
	funcName := RuntimeFunctionName(1 + 1)
	count, ok := s.checkClassAndCount(className, funcName)
	costs := &Costs{
		Count:     count,
		startedAt: time.Now(),
	}
	if !ok {
		return nil, costs
	}
	costs.Names = make([]string, 0, 8)
	costs.Costs = make([]time.Duration, 0, 8)
	return s.l.With([]zapcore.Field{
		zap.String("class", className),
		zap.String("func", funcName),
	}...), costs
}

func (s *Classifier) checkClassAndCount(className, funcName string) (count uint64, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok = s.allowedClasses[className]
	if !ok { // no logging
		return
	}
	count = s.funcCounts[funcName] + 1
	s.funcCounts[funcName] = count
	return
}

// func NewClassifier(allowedClasses Classes) *Classifier {
// 	return &Classifier{
// 		allowedClasses: allowedClasses,
// 		funcCounts:     make(map[string]uint64),
// 	}
// }

// func (s *Classifier) Enable(className string) *Classifier {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	s.allowedClasses[className] = true
// 	return s
// }

// func (s *Classifier) Disable(className string) *Classifier {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	delete(s.allowedClasses, className)
// 	return s
// }
