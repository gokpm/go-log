package log

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	olog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

var ok bool
var provider *log.LoggerProvider

type Config struct {
	Ok          bool
	Name        string
	Environment string
	URL         string
}

func Setup(ctx context.Context, config *Config) (olog.Logger, error) {
	ok = config.Ok
	if !ok {
		return nil, nil
	}
	httpOpts := []otlploghttp.Option{
		otlploghttp.WithEndpointURL(config.URL),
		otlploghttp.WithCompression(otlploghttp.GzipCompression),
	}
	exporter, err := otlploghttp.New(ctx, httpOpts...)
	if err != nil {
		return nil, err
	}
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	base := resource.Default()
	newResource := resource.NewWithAttributes(
		base.SchemaURL(),
		semconv.ServiceName(config.Name),
		semconv.DeploymentEnvironmentName(config.Environment),
		semconv.HostName(hostname),
	)
	mergedResource, err := resource.Merge(base, newResource)
	if err != nil {
		return nil, err
	}
	processor := log.NewBatchProcessor(exporter)
	providerOpts := []log.LoggerProviderOption{
		log.WithResource(mergedResource),
		log.WithProcessor(processor),
	}
	provider = log.NewLoggerProvider(providerOpts...)
	return provider.Logger(config.Name), nil
}

func Shutdown(timeout time.Duration) error {
	if !ok {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()
	err := provider.ForceFlush(ctx)
	if err != nil {
		return err
	}
	err = provider.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
