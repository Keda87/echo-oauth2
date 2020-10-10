package models

type User struct {
	ID       uint   `json:"id,omitempty" db:"id" sql:"id"`
	Email    string `json:"email" db:"email" sql:"email" validate:"required,email"`
	Password string `json:"-" db:"password" sql:"password" validate:"required"`
}

type UserPayload struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}
