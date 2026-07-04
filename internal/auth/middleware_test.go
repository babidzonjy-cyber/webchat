package auth

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	token, err := GenerateToken(42)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := UserIDFromContext(r.Context())
		w.Write([]byte(strconv.Itoa(id)))
	})

	wrapped := AuthMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r := httptest.NewRecorder()

	wrapped.ServeHTTP(r, req)

	if r.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", r.Code)
	}
}

func TestAuthMiddleware_QueryParam(t *testing.T) {
	token, err := GenerateToken(42)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := UserIDFromContext(r.Context())
		w.Write([]byte(strconv.Itoa(id)))
	})

	wrapped := AuthMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test?token="+token, nil)
	req.Header.Set("Authorization", "Bearer")
	r := httptest.NewRecorder()

	wrapped.ServeHTTP(r, req)

	if r.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", r.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	token, err := GenerateToken(42)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := UserIDFromContext(r.Context())
		w.Write([]byte(strconv.Itoa(id)))
	})

	wrapped := AuthMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer garbage"+token)
	r := httptest.NewRecorder()

	wrapped.ServeHTTP(r, req)

	if r.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", r.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	token, err := generateExpiredToken(42)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := UserIDFromContext(r.Context())
		w.Write([]byte(strconv.Itoa(id)))
	})

	wrapped := AuthMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r := httptest.NewRecorder()

	wrapped.ServeHTTP(r, req)

	if r.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", r.Code)
	}
}

func generateExpiredToken(userID int) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
