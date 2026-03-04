package repository

import (
	"context"
	"database/sql"
	"errors"

	"dependency-dashboard/internal/domain"
	"dependency-dashboard/internal/model"

	"github.com/rs/zerolog/log"
)

var ErrNotFound = domain.ErrNotFound

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) InitSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS dependencies (
		name TEXT PRIMARY KEY,
		version TEXT NOT NULL,
		openssf_score REAL NOT NULL CHECK(openssf_score >= 0),
		last_updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_score ON dependencies(openssf_score);
	`
	_, err := r.db.Exec(schema)
	return err
}

func (r *Repository) Upsert(ctx context.Context, d *model.Dependency) error {
	query := `
	INSERT INTO dependencies (name, version, openssf_score, last_updated)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(name) DO UPDATE SET
	version=excluded.version,
	openssf_score=excluded.openssf_score,
	last_updated=excluded.last_updated;
	`

	_, err := r.db.ExecContext(ctx, query,
		d.Name,
		d.Version,
		d.OpenSSFScore,
		d.LastUpdated,
	)

	return err
}

func (r *Repository) List(ctx context.Context, name string, minScore float64) ([]model.Dependency, error) {
	query, args := buildQuery(name, minScore)
	log.Debug().Str("query", query).Interface("args", args).Send()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deps []model.Dependency

	for rows.Next() {
		var d model.Dependency
		if err := rows.Scan(
			&d.Name,
			&d.Version,
			&d.OpenSSFScore,
			&d.LastUpdated,
		); err != nil {
			return nil, err
		}
		deps = append(deps, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deps, nil
}

func (r *Repository) Update(ctx context.Context, d *model.Dependency) error {
	query := `
		UPDATE dependencies
		SET version = ?, openssf_score = ?, last_updated = ?
		WHERE name = ?
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		d.Version,
		d.OpenSSFScore,
		d.LastUpdated,
		d.Name,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, name string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM dependencies WHERE name = ?", name)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (*model.Dependency, error) {
	query := `
		SELECT name, version, openssf_score, last_updated
		FROM dependencies
		WHERE name = ?
	`

	var d model.Dependency
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&d.Name,
		&d.Version,
		&d.OpenSSFScore,
		&d.LastUpdated,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func buildQuery(name string, minScore float64) (string, []any) {
	query := `
		SELECT name, version, openssf_score, last_updated
		FROM dependencies
		WHERE 1=1
	`
	args := []any{}

	if name != "" {
		query += " AND name LIKE ?"
		// Could be use "%"+name+"%" to search for name substring
		args = append(args, name+"%")
	}

	if minScore > 0 {
		query += " AND openssf_score >= ?"
		args = append(args, minScore)
	}

	query += " ORDER BY openssf_score DESC"

	return query, args
}
