package main

import (
	"context"
	"fmt"
	"github.com/biryanim/SongLibrary/config"
	"github.com/biryanim/SongLibrary/internal/adapters/db"
	"github.com/biryanim/SongLibrary/internal/adapters/http"
	"github.com/biryanim/SongLibrary/internal/usecases"
	"github.com/biryanim/SongLibrary/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func initDB(cfg *config.Config) (*pgx.Conn, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	m, err := migrate.New("file://migrations", connStr)
	if err != nil {
		return nil, err
	}

	if err = m.Up(); err != nil && err.Error() != "no change" {
		return nil, err
	}
	return conn, nil
}

// @title Song Library  API
// @version 1.0

// @host localhost:8080
// @BasePath /
func main() {
	if err := logger.Initialize(); err != nil {
		log.Fatal(err)
	}
	cfg := config.New()
	conn, err := initDB(cfg)
	if err != nil {
		logger.Log.Fatal("cannot connect to database", zap.Error(err))
	}
	defer conn.Close(context.Background())
	storage := db.New(conn)
	controller := usecases.New(storage)
	server := http.New(controller)
	http.StartServer(8080, server)
}
