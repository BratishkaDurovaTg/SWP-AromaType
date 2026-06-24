package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user User) (User, error) {
	err := r.db.QueryRow(ctx, `
INSERT INTO users (id, email, password_hash, role)
VALUES ($1, $2, $3, $4)
RETURNING id, email, password_hash, role, created_at, updated_at
`, user.ID, user.Email, user.PasswordHash, user.Role).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, ErrEmailAlreadyExists
		}
		return User{}, err
	}

	return user, nil
}

func (r *Repository) UpsertAdmin(ctx context.Context, user User) (User, error) {
	err := r.db.QueryRow(ctx, `
INSERT INTO users (id, email, password_hash, role)
VALUES ($1, $2, $3, $4)
ON CONFLICT (email) DO UPDATE SET
	password_hash = EXCLUDED.password_hash,
	role = EXCLUDED.role,
	updated_at = now()
RETURNING id, email, password_hash, role, created_at, updated_at
`, user.ID, user.Email, user.PasswordHash, user.Role).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := r.db.QueryRow(ctx, `
SELECT id, email, password_hash, role, created_at, updated_at
FROM users
WHERE email = $1
`, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrInvalidCredentials
	}
	if err != nil {
		return User{}, err
	}

	return user, nil
}
