package main

import (
	"context"
	"fmt"
	"os"
	"log"
	"net/http"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)
// Redis client
var rdb *redis.Client

func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// Configure the OTLP gRPC exporter to send trace data to the OpenTelemetry Collector
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint("collector-agent:4317"),
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
func initRedis(){
	redisHost := os.Getenv("REDIS_HOST")
    redisPort := os.Getenv("REDIS_PORT")
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB: 0,
	})
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		log.Fatalf("failed to instrument Redis with OpenTelemetry: %v", err)
	}
}
func helloHandler(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	//Store data in Redis
	err:= rdb.Set(ctx, "greeting", "Hello Redis!",0).Err()
	if err != nil{
		log.Fatalf("failed to set key in Redis: %v", err)
	}

	//Retrieve data from Redis
	val, err := rdb.Get(ctx, "greeting").Result()
	if err != nil{
		log.Fatalf("failed to get key from Redis: %v", err)
	}

	w.Write([]byte(val))
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
	initRedis()
	wrappedHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "/hello")
	http.Handle("/hello", wrappedHandler)
	fmt.Printf("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
