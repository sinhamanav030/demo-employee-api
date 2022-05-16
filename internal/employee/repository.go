package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"githb.com/demo-employee-api/internal/entity"
	"githb.com/demo-employee-api/pkg/customErrors"
)

type Repository interface {
	Create(ctx context.Context, employee entity.Employee) error
	GetAll(ctx context.Context, archieved bool) ([]entity.Employee, error)
	GetById(ctx context.Context, id uint, archieved bool) (entity.Employee, error)
	Update(ctx context.Context, employee entity.Employee) error
	Delete(ctx context.Context, id string) error
	Migrations(ctx context.Context) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return repository{db}
}

func (rep repository) Create(ctx context.Context, emp entity.Employee) error {
	query := fmt.Sprintf(`
	insert into employees(first_name,last_name,email,password,role,created_at,last_access_at,archieved) 
	values('%v','%v','%v','%v','%v','now','now',FALSE);`, emp.FirstName, emp.LastName, emp.Email, emp.Password, emp.Role)
	result, err := rep.db.ExecContext(ctx, query)
	if err != nil {
		// log.Println(err)
		return err
	}
	rc, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rc == 0 {
		return errors.New(customErrors.ErrorInternalServer)
	}

	log.Println(rc)
	return nil

}

func (rep repository) GetAll(ctx context.Context, archieved bool) ([]entity.Employee, error) {
	emps := make([]entity.Employee, 0)
	query := fmt.Sprintf(`SELECT user_id,first_name,last_name,email,role FROM employees WHERE archieved=%v;`, archieved)
	rows, err := rep.db.QueryContext(ctx, query)
	if err != nil {
		return emps, err
	}
	for rows.Next() {
		var employee entity.Employee
		rows.Scan(&employee.UserId, &employee.FirstName, &employee.LastName, &employee.Email, &employee.Role)
		emps = append(emps, employee)
	}
	return emps, nil
}

func (rep repository) GetById(ctx context.Context, id uint, archieved bool) (entity.Employee, error) {
	// emps := make([]entity.Employee, 0)
	query := fmt.Sprintf("SELECT user_id,first_name,last_name,email,role FROM employees WHERE user_id=%v AND archieved=%v;", id, archieved)
	rows, err := rep.db.QueryContext(ctx, query)
	// rows, err := rep.db.QueryContext(ctx, "SELECT USER_ID,FIRST_NAME,LAST_NAME,EMAIL FROM EMPLOYEES WHERE USER_ID=? AND ARCHIEVED=FALSE;", id)
	if err != nil {
		return entity.Employee{}, err
	}
	var employee entity.Employee
	for rows.Next() {
		rows.Scan(&employee.UserId, &employee.FirstName, &employee.LastName, &employee.Email, &employee.Role)
	}
	if employee.UserId == 0 {
		return entity.Employee{}, errors.New(customErrors.ErrorDataNotFound)
	}
	return employee, nil

}

func (rep repository) Update(ctx context.Context, emp entity.Employee) error {
	query := fmt.Sprintf(`
	UPDATE employees 
	SET first_name = '%v',
	last_name= '%v',
	role = '%v',
	updated_at = 'now',
	last_access_at = 'now'
	WHERE user_id =%v AND archieved=false
	`, emp.FirstName, emp.LastName, emp.Role, emp.UserId)
	result, err := rep.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	rc, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rc == 0 {
		return errors.New(customErrors.ErrorDataNotFound)
	}

	fmt.Println(rc)
	return nil

}

func (rep repository) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("UPDATE employees SET archieved=true WHERE user_id=%v AND archieved=false AND NOT role='admin'", id)
	result, err := rep.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	rc, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println(rc)

	if rc == 0 {
		return errors.New(customErrors.ErrorDataNotFound)
	}
	return nil
}

func (rep repository) Migrations(ctx context.Context) error {
	result, err := rep.db.ExecContext(ctx, `CREATE TABLE employees(
		user_id serial PRIMARY KEY,
		first_name VARCHAR(15) NOT NULL,
		last_name VARCHAR(15) NOT NULL,
		email VARCHAR(50) UNIQUE NOT NULL,
		role VARCHAR(15) NOT NULL,
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
