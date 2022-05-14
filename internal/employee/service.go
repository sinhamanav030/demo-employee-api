package employee

import "context"

type Service interface {
	CreateEmployee(ctx context.Context, employee *Employee) error
	ListEmployee(ctx context.Context) ([]Employee, error)
	ListEmployeeById(ctx context.Context, id uint) (Employee, error)
	UpdateEmployee(ctx context.Context, employee *Employee) error
	DeleteEmployee(ctx context.Context, id string) error
	Migrations(ctx context.Context) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) CreateEmployee(ctx context.Context, employee *Employee) error {
	err := s.repo.Create(ctx, *employee)
	if err != nil {
		return err
	}
	return nil
}

func (s service) ListEmployee(ctx context.Context) ([]Employee, error) {
	emps, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return emps, nil
}

func (s service) ListEmployeeById(ctx context.Context, id uint) (Employee, error) {
	emp, err := s.repo.GetById(ctx, id)
	if err != nil {
		return Employee{}, err
	}
	return emp, err
}

func (s service) UpdateEmployee(ctx context.Context, emp *Employee) error {
	err := s.repo.Update(ctx, *emp)
	if err != nil {
		return err
	}
	return nil
}

func (s service) DeleteEmployee(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s service) Migrations(ctx context.Context) error {
	err := s.repo.Migrations(ctx)

	if err != nil {
		return err
	}

	return nil
}
