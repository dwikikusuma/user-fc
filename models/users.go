package models

type (
	RegisterParameter struct {
		Name            string `json:"name" validate:"required"`
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required,min=8"`
		ConfirmPassword string `json:"confirm_password" validate:"required,min=68`
	}

	LoginParameter struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"password"`
}
