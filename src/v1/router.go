package v1

import (
	"net/http"

	"test-tablelink/src/v1/handler"
	"test-tablelink/src/v1/service"

	"github.com/go-chi/chi/v5"
)

func NewRouter(userService *service.UserService, roleService *service.RoleService, authService *service.AuthService) http.Handler {
	r := chi.NewRouter()

	userHandler := handler.NewUserHandler(userService)
	roleHandler := handler.NewRoleHandler(roleService)
	authHandler := handler.NewAuthHandler(authService)

	// Register routes
	userHandler.RegisterRoutes(r)
	roleHandler.RegisterRoutes(r)
	authHandler.RegisterRoutes(r)
	return r
}
