package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/damirm/links-warehouse/internal/model"

	"github.com/damirm/links-warehouse/internal/postgres/querygen"
	"github.com/damirm/links-warehouse/internal/storage"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "github.com/tabbed/pqtype"
)

func Connect(config *Config) (*sql.DB, error) {
	return sql.Open("postgres", config.ConnString())
}

func InitStorage(ctx context.Context, db *sql.DB, config *Config) (*Storage, error) {
	migrator := &Migrator{db, config}

	err := migrator.Migrate(ctx)
	if err != nil {
		return nil, err
	}

	return NewStorage(db, config), nil
}

type Migrator struct {
	db     *sql.DB
	config *Config
}

func (m *Migrator) Migrate(ctx context.Context) error {
	driver, err := postgres.WithInstance(m.db, &postgres.Config{
		SchemaName:   m.config.Schema,
		DatabaseName: m.config.Database,
	})
	if err != nil {
		return err
	}

	mg, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", m.config.MigrationsPath),
		m.config.Database,
		driver,
	)
	if err != nil {
		return err
	}

	err = mg.Up()
	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}

type Storage struct {
	queries *querygen.Queries
	config  *Config
	db      *sql.DB
}

func NewStorage(db *sql.DB, config *Config) *Storage {
	return &Storage{
		db:      db,
		queries: querygen.New(db),
		config:  config,
	}
}

func (s *Storage) withTx(tx *sql.Tx) *Storage {
	return &Storage{
		db:      s.db,
		queries: s.queries.WithTx(tx),
	}
}

func (s *Storage) Transaction(ctx context.Context, fn func(context.Context, storage.Storage) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			// TODO: Log error.
		}
	}(tx)
	err = fn(ctx, s.withTx(tx))
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Storage) SaveLink(ctx context.Context, link *model.Link) error {
	params := insertLinkParams(link)
	return s.queries.InsertLink(ctx, params)
}

func (s *Storage) EnqueueURL(ctx context.Context, u *url.URL) error {
	return s.queries.EnqueueUrl(ctx, u.String())
}

func (s *Storage) DequeueURL(ctx context.Context) (*url.URL, error) {
	res, err := s.queries.DequeueUrl(ctx)
	if err != nil {
		return nil, err
	}
	return url.Parse(res.Url)
}

func (s *Storage) DeleteProcessedURL(ctx context.Context, u *url.URL) error {
	return s.queries.DeleteQueuedUrl(ctx, u.String())
}
