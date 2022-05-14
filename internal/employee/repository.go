package employee

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Employee struct {
	UserId       int    `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	LastAccessAt string `json:"last_access_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
	Archieved    bool   `json:"archieved,omitempty"`
}

type Repository interface {
	Create(ctx context.Context, employee Employee) error
	GetAll(ctx context.Context) ([]Employee, error)
	GetById(ctx context.Context, id uint) (Employee, error)
	Update(ctx context.Context, employee Employee) error
	Delete(ctx context.Context, id string) error
	Migrations(ctx context.Context) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return repository{db}
}

func (rep repository) Create(ctx context.Context, emp Employee) error {
	query := fmt.Sprintf(`
	insert into employees(first_name,last_name,email,password,created_at,last_access_at,archieved) 
	values('%v','%v','%v','%v','now','now',FALSE);`, emp.FirstName, emp.LastName, emp.Email, emp.Password)
	result, err := rep.db.ExecContext(ctx, query)
	if err != nil {
		// log.Println(err)
		return err
	}
	rc, err := result.RowsAffected()
	if err != nil {
		return err
	}
	log.Println(rc)
	return nil

}

func (rep repository) GetAll(ctx context.Context) ([]Employee, error) {
	emps := make([]Employee, 0)
	rows, err := rep.db.QueryContext(ctx, `SELECT USER_ID,FIRST_NAME,LAST_NAME,EMAIL FROM EMPLOYEES WHERE ARCHIEVED=FALSE;`)
	if err != nil {
		return emps, err
	}
	for rows.Next() {
		var employee Employee
		rows.Scan(&employee.UserId, &employee.FirstName, &employee.LastName, &employee.Email)
		emps = append(emps, employee)
	}
	return emps, nil
}

func (rep repository) GetById(ctx context.Context, id uint) (Employee, error) {
	// emps := make([]Employee, 0)
	query := fmt.Sprintf("SELECT USER_ID,FIRST_NAME,LAST_NAME,EMAIL FROM EMPLOYEES WHERE USER_ID=%v AND ARCHIEVED=FALSE;", id)
	rows, err := rep.db.QueryContext(ctx, query)
	// rows, err := rep.db.QueryContext(ctx, "SELECT USER_ID,FIRST_NAME,LAST_NAME,EMAIL FROM EMPLOYEES WHERE USER_ID=? AND ARCHIEVED=FALSE;", id)
	if err != nil {
		return Employee{}, err
	}
	var employee Employee
	for rows.Next() {
		rows.Scan(&employee.UserId, &employee.FirstName, &employee.LastName, &employee.Email)
		// emps = append(emps, employee)
		// fmt.Println("in here")
	}
	return employee, nil

}

func (rep repository) Update(ctx context.Context, emp Employee) error {
	query := fmt.Sprintf(`
	UPDATE employees 
	SET first_name = '%v',
	last_name= '%v',
	email = '%v',
	updated_at = 'now',
	last_access_at = 'now'
	WHERE user_id =%v
	`, emp.FirstName, emp.LastName, emp.Email, emp.UserId)
	result, err := rep.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	rc, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println(rc)
	return nil

}

func (rep repository) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("UPDATE employees SET archieved=true WHERE user_id=%v", id)
	result, err := rep.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	rc, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println(rc)
	return nil
}

func (rep repository) Migrations(ctx context.Context) error {
	result, err := rep.db.Exec(`CREATE TABLE employees(
		user_id serial PRIMARY KEY,
		first_name VARCHAR(15) NOT NULL,
		last_name VARCHAR(15) NOT NULL,
		email VARCHAR(20) UNIQUE NOT NULL,
		password VARCHAR(256) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		last_access_at TIMESTAMP,
		updated_at TIMESTAMP,
		archieved BOOLEAN
	);`)

	if err != nil {
		return err
	}

	rc, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println(rc)
	return nil
}
