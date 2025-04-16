package handler

import (
	"encoding/json"
	"net/http"

	"test-tablelink/src/v1/contract"
	"test-tablelink/src/v1/service"

	"github.com/go-chi/chi/v5"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func (h *RoleHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")
	if roleID == "" {
		http.Error(w, "Role ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.roleService.GetRole(r.Context(), roleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req contract.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.roleService.CreateRole(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")
	if roleID == "" {
		http.Error(w, "Role ID is required", http.StatusBadRequest)
		return
	}

	var req contract.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.roleService.UpdateRole(r.Context(), roleID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")
	if roleID == "" {
		http.Error(w, "Role ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.roleService.DeleteRole(r.Context(), roleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *RoleHandler) RegisterRoutes(r chi.Router) {
	r.Get("/roles/{id}", h.GetRole)
	r.Post("/roles", h.CreateRole)
	r.Put("/roles/{id}", h.UpdateRole)
	r.Delete("/roles/{id}", h.DeleteRole)
}
