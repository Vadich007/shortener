package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/Vadich007/shortener/internal/config"
	"github.com/Vadich007/shortener/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresLinkRepository struct {
	db *sql.DB
}

func NewPostrgesLinkRepository(conf config.Config) (*PostgresLinkRepository, error) {
	db, err := sql.Open("pgx", conf.DatabaseDsn)

	return &PostgresLinkRepository{
		db: db,
	}, err
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

func (r *PostgresLinkRepository) AddLinksBatch(request *model.BatchRequest, m map[string]string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	for _, record := range request.Records {
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
