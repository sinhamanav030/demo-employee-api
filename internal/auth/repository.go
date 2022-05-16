package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"githb.com/demo-employee-api/internal/entity"
	"githb.com/demo-employee-api/pkg/customErrors"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	GetByEmail(ctx context.Context, email string, password string) (entity.Employee, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return repository{db}
}

func (rep repository) GetByEmail(ctx context.Context, email string, password string) (entity.Employee, error) {
	query := fmt.Sprintf(`
	SELECT user_id,email,password,role FROM employees 
	WHERE email='%v' and archieved=false;
	`, email)
	row, err := rep.db.QueryContext(ctx, query)

	if err != nil {
		return entity.Employee{}, err
	}

	var emp entity.Employee
	r := 0
	var pass string
	for row.Next() {
		r++
		row.Scan(&emp.UserId, &emp.Email, &pass, &emp.Role)
	}
	if emp.UserId == 0 {
		return entity.Employee{}, errors.New(customErrors.ErrorDataNotFound)
	}

	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))

	if err != nil {
		return entity.Employee{}, errors.New(customErrors.ErrorAuthFailed)
	}

	query = fmt.Sprintf(`
		UPDATE employees
		SET last_access_at='now'
		WHERE email='%v' AND archieved=false' 
	`, email)
	_, _ = rep.db.QueryContext(ctx, query)

	return emp, nil
}
