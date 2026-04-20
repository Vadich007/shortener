package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/Vadich007/shortener/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresLinkRepository struct {
	db *sql.DB
}

func NewPostrgesLinkRepository(conf config.Config) *PostgresLinkRepository {
	db, err := sql.Open("pgx", conf.DatabaseDsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	return &PostgresLinkRepository{
		db: db,
	}
}

func (r *PostgresLinkRepository) GetLink(shortedLink string) (string, error) {
	row := r.db.QueryRow("SELECT original_url FROM links WHERE shorted_url = ?", shortedLink)
	// готовим переменную для чтения результата
	var originalLink string
	err := row.Scan(&originalLink)
	return originalLink, err
}

func (r *PostgresLinkRepository) AddLink(shortedLink string, originalLink string) error {
	_, err := r.db.Exec(
		"INSERT INTO links (shorted_url, original_url) VALUES ($1, $2)",
		shortedLink,
		originalLink,
	)
	return err
}

func (r *PostgresLinkRepository) PingDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := r.db.PingContext(ctx); err != nil {
		return sql.ErrConnDone
	}
	return nil
}
