package employee

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"githb.com/demo-employee-api/internal/entity"
	"githb.com/demo-employee-api/pkg/customErrors"
	"githb.com/demo-employee-api/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateEmployee(ctx context.Context, employee *entity.Employee) error
	ListEmployee(ctx context.Context, params map[string]string) ([]entity.Employee, error)
	ListEmployeeByParams(ctx context.Context, params map[string]string) ([]entity.Employee, error)
	UpdateEmployee(ctx context.Context, id int, employee *UpdateEmpReq) error
	DeleteEmployee(ctx context.Context, id string) error
	LoginEmployee(ctx context.Context, employee *LoginReq) (string, error)
}

type service struct {
	repo       Repository
	logger     *log.Logger
	tokenMaker token.Maker
}

func NewService(repo Repository, logger *log.Logger, tokenMaker token.Maker) Service {
	return service{repo, logger, tokenMaker}
}

func (s service) LoginEmployee(ctx context.Context, employee *LoginReq) (string, error) {

	emp, err := s.repo.GetByEmail(ctx, employee.Email)
	if err != nil {
		s.logger.Println(err)
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(emp.Password), []byte(employee.Password))
	if err != nil {
		s.logger.Println(err)
		return "", errors.New(customErrors.ErrorAuthFailed)
	}
	token, _, err := s.tokenMaker.CreateToken(emp.UserId, emp.Role, time.Duration(5*time.Minute))

	if err != nil {
		s.logger.Println(err)
		return "", errors.New(customErrors.ErrorInternalServer)
	}
	return token, nil
}

func (s service) CreateEmployee(ctx context.Context, emp *entity.Employee) error {

	passByteSlice, err := bcrypt.GenerateFromPassword([]byte(emp.Password), bcrypt.MinCost)

	if err != nil {
		s.logger.Println(err)
		return err
	}

	emp.Password = string(passByteSlice)

	err = s.repo.Create(ctx, *emp)
	if err != nil {
		s.logger.Println(err)
	}
	return nil
}

func (s service) ListEmployee(ctx context.Context, params map[string]string) ([]entity.Employee, error) {

	var page, perPage int
	var err error
	if page, err = strconv.Atoi(params["page"]); err != nil {
		s.logger.Println(err)
		return nil, errors.New(customErrors.ErrorInternalServer)
	}
	if perPage, err = strconv.Atoi(params["per_page"]); err != nil {
		s.logger.Println(err)
		return nil, errors.New(customErrors.ErrorInternalServer)
	}
	emps, err := s.repo.Get(ctx, params["archieved"] == "true", page, perPage)
	if err != nil {
		s.logger.Println(err)
		return nil, err
	}
	return emps, nil
}

func (s service) ListEmployeeByParams(ctx context.Context, params map[string]string) ([]entity.Employee, error) {

	var (
		page, perPage int
		err           error
	)

	if page, err = strconv.Atoi(params["page"]); err != nil {
		s.logger.Println(err)
		return nil, errors.New(customErrors.ErrorInternalServer)
	}
	if perPage, err = strconv.Atoi(params["per_page"]); err != nil {
		s.logger.Println(err)
		return nil, errors.New(customErrors.ErrorInternalServer)
	}
	emps, err := s.repo.GetByParams(ctx, params, page, perPage)
	if err != nil {
		s.logger.Println(err)
		return nil, err
	}
	return emps, err
}

func (s service) UpdateEmployee(ctx context.Context, id int, emp *UpdateEmpReq) error {

	employee := entity.Employee{
		UserId:    id,
		FirstName: emp.FirstName,
		LastName:  emp.LastName,
		Role:      emp.Role,
	}

	err := s.repo.Update(ctx, employee)

	if err != nil {
		s.logger.Println(err)
		return err
	}
	return nil
}

func (s service) DeleteEmployee(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Println(err)
		return err
	}
	return nil
}
