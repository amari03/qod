// Filename: cmd/api/main.go

package main

import (
	"flag"
    "log/slog"
    "os"
)

const version = "1.0.0"

type configuration struct {
	port int
	env  string
}

type application struct {
	config configuration
	logger *slog.Logger
}

func main() {
	// Load config from flags
	cfg := loadConfig()

	// Create logger (text to stdout)
	logger := setupLogger(cfg.env)

	// Wire dependencies
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Start server
	if err := app.serve(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// loadConfig reads configuration from command line flags
func loadConfig() configuration {
    var cfg configuration
	
flag.IntVar(&cfg.port, "port", 4000, "API server port")
flag.StringVar(&cfg.env,"env","development","Environment(development|staging|production)")
flag.Parse()
	
return cfg
}

// setupLogger configures the application logger based on environment
func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
		
		
	return logger
}

