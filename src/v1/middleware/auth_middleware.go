package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"test-tablelink/src/entity"
	"test-tablelink/src/repository"
	"test-tablelink/src/v1/service"
)

type AuthMiddleware struct {
	authService   *service.AuthService
	roleRightRepo *repository.RoleRightRepository
}

func NewAuthMiddleware(authService *service.AuthService, roleRightRepo *repository.RoleRightRepository) *AuthMiddleware {
	return &AuthMiddleware{
		authService:   authService,
		roleRightRepo: roleRightRepo,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check X-Link-Service header
		section := r.Header.Get("section")
		if section != "be" {
			http.Error(w, "Invalid service section", http.StatusUnauthorized)
			return
		}

		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Extract token from Bearer
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}
		token := parts[1]

		// Validate token and get user
		user, err := m.authService.ValidateToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user from context
		user, ok := r.Context().Value("user").(*entity.User)
		if !ok {
			http.Error(w, "User not found in context", http.StatusUnauthorized)
			return
		}

		// Get route pattern
		routePattern := r.URL.Path
		log.Printf("Checking permission for route: %s, method: %s, user role: %d", routePattern, r.Method, user.RoleID)

		// Check permission
		hasPermission, err := m.roleRightRepo.CheckPermission(
			r.Context(),
			user.RoleID,
			"be",
			routePattern,
			r.Method,
		)
		if err != nil {
			log.Printf("Error checking permission: %v", err)
			http.Error(w, fmt.Sprintf("Error checking permission: %v", err), http.StatusInternalServerError)
			return
		}
		if !hasPermission {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
