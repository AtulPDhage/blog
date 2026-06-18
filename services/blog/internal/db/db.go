package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"blog/internal/config"
	"blog/internal/logger"
)

var Pool *pgxpool.Pool

// ConnectDB establishes a connection pool to the database using the configured values
func ConnectDB(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgxCfg, err := pgxpool.ParseConfig(cfg.DBURL)
	if err != nil {
		return fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Connection pooling configuration driven by environment variables with fallbacks
	pgxCfg.MaxConns = cfg.DBMaxConns
	pgxCfg.MinConns = cfg.DBMinConns
	pgxCfg.MaxConnLifetime = cfg.DBMaxConnLifetime
	pgxCfg.MaxConnIdleTime = cfg.DBMaxConnIdleTime

	p, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return fmt.Errorf("unable to create database pool: %w", err)
	}

	// Ping database to ensure connectivity
	if err := p.Ping(ctx); err != nil {
		p.Close()
		return fmt.Errorf("unable to ping database: %w", err)
	}

	Pool = p
	logger.Logger.Info("Database connection established successfully")
	return nil
}

// InitTables initializes database tables if they do not exist
func InitTables(ctx context.Context) error {
	queries := []string{
		CreateBlogsTableQuery,
		CreateCommentsTableQuery,
		CreateSavedBlogsTableQuery,
		CreateLikedBlogsTableQuery,
		AlterBlogsAddViewsQuery,
	}

	for _, query := range queries {
		if _, err := Pool.Exec(ctx, query); err != nil {
			return fmt.Errorf("error creating database table or column: %w", err)
		}
	}

	logger.Logger.Info("Database tables initialized successfully")
	return nil
}
