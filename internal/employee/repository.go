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
	GetByEmail(ctx context.Context, email string) (entity.Employee, error)
}

type repository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewRepository(db *sql.DB, logger *log.Logger) Repository {
	return repository{db, logger}
}

func (rep repository) Create(ctx context.Context, emp entity.Employee) error {
	query := fmt.Sprintf(`
	insert into employees(first_name,last_name,email,password,role,created_at,last_access_at,archieved) 
	values('%v','%v','%v','%v','%v','now','now',FALSE);`, emp.FirstName, emp.LastName, emp.Email, emp.Password, emp.Role)
	result, err := rep.db.ExecContext(ctx, query)
	if err != nil {
		rep.logger.Println(err)
		return errors.New(customErrors.ErrorUserExists)
	}
	rc, err := result.RowsAffected()
	if err != nil {
		rep.logger.Println(err)
		return errors.New(customErrors.ErrorInternalServer)
	}

	if rc == 0 {
		rep.logger.Println(err)
		return errors.New(customErrors.ErrorInternalServer)
	}

	rep.logger.Printf("%d rows inseterd ", rc)
	return nil

}

func (rep repository) Get(ctx context.Context, include_archieved bool, page, perPage int) ([]entity.Employee, error) {
	emps := make([]entity.Employee, 0)
	query := fmt.Sprintf(`SELECT user_id,first_name,last_name,email,role FROM employees WHERE archieved=false or archieved=%v limit %d offset %d ;`, include_archieved, perPage, (page-1)*perPage)
	rows, err := rep.db.QueryContext(ctx, query)
	if err != nil {
		rep.logger.Println(err)
		return emps, errors.New(customErrors.ErrorInternalServer)
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

	if id, ok := params["user_id"]; ok {
		if _, err := strconv.Atoi(params["user_id"]); err == nil {
			query = query + " user_id=" + id + " AND"
		} else {
			rep.logger.Println(err)
		}

	}
	for key, value := range params {
		if key == "user_id" || key == "page" || key == "per_page" || key == "sort_by" || key == "sort_order" {
			continue
		}
		if key != "archieved" {
			query = query + key + " ILIKE '%" + value + "%' AND "
		}
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

	rep.logger.Printf("Executed Query %v", query)

	rows, err := rep.db.QueryContext(ctx, query)
	if err != nil {
		rep.logger.Println(err)
		return nil, errors.New(customErrors.ErrorInternalServer)
	}
	employees := make([]entity.Employee, 0)
	for rows.Next() {
		var employee entity.Employee
		rows.Scan(&employee.UserId, &employee.FirstName, &employee.LastName, &employee.Email, &employee.Role)
		employees = append(employees, employee)
	}
	if len(employees) == 0 {
		rep.logger.Println(err)
		return nil, errors.New(customErrors.ErrorDataNotFound)
	}
	return employees, nil

}

func (rep repository) GetByEmail(ctx context.Context, email string) (entity.Employee, error) {
	query := fmt.Sprintf(`
	SELECT user_id,email,password,role FROM employees 
	WHERE email='%v' and archieved=false;
	`, email)
	row, err := rep.db.QueryContext(ctx, query)

	if err != nil {
		rep.logger.Println(err)
		return entity.Employee{}, errors.New(customErrors.ErrorInternalServer)
	}

	var emp entity.Employee
	for row.Next() {
		row.Scan(&emp.UserId, &emp.Email, &emp.Password, &emp.Role)
	}
	if emp.UserId == 0 {
		return entity.Employee{}, errors.New(customErrors.ErrorDataNotFound)
	}

	query = fmt.Sprintf(`
		UPDATE employees
		SET last_access_at='now'
		WHERE email='%v' AND archieved=false' 
	`, email)
	_, _ = rep.db.QueryContext(ctx, query)
	return emp, nil
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
		rep.logger.Println(err)
		return errors.New(customErrors.ErrorInternalServer)
	}

	rc, err := result.RowsAffected()
	if err != nil {
		rep.logger.Println(err)
		return errors.New(customErrors.ErrorInternalServer)
	}

	if rc == 0 {
		return errors.New(customErrors.ErrorDataNotFound)
	}

	rep.logger.Printf("%d rows updated\n", rc)
	return nil

}

func (rep repository) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("UPDATE employees SET archieved=true WHERE user_id=%v AND archieved=false AND NOT role='admin'", id)
	result, err := rep.db.ExecContext(ctx, query)
	if err != nil {
		rep.logger.Println(err)
		return errors.New(customErrors.ErrorInternalServer)
	}
	rc, err := result.RowsAffected()
	if err != nil {
		rep.logger.Println(err)
		return errors.New(customErrors.ErrorInternalServer)
	}
	rep.logger.Printf("%d rows updated\n", rc)

	if rc == 0 {
		return errors.New(customErrors.ErrorDataNotFound)
	}
	return nil
}
