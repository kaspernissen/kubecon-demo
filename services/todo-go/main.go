package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

var (
	db *sql.DB
	// Acquire tracer of hybrid instrumentation
	tracer = otel.Tracer("todo")
)

func main() {
	// Read DB configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "mysecretpassword")
	dbName := getEnv("DB_NAME", "todo")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	// Construct the connection string
	psqlconn := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=" + dbSSLMode

	var err error
	db, err = sql.Open("postgres", psqlconn)
	if err != nil {
		Logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	setupDB()
	setupMeter()

	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/todos/all", AllTodos).Methods("GET")
	r.HandleFunc("/todos/{id}", GetTodo).Methods("GET")
	r.HandleFunc("/todos", CreateTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", DeleteTodo).Methods("DELETE")
	r.HandleFunc("/todos/{id}", UpdateTodo).Methods("PUT")

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		Logger.Fatal().Err(err).Msg("Server failed to start")
	}
	Logger.Info().Msg("Server running")
}

// setupDB ensures the database table exists
func setupDB() {
	query := `
		CREATE TABLE IF NOT EXISTS todo (
			id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
			name TEXT NOT NULL
		);
	`
	if _, err := db.Exec(query); err != nil {
		Logger.Fatal().Err(err).Msg("Failed to initialize database")
	}
}

// getEnv retrieves environment variables with a default fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func setupMeter() {
	// Initialize the Prometheus exporter
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatalf("failed to initialize Prometheus exporter: %v", err)
	}

	// Create a MeterProvider with the exporter
	meterProvider := metric.NewMeterProvider(metric.WithReader(exporter))
	defer func() {
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shut down MeterProvider: %v", err)
		}
	}()

	// Start collecting runtime metrics
	if err := runtime.Start(runtime.WithMeterProvider(meterProvider)); err != nil {
		log.Fatalf("failed to start runtime instrumentation: %v", err)
	}
}
