package service

import (
	"context"
	"strconv"
	"test-tablelink/src/entity"
	"test-tablelink/src/repository"
	"test-tablelink/src/v1/contract"
)

type RoleService struct {
	roleRepo *repository.RoleRepository
}

func NewRoleService(roleRepo *repository.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (s *RoleService) GetRole(ctx context.Context, id string) (*contract.RoleResponse, error) {
	roleID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	return &contract.RoleResponse{
		Status:  true,
		Message: "Successfully",
		Data:    role,
	}, nil
}

func (s *RoleService) CreateRole(ctx context.Context, req *contract.CreateRoleRequest) (*contract.RoleResponse, error) {
	role := &entity.Role{
		Name: req.Name,
	}

	err := s.roleRepo.Create(ctx, role)
	if err != nil {
		return nil, err
	}

	return &contract.RoleResponse{
		Status:  true,
		Message: "Successfully",
		Data:    role,
	}, nil
}

func (s *RoleService) UpdateRole(ctx context.Context, id string, req *contract.UpdateRoleRequest) (*contract.RoleResponse, error) {
	roleID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	role.Name = req.Name

	err = s.roleRepo.Update(ctx, role)
	if err != nil {
		return nil, err
	}

	return &contract.RoleResponse{
		Status:  true,
		Message: "Successfully",
		Data:    role,
	}, nil
}

func (s *RoleService) DeleteRole(ctx context.Context, id string) (*contract.RoleResponse, error) {
	roleID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	err = s.roleRepo.Delete(ctx, roleID)
	if err != nil {
		return nil, err
	}

	return &contract.RoleResponse{
		Status:  true,
		Message: "Successfully",
	}, nil
}
