package jwt_test

import (
	"testing"
	"time"

	"hosim-go/pkg/jwt"
)

func TestJWT_GenerateAndValidate(t *testing.T) {
	secret := "test-secret-key-12345"
	userID := "usr-123"
	username := "john_doe"
	role := "DOCTOR"
	permissions := []string{"patient:read", "patient:create"}
	duration := 15 * time.Minute

	tokenString, err := jwt.GenerateAccessToken(userID, username, role, permissions, secret, duration)
	if err != nil {
		t.Fatalf("GenerateAccessToken gagal: %v", err)
	}

	if tokenString == "" {
		t.Fatal("tokenString tidak boleh kosong")
	}

	claims, err := jwt.ValidateToken(tokenString, secret)
	if err != nil {
		t.Fatalf("ValidateToken gagal: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID tidak cocok. Diharapkan: %s, Didapat: %s", userID, claims.UserID)
	}
	if claims.Username != username {
		t.Errorf("Username tidak cocok. Diharapkan: %s, Didapat: %s", username, claims.Username)
	}
	if claims.Role != role {
		t.Errorf("Role tidak cocok. Diharapkan: %s, Didapat: %s", role, claims.Role)
	}
	if len(claims.Permissions) != 2 || claims.Permissions[0] != "patient:read" {
		t.Errorf("Permissions tidak cocok. Didapat: %v", claims.Permissions)
	}
}

func TestJWT_ExpiredToken(t *testing.T) {
	secret := "test-secret-key-12345"
	userID := "usr-expired"
	username := "expired_user"
	role := "STAFF"
	permissions := []string{"department:read"}
	duration := -1 * time.Minute // Token sudah kedaluwarsa

	tokenString, err := jwt.GenerateAccessToken(userID, username, role, permissions, secret, duration)
	if err != nil {
		t.Fatalf("GenerateAccessToken gagal: %v", err)
	}

	_, err = jwt.ValidateToken(tokenString, secret)
	if err == nil {
		t.Fatal("ValidateToken seharusnya gagal karena token expired, tapi tidak ada error")
	}
	if err != jwt.ErrExpiredToken {
		t.Errorf("Diharapkan ErrExpiredToken, didapat: %v", err)
	}
}

func TestJWT_InvalidSecret(t *testing.T) {
	secret := "test-secret-key-12345"
	wrongSecret := "wrong-secret-key"

	tokenString, err := jwt.GenerateAccessToken("usr-1", "admin", "ADMIN", []string{"*"}, secret, 10*time.Minute)
	if err != nil {
		t.Fatalf("GenerateAccessToken gagal: %v", err)
	}

	_, err = jwt.ValidateToken(tokenString, wrongSecret)
	if err == nil {
		t.Fatal("ValidateToken dengan secret salah seharusnya gagal")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	token1, err := jwt.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken token1 gagal: %v", err)
	}

	token2, err := jwt.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("GenerateRefreshToken token2 gagal: %v", err)
	}

	if token1 == "" || token2 == "" {
		t.Fatal("Refresh token tidak boleh kosong")
	}

	if token1 == token2 {
		t.Fatal("Dua refresh token yang digenerate tidak boleh identik")
	}

	// 32 bytes hex encoded = 64 characters
	if len(token1) != 64 {
		t.Errorf("Panjang token seharusnya 64 karakter hex, didapat: %d", len(token1))
	}
}
