package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"githb.com/demo-employee-api/internal/entity"
	"githb.com/demo-employee-api/pkg/customErrors"
)

type Repository interface {
	Create(ctx context.Context, employee entity.Employee) error
	Get(ctx context.Context, include_archieved bool, page int, perPage int) ([]entity.Employee, error)
	GetByParams(ctx context.Context, params map[string]string, page int, perPage int) ([]entity.Employee, error)
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

func (rep repository) Get(ctx context.Context, include_archieved bool, page int, perPage int) ([]entity.Employee, error) {
	emps := make([]entity.Employee, 0)
	query := fmt.Sprintf(`SELECT user_id,first_name,last_name,email,role FROM employees WHERE archieved=false or archieved=%v limit %d offset %d ;`, include_archieved, perPage, (page-1)*perPage)
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

func (rep repository) GetByParams(ctx context.Context, params map[string]string, page int, perPage int) ([]entity.Employee, error) {
	// emps := make([]entity.Employee, 0)

	// fmt.Println(params)
	query := "SELECT user_id,first_name,last_name,email,role FROM employees WHERE "

	if _, ok := params["user_id"]; ok {
		if idv, err := strconv.Atoi(params["user_id"]); err != nil {
			query = fmt.Sprintf(query, " user_id=%v AND ", idv)
		}

	}
	for key, value := range params {
		fmt.Println(key, value)
		if key == "user_id" || key == "page" || key == "per_page" || key == "sort_by" || key == "sort_order" {
			continue
		}
		if key != "archieved" {
			query = query + key + " ILIKE '%" + value + "%' AND "
		}
		fmt.Println(query)
	}

	perPageV := fmt.Sprintf("%d", (page-1)*perPage)
	pageV := fmt.Sprintf("%d", perPage)
	query = query + " archieved=" + params["archieved"]
	if v, ok := params["sort_by"]; ok {
		validColumn := false
		columns := []string{"user_id", "first_name", "last_name", "email", "created_at", "last_access_at", "updated_at"}
		for _, col := range columns {
			if v == col {
				validColumn = true
				break
			}
		}
		if validColumn {
			query = query + " order by " + v + " " + params["sort_order"]
		}
	}
	query = query + " offset " + perPageV + " limit " + pageV

	// fmt.Println(query)

	rows, err := rep.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	employees := make([]entity.Employee, 0)
	for rows.Next() {
		var employee entity.Employee
		rows.Scan(&employee.UserId, &employee.FirstName, &employee.LastName, &employee.Email, &employee.Role)
		employees = append(employees, employee)
	}
	if len(employees) == 0 {
		return nil, errors.New(customErrors.ErrorDataNotFound)
	}
	return employees, nil

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
