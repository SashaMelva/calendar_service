package sqlstorage

import (
	"context"
	"database/sql"

	config "github.com/SashaMelva/calendar_service/internal/config"
	_ "github.com/jackc/pgx/stdlib"
	"go.uber.org/zap"
)

type Storage struct {
	StorageDb *sql.DB
}

func New(confDB *config.ConfigDB, log *zap.SugaredLogger) *Storage {
	dsn := "user=" + confDB.User + " dbname=" + confDB.NameDB + " sslmode=verify-full password=" + confDB.Password
	storage, err := sql.Open("pgx", dsn)

	if err != nil {
		log.Fatal("Cannot open pgx driver: %w", err)
	}

	return &Storage{storage}
}

func (s *Storage) Connect(ctx context.Context) error {
	err := s.StorageDb.PingContext(ctx)
	return err
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}
