package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type serveFunc func(handler http.Handler) error

func main() {
	runService(func(handler http.Handler) error {
		httpAddr := fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))
		return http.ListenAndServe(httpAddr, handler)
	})
}

// runService runs a service.
// The main purpose of this function is to have an ability to run
// functional testing by passing the serving function.
func runService(serve serveFunc) {
	logger := logrus.New()

	db, err := dbConnect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	service := createService(db, logger)
	handler := drawRoutes(service)

	if err := serve(handler); err != nil {
		logger.WithError(err).Error("error during serving")
	}
}

// createService allocates the service instance.
func createService(db *pgx.ConnPool, logger *logrus.Logger) ConfigServicer {
	pgRepo := NewPostgresConfigRepo(db)
	repo := NewConfigRepoLogger(pgRepo, logger)

	return NewConfigService(repo)
}

// drawRoutes draws service routes.
func drawRoutes(service ConfigServicer) http.Handler {
	r := mux.NewRouter()
	service.SetupRoutes(r.PathPrefix("/configs").Subrouter())
	return handlers.LoggingHandler(os.Stdout, r)
}

// dbConnect creates pgx connection pool.
func dbConnect() (*pgx.ConnPool, error) {
	port, err := strconv.ParseUint(os.Getenv("PG_PORT"), 10, 16)
	if err != nil {
		return nil, fmt.Errorf("Can't cast %q to uint16: %v", os.Getenv("PG_PORT"), err)
	}

	return pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     os.Getenv("PG_HOST"),
			Port:     uint16(port),
			User:     os.Getenv("PG_USER"),
			Password: os.Getenv("PG_PASSWORD"),
			Database: os.Getenv("PG_DATABASE"),
		},
		AcquireTimeout: time.Second,
		MaxConnections: 4,
	})
}
