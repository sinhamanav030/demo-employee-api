package entity

const (
	RoleAdmin    = "admin"
	RoleEmployee = "employee"
)

type Employee struct {
	UserId       int    `json:"user_id,omitempty"`
	FirstName    string `json:"first_name,omitempty" validate:"required"`
	LastName     string `json:"last_name,omitempty" validate:"required"`
	Email        string `json:"email,omitempty" validate:"required,email"`
	Password     string `json:"password,omitempty" validate:"required,password"`
	Role         string `json:"role" validate:"required,role"`
	CreatedAt    string `json:"-"`
	LastAccessAt string `json:"-"`
	UpdatedAt    string `json:"-"`
	Archieved    bool   `json:"-"`
}
