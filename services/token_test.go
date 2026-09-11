package services

import "testing"

func TestAccessTokenUsesRuntimeSecret(t *testing.T) {
	t.Setenv("TOKEN_SECRET", "test-secret")
	t.Setenv("TOKEN_VALIDITY_MINS", "15")

	token, err := GenerateAccessToken(42, "user@example.com")
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	claims, err := ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken returned error: %v", err)
	}
	if claims["user_id"] != float64(42) {
		t.Fatalf("expected user_id claim 42, got %v", claims["user_id"])
	}
}

func TestAccessTokenRequiresSecret(t *testing.T) {
	t.Setenv("TOKEN_SECRET", "")

	if _, err := GenerateAccessToken(42, "user@example.com"); err == nil {
		t.Fatal("expected token generation to fail without TOKEN_SECRET")
	}
}
