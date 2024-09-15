package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// Configure the OTLP gRPC exporter to send trace data to the OpenTelemetry Collector
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint("collector-gateway:4317"),
		otlptracegrpc.WithInsecure(),
	)
	// Create a new OTLP trace exporter
	otlpExporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}
	// Set up resource attributes (such as service name)
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("hello-Service"),
		),
	)
	if err != nil {
		return nil, err
	}
	// Create a TracerProvider with the OTLP exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(otlpExporter),
		sdktrace.WithResource(res),
	)
	// Set the TracerProvider globally
	otel.SetTracerProvider(tp)
	return tp, nil

}
func helloHandler(w http.ResponseWriter, req *http.Request) {
	// Use the global tracer to start a new span
	// ctx := req.Context()
	// tracer := otel.Tracer("hello-handler")
	// _, span := tracer.Start(ctx, "helloHandler")
	// defer span.End()
	// log.Println("Span created for helloHandler")

	w.Write([]byte("Hello, World"))
}

func main() {
	tp, err := initTracer()
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		// Ensure all spans are flushed before shutdown
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}()
	wrappedHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "/hello")
	http.Handle("/hello", wrappedHandler)
	fmt.Printf("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
