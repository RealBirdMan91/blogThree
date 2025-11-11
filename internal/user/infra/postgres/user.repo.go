package postgres

import (
	"blogThree/internal/user/app"
	"blogThree/internal/user/domain"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresUserRepo struct{ db *sql.DB }

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo { return &PostgresUserRepo{db: db} }

var _ app.UserRepository = (*PostgresUserRepo)(nil)

func (r *PostgresUserRepo) CreateUser(ctx context.Context, user *domain.User) error {
	const query = `
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID(),
		user.Email().String(),
		user.PasswordHash().String(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)

	if err != nil {
		if isEmailUniqueViolation(err) {
			return app.ErrEmailAlreadyExists
		}
		return app.NewUserInsertFailed(err)
	}

	return nil
}

func (r *PostgresUserRepo) ExistsByEmail(ctx context.Context, email domain.Email) (bool, error) {
	const q = `SELECT 1 FROM users WHERE lower(email)=lower($1) LIMIT 1`
	var one int
	err := r.db.QueryRowContext(ctx, q, email.String()).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, app.NewUserExistsCheckFailed(err)
	}

	return true, nil
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE lower(email) = lower($1)
		LIMIT 1
	`

	var (
		id        uuid.UUID
		emailStr  string
		hashStr   string
		createdAt time.Time
		updatedAt time.Time
	)

	err := r.db.QueryRowContext(ctx, q, email.String()).
		Scan(&id, &emailStr, &hashStr, &createdAt, &updatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, app.ErrUserNotFound
	}
	if err != nil {
		return nil, app.NewUserSelectFailed(err)
	}
	return r.hydrateUser(id, emailStr, hashStr, createdAt, updatedAt)
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	var (
		emailStr  string
		hashStr   string
		createdAt time.Time
		updatedAt time.Time
	)

	err := r.db.QueryRowContext(ctx, q, id).
		Scan(&id, &emailStr, &hashStr, &createdAt, &updatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, app.ErrUserNotFound
	}
	if err != nil {
		return nil, app.NewUserSelectFailed(err)
	}

	return r.hydrateUser(id, emailStr, hashStr, createdAt, updatedAt)
}

func (r *PostgresUserRepo) List(ctx context.Context) ([]*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, app.NewUserListFailed(err)
	}
	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		var (
			id        uuid.UUID
			emailStr  string
			hashStr   string
			createdAt time.Time
			updatedAt time.Time
		)

		if err := rows.Scan(&id, &emailStr, &hashStr, &createdAt, &updatedAt); err != nil {
			return nil, app.NewUserListFailed(err)
		}

		if err := rows.Scan(&id, &emailStr, &hashStr, &createdAt, &updatedAt); err != nil {
			return nil, app.NewUserListFailed(err)
		}

		u, err := r.hydrateUser(id, emailStr, hashStr, createdAt, updatedAt)
		if err != nil {
			return nil, app.NewUserListFailed(err)
		}

		users = append(users, u)

	}

	if err := rows.Err(); err != nil {
		return nil, app.NewUserListFailed(err)
	}

	return users, nil
}

//helpers

func isEmailUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_lower_unique"
}

func (r *PostgresUserRepo) hydrateUser(
	id uuid.UUID,
	emailStr, hashStr string,
	createdAt, updatedAt time.Time,
) (*domain.User, error) {
	emailVO, err := domain.NewEmail(emailStr)
	if err != nil {
		return nil, app.NewUserSelectFailed(err)
	}

	hashVO, err := domain.NewPasswordHash(hashStr)
	if err != nil {
		return nil, app.NewUserSelectFailed(err)
	}

	user, err := domain.RehydrateUser(id, emailVO, hashVO, createdAt, updatedAt)
	if err != nil {
		return nil, app.NewUserSelectFailed(err)
	}

	return user, nil
}
