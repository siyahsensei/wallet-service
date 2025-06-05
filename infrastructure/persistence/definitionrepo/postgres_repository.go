package definitionrepo

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"siyahsensei/wallet-service/domain/definition"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, def *definition.Definition) error {
	query := `
		INSERT INTO definitions (
			id, name, abbreviation, suffix, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		def.ID,
		def.Name,
		def.Abbreviation,
		def.Suffix,
		def.CreatedAt,
		def.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*definition.Definition, error) {
	query := `
		SELECT id, name, abbreviation, suffix, created_at, updated_at
		FROM definitions
		WHERE id = $1
	`
	var def definition.Definition
	err := r.db.GetContext(ctx, &def, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("definition not found")
		}
		return nil, err
	}
	return &def, nil
}

func (r *PostgresRepository) GetByAbbreviation(ctx context.Context, abbreviation string) (*definition.Definition, error) {
	query := `
		SELECT id, name, abbreviation, suffix, created_at, updated_at
		FROM definitions
		WHERE abbreviation = $1
	`
	var def definition.Definition
	err := r.db.GetContext(ctx, &def, query, abbreviation)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("definition not found")
		}
		return nil, err
	}
	return &def, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context, limit, offset int) ([]*definition.Definition, error) {
	query := `
		SELECT id, name, abbreviation, suffix, created_at, updated_at
		FROM definitions
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`
	var definitions []*definition.Definition
	err := r.db.SelectContext(ctx, &definitions, query, limit, offset)
	if err != nil {
		return nil, err
	}
	return definitions, nil
}

func (r *PostgresRepository) Update(ctx context.Context, def *definition.Definition) error {
	def.UpdatedAt = time.Now()
	query := `
		UPDATE definitions
		SET name = $1, abbreviation = $2, suffix = $3, updated_at = $4
		WHERE id = $5
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		def.Name,
		def.Abbreviation,
		def.Suffix,
		def.UpdatedAt,
		def.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("definition not found")
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM definitions
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("definition not found")
	}
	return nil
}

func (r *PostgresRepository) Search(ctx context.Context, searchTerm string, limit, offset int) ([]*definition.Definition, error) {
	searchPattern := "%" + strings.ToLower(searchTerm) + "%"
	query := `
		SELECT id, name, abbreviation, suffix, created_at, updated_at
		FROM definitions
		WHERE LOWER(name) LIKE $1 OR LOWER(abbreviation) LIKE $1
		ORDER BY 
			CASE 
				WHEN LOWER(abbreviation) = LOWER($2) THEN 1
				WHEN LOWER(abbreviation) LIKE LOWER($2) || '%' THEN 2
				WHEN LOWER(name) = LOWER($2) THEN 3
				WHEN LOWER(name) LIKE LOWER($2) || '%' THEN 4
				ELSE 5
			END,
			name ASC
		LIMIT $3 OFFSET $4
	`
	var definitions []*definition.Definition
	err := r.db.SelectContext(ctx, &definitions, query, searchPattern, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	return definitions, nil
}
