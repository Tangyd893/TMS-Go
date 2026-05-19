package jwt

import (
	"testing"
	"time"
)

func newTestManager() *Manager {
	return NewManager("access-secret", "refresh-secret", 2*time.Hour, 168*time.Hour, "tms-test")
}

func TestGenerateAccessToken(t *testing.T) {
	m := newTestManager()
	token, err := m.GenerateAccessToken("user-1", "admin")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}
	if token == "" {
		t.Fatal("GenerateAccessToken() returned empty token")
	}
}

func TestParseAccessToken(t *testing.T) {
	m := newTestManager()
	token, _ := m.GenerateAccessToken("user-1", "admin")
	claims, err := m.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("ParseAccessToken() error = %v", err)
	}
	if claims.UserID != "user-1" {
		t.Fatalf("ParseAccessToken() UserID = %s, want user-1", claims.UserID)
	}
	if claims.Username != "admin" {
		t.Fatalf("ParseAccessToken() Username = %s, want admin", claims.Username)
	}
	if claims.Issuer != "tms-test" {
		t.Fatalf("ParseAccessToken() Issuer = %s, want tms-test", claims.Issuer)
	}
}

func TestParseAccessTokenWrongSecret(t *testing.T) {
	m := newTestManager()
	token, _ := m.GenerateAccessToken("user-1", "admin")
	m2 := NewManager("wrong-secret", "refresh-secret", 2*time.Hour, 168*time.Hour, "tms-test")
	_, err := m2.ParseAccessToken(token)
	if err == nil {
		t.Fatal("ParseAccessToken() should fail with wrong secret")
	}
}

func TestParseAccessTokenExpired(t *testing.T) {
	m := NewManager("access-secret", "refresh-secret", -1*time.Hour, 168*time.Hour, "tms-test")
	token, _ := m.GenerateAccessToken("user-1", "admin")
	_, err := m.ParseAccessToken(token)
	if err == nil {
		t.Fatal("ParseAccessToken() should fail for expired token")
	}
}

func TestParseAccessTokenTampered(t *testing.T) {
	m := newTestManager()
	token, _ := m.GenerateAccessToken("user-1", "admin")
	tampered := token + "x"
	_, err := m.ParseAccessToken(tampered)
	if err == nil {
		t.Fatal("ParseAccessToken() should fail for tampered token")
	}
}

func TestGenerateAndParseRefreshToken(t *testing.T) {
	m := newTestManager()
	token, err := m.GenerateRefreshToken("user-1", "admin")
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}
	claims, err := m.ParseRefreshToken(token)
	if err != nil {
		t.Fatalf("ParseRefreshToken() error = %v", err)
	}
	if claims.UserID != "user-1" {
		t.Fatalf("ParseRefreshToken() UserID = %s, want user-1", claims.UserID)
	}
}

func TestAccessTokenCannotBeParsedAsRefresh(t *testing.T) {
	m := newTestManager()
	accessToken, _ := m.GenerateAccessToken("user-1", "admin")
	_, err := m.ParseRefreshToken(accessToken)
	if err == nil {
		t.Fatal("Access token should not be parseable as refresh token")
	}
}
