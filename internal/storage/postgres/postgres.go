package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/behummble/Questions-answers/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	log *slog.Logger
	conn *gorm.DB
}

func NewStorage(ctx context.Context, log *slog.Logger, cfg config.StorageConfig) *Storage {
	dsn := parseConnectStr(cfg)
	conn, err := gorm.Open(
		postgres.Open(dsn), 
		&gorm.Config{},
	)
	if err != nil {
		panic(err)
	}
	return &Storage{
		log: log,
		conn: conn,
	}
}

func(storage *Storage) Shutdown(ctx context.Context) {
	db, err := storage.conn.DB()
	if err != nil {
		return
	}
	db.Close()
}

func parseConnectStr(cfg config.StorageConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		cfg.Host,
		cfg.Username,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.TimeZone,
	)
}