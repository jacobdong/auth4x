package auth

import (
	"context"
	"testing"

	"auth4x/internal/config"
	userslogic "auth4x/internal/logic/users"
	"auth4x/internal/svc"
	"auth4x/internal/types"
)

func TestAuthFlow(t *testing.T) {
	cfg := config.Config{
		Auth: config.AuthConfig{
			AccessSecret:  "test-access-secret",
			AccessExpire:  3600,
			RefreshSecret: "test-refresh-secret",
			RefreshExpire: 7200,
		},
	}
	svcCtx := svc.NewServiceContext(cfg)
	ctx := context.Background()

	signUpLogic := NewSignUpLogic(ctx, svcCtx)
	signUpResp, err := signUpLogic.SignUp(&types.SignUpReq{
		Username: "alice",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("sign up failed: %v", err)
	}
	if signUpResp.UserId == "" {
		t.Fatalf("expected user id")
	}

	signInLogic := NewSignInLogic(ctx, svcCtx)
	signInResp, err := signInLogic.SignIn(&types.SignInReq{
		Username: "alice",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("sign in failed: %v", err)
	}
	if signInResp.AccessToken == "" || signInResp.RefreshToken == "" {
		t.Fatalf("expected tokens")
	}

	tokenLogic := NewTokenLogic(ctx, svcCtx)
	refreshResp, err := tokenLogic.Token(&types.TokenReq{
		RefreshToken: signInResp.RefreshToken,
	})
	if err != nil {
		t.Fatalf("refresh token failed: %v", err)
	}
	if refreshResp.AccessToken == "" || refreshResp.RefreshToken == "" {
		t.Fatalf("expected refreshed tokens")
	}

	userMeLogic := userslogic.NewUserMeLogic(ctx, svcCtx)
	userResp, err := userMeLogic.UserMe(&types.UserMeReq{
		Authorization: "Bearer " + refreshResp.AccessToken,
	})
	if err != nil {
		t.Fatalf("user me failed: %v", err)
	}
	if userResp.Username != "alice" {
		t.Fatalf("expected username alice, got %q", userResp.Username)
	}
}
