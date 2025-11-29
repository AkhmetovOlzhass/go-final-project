package otelinit

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
)

func Init(serviceName string) func() {
	ctx := context.Background()

	conn, err := grpc.DialContext(
		ctx,
		"localhost:14317",
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("[OTEL] Failed to connect to collector: %v", err)
	}

	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithGRPCConn(conn),
	)
	if err != nil {
		log.Fatalf("[OTEL] Failed to create exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)

	log.Println("[OTEL] Tracer initialized")

	return func() {
		_ = tp.Shutdown(ctx)
		_ = conn.Close()
	}
}
