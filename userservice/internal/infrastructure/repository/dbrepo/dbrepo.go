package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/wanomir/e"
	"time"
	"userservice/internal/entity"
)

const dbTimeout = time.Second * 3

type PostgresDBRepo struct {
	conn    *sql.DB
	timeout time.Duration
}

func WithDBTimeout(timeout time.Duration) func(*PostgresDBRepo) {
	return func(db *PostgresDBRepo) {
		db.timeout = timeout
	}

}

func NewPostgresDBRepo(conn *sql.DB, options ...func(repo *PostgresDBRepo)) *PostgresDBRepo {
	db := &PostgresDBRepo{
		conn:    conn,
		timeout: time.Second * 3, // default timout
	}

	for _, option := range options {
		option(db)
	}

	return db
}

func (db *PostgresDBRepo) Connection() *sql.DB {
	return db.conn
}

func (db *PostgresDBRepo) ListUsers(ctx context.Context, offset, limit int) ([]entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, db.timeout)
	defer cancel()

	query := `SELECT id, email, password, name, age
				 FROM users 
				 WHERE is_deleted = FALSE
				 LIMIT $1 OFFSET $2`

	var users []entity.User
	rows, err := db.conn.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, e.WrapIfErr("error executing query", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.User
		err = rows.Scan(
			&user.Id,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.Age,
		)
		if err != nil {
			return nil, e.WrapIfErr("error scanning row", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (db *PostgresDBRepo) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, db.timeout)
	defer cancel()

	query := `SELECT id, email, password, name, age
				 FROM users
				 WHERE email = $1 AND is_deleted = FALSE`

	var user entity.User
	err := db.conn.QueryRowContext(ctx, query, email).Scan(
		&user.Id, &user.Email, &user.Password, &user.Name, &user.Age,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.User{}, errors.New("user not found")
	} else if err != nil {
		return entity.User{}, e.WrapIfErr("failed to execute query", err)
	}

	return user, nil
}

func (db *PostgresDBRepo) CreateUser(ctx context.Context, user entity.User) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, db.timeout)
	defer cancel()

	query := `INSERT INTO users (email, password, name, age, is_deleted) 
				 VALUES ($1, $2, $3, $4, FALSE)
				 RETURNING id`

	var userId int
	err := db.conn.QueryRowContext(ctx, query,
		user.Email,
		user.Password,
		user.Name,
		user.Age,
	).Scan(&userId)

	if err != nil {
		return 0, e.WrapIfErr("failed to execute query", err)
	}

	return userId, nil
}

func (db *PostgresDBRepo) UpdateUser(ctx context.Context, user entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, db.timeout)
	defer cancel()

	query := `UPDATE users 
				 SET password = $1, name = $2, age = $3
				 WHERE email = $4`

	if _, err := db.conn.ExecContext(ctx, query,
		user.Password,
		user.Name,
		user.Age,
		user.Email,
	); err != nil {
		return e.WrapIfErr("error executing query", err)
	}

	return nil
}

func (db *PostgresDBRepo) DeleteUser(ctx context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctx, db.timeout)
	defer cancel()

	query := `UPDATE users 
				 SET is_deleted = TRUE
				 WHERE email = $1`

	if _, err := db.conn.ExecContext(ctx, query, email); err != nil {
		return e.WrapIfErr("error executing query", err)
	}

	return nil
}
