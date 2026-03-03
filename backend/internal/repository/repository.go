package repository

import (
	"context"
	"database/sql"

	"dependency-dashboard/internal/model"
)

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
	// TODO_TOM paging/metrics
	query := `
	SELECT name, version, openssf_score, last_updated
	FROM dependencies
	WHERE (? = '' OR name LIKE ?)
	AND (? = 0 OR openssf_score >= ?)
	ORDER BY openssf_score DESC;
	`

	rows, err := r.db.QueryContext(ctx, query,
		name,
		"%"+name+"%",
		minScore,
		minScore,
	)
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

	return deps, nil
}

func (r *Repository) Delete(ctx context.Context, name string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM dependencies WHERE name = ?", name)
	return err
}
