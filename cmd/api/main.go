// Filename: cmd/api/main.go
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// A global variable to hold the application
// version number
const version = "1.0.0"

// Setup a struct to hold server configuration
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

// Setup dependency injection
type application struct {
	config config
	logger *log.Logger
}

// setup the main() function
func main() {
	var cfg config
	// Get the arguments for the user for
	// the server configuration
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("ORANGE_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()
	// Create a logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// setup the database connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	// close the database connection pool
	defer db.Close()
	// Create an object of type application
	app := &application{
		config: cfg,
		logger: logger,
	}
	// Create our server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,      // inactive connections
		ReadTimeout:  10 * time.Second, // time to read request body or header
		WriteTimeout: 10 * time.Second,
	}
	// Start our server
	logger.Printf("starting %s server on port %d", cfg.env, cfg.port)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// Setup up a database connection
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	// create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// ping the database
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
