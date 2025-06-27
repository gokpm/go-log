package log

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

type Config struct {
	Name        string
	Environment string
	URL         string
}

func Setup(ctx context.Context, config *Config) error {
	httpOpts := []otlploghttp.Option{
		otlploghttp.WithEndpointURL(config.URL),
		otlploghttp.WithCompression(otlploghttp.GzipCompression),
	}
	exporter, err := otlploghttp.New(ctx, httpOpts...)
	if err != nil {
		return err
	}
	hostname, err := os.Hostname()
	if err != nil {
		return err
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
		return err
	}
	processor := log.NewBatchProcessor(exporter)
	providerOpts := []log.LoggerProviderOption{
		log.WithResource(mergedResource),
		log.WithProcessor(processor),
	}
	_ = log.NewLoggerProvider(providerOpts...).Logger(config.Name)
	return nil
}
