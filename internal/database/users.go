package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"lesson-proj/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (userRepository *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, name, hashed_password
		FROM users
		WHERE email = $1;
	`
	err := userRepository.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.HashedPassword,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("user with email %s not found", email)
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (userRepository *UserRepository) GetAllUsers(ctx context.Context) ([]models.UserWithoutPassword, error) {
	var users []models.UserWithoutPassword
	query := `
		SELECT id, email, name FROM users;`
	rows, err := userRepository.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.UserWithoutPassword
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (userRepository *UserRepository) GetUserByID(ctx context.Context, id int) (*models.UserWithoutPassword, error) {
	var user models.UserWithoutPassword
	query := `
		SELECT id, email, name
		FROM users
		WHERE id = $1;
	`
	err := userRepository.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (userRepository *UserRepository) CreateUser(ctx context.Context, inputUser models.CreateUser) (*models.UserWithoutPassword, error) {
	var user models.UserWithoutPassword

	query := `
		INSERT INTO users (email, name, hashed_password)
		VALUES ($1, $2, $3)
		RETURNING id, email, name;`
	err := userRepository.db.QueryRow(ctx, query,
		inputUser.Email,
		inputUser.Name,
		inputUser.Password,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepository *UserRepository) UpdateUser(ctx context.Context, id int, inputUser models.UpdateUser) (*models.UserWithoutPassword, error) {

	query := `
		UPDATE users
		SET
			email = COALESCE($1, email),
			name  = COALESCE($2, name),
			hashed_password = COALESCE($3, hashed_password)
		WHERE id = $4
		RETURNING id, email, name;
	`
	var updatedUser models.UserWithoutPassword
	err := userRepository.db.QueryRow(
		ctx,
		query,
		inputUser.Email,
		inputUser.Name,
		inputUser.Password,
		id,
	).Scan(
		&updatedUser.ID,
		&updatedUser.Email,
		&updatedUser.Name,
	)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (userRepository *UserRepository) DeleteUser(ctx context.Context, id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1;`

	result, err := userRepository.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	rows := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}
	return nil
}
