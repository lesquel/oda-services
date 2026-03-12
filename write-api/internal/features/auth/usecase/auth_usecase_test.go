package usecase_test

import (
	"errors"
	"testing"

	"github.com/lesquel/oda-shared/domain"
	"github.com/lesquel/oda-write-api/internal/features/auth/usecase"
)

// ── Mock repositories ─────────────────────────────────────────────────────────

type mockUserRepo struct {
	users map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*domain.User)}
}

func (m *mockUserRepo) Create(u *domain.User) error {
	if _, exists := m.users[u.Email]; exists {
		return errors.New("email already exists")
	}
	m.users[u.Email] = u
	m.users[u.Username] = u
	return nil
}

func (m *mockUserRepo) FindByID(id string) (*domain.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepo) FindByEmail(email string) (*domain.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepo) FindByUsername(username string) (*domain.User, error) {
	if u, ok := m.users[username]; ok {
		return u, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepo) Update(u *domain.User) error {
	m.users[u.Email] = u
	return nil
}

func (m *mockUserRepo) Delete(id string) error { return nil }

func (m *mockUserRepo) Search(query string, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}

type mockTokenRepo struct {
	tokens map[string]*domain.RefreshToken
}

func newMockTokenRepo() *mockTokenRepo {
	return &mockTokenRepo{tokens: make(map[string]*domain.RefreshToken)}
}

func (m *mockTokenRepo) Create(t *domain.RefreshToken) error {
	m.tokens[t.Token] = t
	return nil
}

func (m *mockTokenRepo) FindByToken(token string) (*domain.RefreshToken, error) {
	if t, ok := m.tokens[token]; ok {
		return t, nil
	}
	return nil, errors.New("token not found")
}

func (m *mockTokenRepo) DeleteByToken(token string) error {
	delete(m.tokens, token)
	return nil
}

func (m *mockTokenRepo) DeleteByUserID(userID string) error {
	for k, t := range m.tokens {
		if t.UserID == userID {
			delete(m.tokens, k)
		}
	}
	return nil
}

// ── Tests ─────────────────────────────────────────────────────────────────────

const testJWTSecret = "test-jwt-secret-key"

func TestRegister_Success(t *testing.T) {
	uc := usecase.NewAuthUseCase(newMockUserRepo(), newMockTokenRepo(), testJWTSecret)
	req := &domain.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	resp, err := uc.Register(req)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if resp.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
	if resp.User == nil {
		t.Error("expected non-nil user in response")
	}
	if resp.User.Email != req.Email {
		t.Errorf("expected email %s, got %s", req.Email, resp.User.Email)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	userRepo := newMockUserRepo()
	uc := usecase.NewAuthUseCase(userRepo, newMockTokenRepo(), testJWTSecret)
	req := &domain.RegisterRequest{
		Username: "user1",
		Email:    "dup@example.com",
		Password: "password123",
	}
	_, err := uc.Register(req)
	if err != nil {
		t.Fatalf("first Register failed: %v", err)
	}
	_, err = uc.Register(req)
	if err == nil {
		t.Error("expected error for duplicate registration, got nil")
	}
}

func TestLogin_Success(t *testing.T) {
	uc := usecase.NewAuthUseCase(newMockUserRepo(), newMockTokenRepo(), testJWTSecret)
	registerReq := &domain.RegisterRequest{
		Username: "loginuser",
		Email:    "login@example.com",
		Password: "mypassword",
	}
	if _, err := uc.Register(registerReq); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	resp, err := uc.Login(&domain.LoginRequest{
		Email:    "login@example.com",
		Password: "mypassword",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	uc := usecase.NewAuthUseCase(newMockUserRepo(), newMockTokenRepo(), testJWTSecret)
	_, _ = uc.Register(&domain.RegisterRequest{
		Username: "user2",
		Email:    "user2@example.com",
		Password: "correctpassword",
	})

	_, err := uc.Login(&domain.LoginRequest{
		Email:    "user2@example.com",
		Password: "wrongpassword",
	})
	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}
}

func TestLogin_NonExistentUser(t *testing.T) {
	uc := usecase.NewAuthUseCase(newMockUserRepo(), newMockTokenRepo(), testJWTSecret)
	_, err := uc.Login(&domain.LoginRequest{
		Email:    "nobody@example.com",
		Password: "password",
	})
	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}
