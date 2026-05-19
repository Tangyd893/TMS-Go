package crypto

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("my-secret-password")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword() returned empty hash")
	}
	if hash == "my-secret-password" {
		t.Fatal("HashPassword() returned plaintext")
	}
}

func TestCheckPasswordMatch(t *testing.T) {
	password := "correct-password"
	hash, _ := HashPassword(password)
	if !CheckPassword(password, hash) {
		t.Fatal("CheckPassword() should return true for matching password")
	}
}

func TestCheckPasswordMismatch(t *testing.T) {
	hash, _ := HashPassword("correct-password")
	if CheckPassword("wrong-password", hash) {
		t.Fatal("CheckPassword() should return false for wrong password")
	}
}

func TestCheckPasswordEmpty(t *testing.T) {
	hash, _ := HashPassword("")
	if !CheckPassword("", hash) {
		t.Fatal("CheckPassword() should work with empty password")
	}
}

func TestHashPasswordProducesUniqueHashes(t *testing.T) {
	hash1, _ := HashPassword("same-password")
	hash2, _ := HashPassword("same-password")
	if hash1 == hash2 {
		t.Fatal("HashPassword() should produce different hashes due to salt")
	}
}
