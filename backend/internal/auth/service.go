package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AuthResult struct {
	AccessToken string
	User        User
}

type Service struct {
	repo      *Repository
	jwtSecret []byte
	jwtTTL    time.Duration
}

func NewService(repo *Repository, jwtSecret string, jwtTTL time.Duration) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
		jwtTTL:    jwtTTL,
	}
}

func (s *Service) Register(ctx context.Context, email, password, role string) (AuthResult, error) {
	email = normalizeEmail(email)
	role = normalizeRole(role)

	if email == "" || len(password) < 6 || !isAllowedRole(role) {
		return AuthResult{}, ErrInvalidInput
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, err
	}

	user := User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: string(passwordHash),
		Role:         role,
	}

	created, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			return AuthResult{}, ErrEmailAlreadyExists
		}
		return AuthResult{}, err
	}

	token, err := s.issueToken(created)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{AccessToken: token, User: created}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (AuthResult, error) {
	email = normalizeEmail(email)
	if email == "" || password == "" {
		return AuthResult{}, ErrInvalidCredentials
	}

	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return AuthResult{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return AuthResult{}, ErrInvalidCredentials
	}

	token, err := s.issueToken(user)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{AccessToken: token, User: user}, nil
}

func (s *Service) issueToken(user User) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func normalizeRole(role string) string {
	role = strings.ToLower(strings.TrimSpace(role))
	if role == "" {
		return RoleUser
	}
	return role
}

func isAllowedRole(role string) bool {
	return role == RoleUser || role == RoleAdmin
}

type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
