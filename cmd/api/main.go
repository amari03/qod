// Filename: cmd/api/main.go

package main

import (
	"context"
    "database/sql"
	"flag"
	"fmt"
    "log/slog"
    "os"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type configuration struct {
	port int
	env  string
	db struct {
        dsn string
    }
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
// read in the dsn
    flag.StringVar(&settings.db.dsn, "db-dsn", "postgres://comments:fishsticks@localhost/comments","PostgreSQL DSN")

flag.Parse()
	
return cfg
}

// setupLogger configures the application logger based on environment
func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	// the call to openDB() sets up our connection pool
	db, err := openDB(settings)
	if err != nil {
    	logger.Error(err.Error())
    	os.Exit(1)
	}
	// release the database resources before exiting
	defer db.Close()

	logger.Info("database connection pool established")	
		
	return logger
}

func openDB(settings serverConfig) (*sql.DB, error) {
    // open a connection pool
    db, err := sql.Open("postgres", settings.db.dsn)
    if err != nil {
        return nil, err
    }
    
    // set a context to ensure DB operations don't take too long
    ctx, cancel := context.WithTimeout(context.Background(),
                                       5 * time.Second)
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
