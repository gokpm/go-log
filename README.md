# go-log

A simple Go package for OpenTelemetry logging with OTLP HTTP export.

## Installation

```bash
go get github.com/gokpm/go-log
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"
    
    "github.com/gokpm/go-log"
)

func setup() error {
    config := log.Config{
        Ok:          true,
        Name:        "my-service",
        Environment: "production",
        URL:         "http://localhost:4318/v1/logs",
    }
    
    ctx := context.Background()
    _, err := log.Setup(ctx, config)
    return err
}

func main() {
    if err := setup(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    
    defer log.Shutdown(5 * time.Second)
    
    // Use logger for logging...
}
```

## Configuration

- `Ok`: Enable/disable logging
- `Name`: Service name
- `Environment`: Deployment environment
- `URL`: OTLP HTTP endpoint URL (default: `http://localhost:4318/v1/logs`)

## Features

- OTLP HTTP export with gzip compression
- Automatic resource detection (hostname, service info)
- Graceful shutdown with timeout
- Batch processing for performance