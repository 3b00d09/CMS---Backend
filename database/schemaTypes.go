package database

type User struct {
	ID       string `db:"id"`
	Username string `db:"username"`
}

// validate required tag is important for the library that makes sure the request body has all fields
type UserCredentials struct {
	ID        string `json:"id" db:"id"`
	Username  string `json:"username" db:"username" validate:"required"`
	Password  string `json:"password" db:"password" validate:"required"`
	FirstName string `json:"firstName" db:"first_name" validate:"required"`
	LastName  string `json:"lastName" db:"last_name" validate:"required"`
	Email     string `json:"email" db:"email" validate:"required,email"`
}

type UserSession struct {
	ID            string `db:"id"`
	UserID        string `db:"user_id"`
	ActiveExpires int64  `db:"active_expires"`
	IdleExpires   int64  `db:"idle_expires"`
}

type Project struct {
	ID          string `json:"id" db:"id"`
	Username    string `json:"username" db:"username" validate:"required"`
	ProjectName string `json:"project_name" db:"name" validate:"required"`
}

type Page struct {
	ID        string `json:"id" db:"id"`
	ProjectID string `json:"project_id" db:"project_id"`
	Name      string `json:"name" db:"name"`
	Content   string `json:"content" db:"content"`
}