package models

type User struct {
	// name for json
	ID       int    `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Name     string `json:"name" db:"name"`
	HashedPassword string `json:"password" db:"hashed_password"`
}
type UserWithoutPassword struct {
	ID    int    `json:"id" db:"id"`
	Email string `json:"email" db:"email"`
	Name  string `json:"name" db:"name"`
}

type CreateUser struct {
	Email    string `json:"email" db:"email"`
	Name     string `json:"name" db:"name"`
	Password string `json:"password" db:"hashed_password"`
}

type UpdateUser struct {
	// pointer to a string, can be nil if the field is not provided
	Email    *string `json:"email" db:"email"`
	Name     *string `json:"name" db:"name"`
	Password *string `json:"password" db:"hashed_password"`
}

type AuthUser struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"hashed_password"`
}
