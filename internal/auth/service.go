package auth

import (
	"context"

	"githb.com/demo-employee-api/internal/entity"
)

type Service interface {
	LoginEmployee(ctx context.Context, email string, password string) (entity.Employee, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) LoginEmployee(ctx context.Context, email string, password string) (entity.Employee, error) {
	emp, err := s.repo.GetByEmail(ctx, email, password)
	if err != nil {
		return entity.Employee{}, err
	}
	return emp, err
}
