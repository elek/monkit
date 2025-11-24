package motel

import (
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestProvider(t *testing.T) {
	tp := &TraceProvider{}

	otel.SetTracerProvider(tp)
	otel.GetTracerProvider()
	trace.NewTracerProvider()

}
