// Filename: cmd/api/main.go

package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"time"
	"strings"

	"github.com/amari03/qod/internal/data"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type configuration struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	cors struct {
        trustedOrigins []string
    }
	limiter struct {
        rps float64                      // requests per second
        burst int                        // initial requests possible
        enabled bool                     // enable or disable rate limiter
    }
}

type application struct {
	config       configuration
	logger       *slog.Logger
	commentModel data.CommentModel
}

func main() {
	// Load config from flags
	cfg := loadConfig()

	// Create logger
	logger := setupLogger(cfg)

	// Open DB here (so it's in scope for app + stays open)
	db, err := openDB(cfg.db.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	// Wire dependencies (this is where your lines go)
	app := &application{
		config:       cfg,
		logger:       logger,
		commentModel: data.CommentModel{DB: db},
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
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")
	// read in the dsn
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://comments:fishsticks@localhost/comments", "PostgreSQL DSN")

	// We will build a custom command-line flag.  This flag will allow us to access 
	// space-separated origins. We will then put those origins in our slice. Again not // something we can do with the flag functions that we have seen so far. 
	// strings.Fields() splits string (origins) on spaces
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)",
	func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
	 	return nil
	})
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate Limiter maximum requests per second")

    flag.IntVar(&cfg.limiter.burst, "limiter-burst", 5, "Rate Limiter maximum burst")

    flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()

	return cfg
}

// setupLogger configures the application logger based on environment
func setupLogger(cfg configuration) *slog.Logger {
	//var logger *slog.Logger

	if cfg.env == "production" {
        return slog.New(slog.NewJSONHandler(os.Stdout, nil))
    }
    return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func openDB(dsn string) (*sql.DB, error) {
	// open a connection pool
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// set a context to ensure DB operations don't take too long
	ctx, cancel := context.WithTimeout(context.Background(),
		5*time.Second)
	defer cancel()
	// let's test if the connection pool was created
	// we trying pinging it with a 5-second timeout
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// return the connection pool (sql.DB)
	return db, nil

}
