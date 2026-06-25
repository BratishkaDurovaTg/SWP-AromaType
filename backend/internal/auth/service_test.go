package auth

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeUserRepository struct {
	usersByEmail map[string]User
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{usersByEmail: map[string]User{}}
}

func (r *fakeUserRepository) CreateUser(_ context.Context, user User) (User, error) {
	if _, exists := r.usersByEmail[user.Email]; exists {
		return User{}, ErrEmailAlreadyExists
	}
	now := time.Date(2026, 6, 25, 12, 0, 0, 0, time.UTC)
	user.CreatedAt = now
	user.UpdatedAt = now
	r.usersByEmail[user.Email] = user
	return user, nil
}

func (r *fakeUserRepository) UpsertAdmin(_ context.Context, user User) (User, error) {
	now := time.Date(2026, 6, 25, 12, 0, 0, 0, time.UTC)
	if existing, ok := r.usersByEmail[user.Email]; ok {
		user.ID = existing.ID
		user.CreatedAt = existing.CreatedAt
		user.UpdatedAt = now
	} else {
		user.CreatedAt = now
		user.UpdatedAt = now
	}
	r.usersByEmail[user.Email] = user
	return user, nil
}

func (r *fakeUserRepository) FindUserByEmail(_ context.Context, email string) (User, error) {
	user, ok := r.usersByEmail[email]
	if !ok {
		return User{}, ErrInvalidCredentials
	}
	return user, nil
}

func TestRegisterNormalizesEmailAndIssuesUserToken(t *testing.T) {
	repo := newFakeUserRepository()
	service := NewService(repo, "test-secret", time.Hour)

	result, err := service.Register(context.Background(), " USER@Example.COM ", "secret123")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if result.User.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", result.User.Email)
	}
	if result.User.Role != RoleUser {
		t.Fatalf("expected user role, got %q", result.User.Role)
	}
	if result.User.PasswordHash == "secret123" {
		t.Fatal("password was stored without hashing")
	}

	claims, err := service.VerifyToken(result.AccessToken)
	if err != nil {
		t.Fatalf("VerifyToken returned error: %v", err)
	}
	if claims.UserID != result.User.ID || claims.Email != result.User.Email || claims.Role != RoleUser {
		t.Fatalf("token claims do not match user: %#v", claims)
	}
}

func TestRegisterRejectsInvalidInput(t *testing.T) {
	service := NewService(newFakeUserRepository(), "test-secret", time.Hour)

	_, err := service.Register(context.Background(), "user@example.com", "123")
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	service := NewService(newFakeUserRepository(), "test-secret", time.Hour)
	if _, err := service.Register(context.Background(), "user@example.com", "secret123"); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	_, err := service.Login(context.Background(), "user@example.com", "wrong-password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestEnsureAdminCreatesAdminWithStrongPassword(t *testing.T) {
	service := NewService(newFakeUserRepository(), "test-secret", time.Hour)

	admin, err := service.EnsureAdmin(context.Background(), " ADMIN@Example.COM ", "very-strong-password")
	if err != nil {
		t.Fatalf("EnsureAdmin returned error: %v", err)
	}
	if admin.Email != "admin@example.com" {
		t.Fatalf("expected normalized admin email, got %q", admin.Email)
	}
	if admin.Role != RoleAdmin {
		t.Fatalf("expected admin role, got %q", admin.Role)
	}
}

func TestVerifyTokenRejectsMalformedToken(t *testing.T) {
	service := NewService(newFakeUserRepository(), "test-secret", time.Hour)

	_, err := service.VerifyToken("not-a-jwt")
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}
