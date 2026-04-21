package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/model"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresLinkRepository struct {
	db *sql.DB
}

func NewPostrgesLinkRepository(conf config.Config) (*PostgresLinkRepository, error) {
	db, err := sql.Open("pgx", conf.DatabaseDsn)

	// if err := runMigrations(db); err != nil {
	// 	return nil, err
	// }

	return &PostgresLinkRepository{
		db: db,
	}, err
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create database driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	return nil
}

func (r *PostgresLinkRepository) GetLink(shortedLink string) (string, error) {
	row := r.db.QueryRow("SELECT original_url FROM links WHERE shorted_url = $1", shortedLink)
	var originalLink string
	err := row.Scan(&originalLink)
	return originalLink, err
}

func (r *PostgresLinkRepository) AddLink(shortedLink string, originalLink string) error {
	if _, err := r.GetLink(shortedLink); err == nil {
		return model.NewLinkAlreadyExistError(shortedLink)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"INSERT INTO links (shorted_url, original_url) VALUES ($1, $2)",
		shortedLink,
		originalLink,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *PostgresLinkRepository) PingDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := r.db.PingContext(ctx); err != nil {
		return sql.ErrConnDone
	}
	return nil
}

func (r *PostgresLinkRepository) AddLinksBatch(request []model.BatchRecordRequest, m map[string]string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	for _, record := range request {
		_, err := tx.Exec("INSERT INTO links (shorted_url, original_url) VALUES ($1, $2)",
			m[record.CorrelationID],
			record.OriginalURL,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
