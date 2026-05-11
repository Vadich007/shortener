package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return &PostgresLinkRepository{db: db}, err
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
	row := r.db.QueryRow(
		"SELECT original_url, is_deleted FROM links WHERE shorted_url = $1",
		shortedLink,
	)
	var originalLink string
	var isDeleted bool
	if err := row.Scan(&originalLink, &isDeleted); err != nil {
		return "", err
	}
	if isDeleted {
		return "", model.NewLinkDeletedError(shortedLink)
	}
	return originalLink, nil
}

func (r *PostgresLinkRepository) AddLink(shortedLink, originalLink string, userID int) error {
	if _, err := r.GetLink(shortedLink); err == nil {
		return model.NewLinkAlreadyExistError(shortedLink)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"INSERT INTO links (shorted_url, original_url, user_id) VALUES ($1, $2, $3)",
		shortedLink, originalLink, userID,
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

func (r *PostgresLinkRepository) AddLinksBatch(request []model.BatchRecordRequest, m map[string]string, userID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	for _, record := range request {
		_, err := tx.Exec(
			"INSERT INTO links (shorted_url, original_url, user_id) VALUES ($1, $2, $3)",
			m[record.CorrelationID], record.OriginalURL, userID,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *PostgresLinkRepository) GetUserUrls(userID int) ([]model.UserURLResponse, error) {
	rows, err := r.db.Query(
		"SELECT shorted_url, original_url FROM links WHERE user_id = $1 AND is_deleted = false",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.UserURLResponse
	for rows.Next() {
		var rec model.UserURLResponse
		if err := rows.Scan(&rec.ShortURL, &rec.OriginalURL); err != nil {
			return nil, err
		}
		result = append(result, rec)
	}
	return result, rows.Err()
}

func (r *PostgresLinkRepository) DeleteURLsBatch(userID int, shortURLs []string) error {
	if len(shortURLs) == 0 {
		return nil
	}

	args := make([]interface{}, 0, len(shortURLs)+1)
	args = append(args, userID)
	placeholders := make([]string, len(shortURLs))
	for i, u := range shortURLs {
		args = append(args, u)
		placeholders[i] = fmt.Sprintf("$%d", i+2)
	}

	query := fmt.Sprintf(
		"UPDATE links SET is_deleted = true WHERE user_id = $1 AND shorted_url IN (%s)",
		strings.Join(placeholders, ", "),
	)

	_, err := r.db.Exec(query, args...)
	return err
}
