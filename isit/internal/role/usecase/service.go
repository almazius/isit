package usecase

import (
	"context"
	"isit/isit/internal/role/models"
)

type RoleService struct {
	repo interface{}
}

func NewRoleService() *RoleService {
	return &RoleService{
		nil,
	}
}

func (rs *RoleService) CreateRole(ctx context.Context, role *models.Role) error {
	return nil
}

func (rs *RoleService) CreatePermission(ctx context.Context, role *models.Permission) error {
	return nil
}

func (rs *RoleService) GrantPermission(ctx context.Context, params *models.GrantPermission) error {
	return nil
}
