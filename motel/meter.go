package motel

import (
	"context"

	"github.com/spacemonkeygo/monkit/v3"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

type MeterProvider struct {
	noop.MeterProvider
	registry *monkit.Registry
}

func NewMeterProvider(registry *monkit.Registry) *MeterProvider {
	return &MeterProvider{
		registry: registry,
	}
}

func (m *MeterProvider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	return &Meter{
		scope: m.registry.ScopeNamed(name),
	}
}

var _ metric.MeterProvider = &MeterProvider{}

type Meter struct {
	noop.Meter
	scope *monkit.Scope
}

type Int64Counter struct {
	noop.Int64Counter
	value *monkit.IntVal
}

func (i *Int64Counter) Add(ctx context.Context, incr int64, options ...metric.AddOption) {
	i.value.Observe(incr)
}

func (m *Meter) Int64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return &Int64Counter{
		value: m.scope.IntVal(name),
	}, nil
}

var _ metric.Meter = &Meter{}
