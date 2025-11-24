package hotel

import (
	"context"
	"sync"

	"github.com/spacemonkeygo/monkit/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func Package() *Scope {
	tracer := otel.Tracer(monkit.CallerPackage(1))
	meter := otel.Meter(monkit.CallerPackage(1))
	return &Scope{
		tracer:  tracer,
		meter:   meter,
		sources: make(map[string]any),
	}
}

type Scope struct {
	tracer  trace.Tracer
	meter   metric.Meter
	mtx     sync.RWMutex
	sources map[string]any
}

func (s *Scope) Task(tags ...monkit.SeriesTag) monkit.Task {
	return func(ctx *context.Context, args ...interface{}) func(*error) {
		tctx := *ctx
		nctx, task := s.tracer.Start(tctx, monkit.CallerFunc(2))
		ctx = &nctx
		return func(*error) {
			task.End()
		}
	}
}

func (s *Scope) IntVal(name string) *IntVal {
	iv := newSource[*IntVal](s, name, func() *IntVal {
		counter, err := s.meter.Int64Counter(name)
		// TODO: is this safe?
		if err != nil {
			panic(err)
		}
		return &IntVal{
			counter: counter,
		}
	})
	return iv

}

type IntVal struct {
	counter metric.Int64Counter
}

func (i *IntVal) Observe(value int64) {
	i.counter.Add(context.Background(), value)
}

func newSource[T any](s *Scope, name string, constructor func() T) (rv T) {
	s.mtx.RLock()
	source, exists := s.sources[name]
	s.mtx.RUnlock()

	if exists {
		// TODO: is this safe?
		return source.(T)
	}

	s.mtx.Lock()
	if source, exists := s.sources[name]; exists {
		s.mtx.Unlock()
		// TODO: is this safe?
		return source.(T)
	}

	ss := constructor()
	s.sources[name] = ss
	s.mtx.Unlock()

	return ss
}
