package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"life-system-backend/internal/middleware"
	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
)

type AuthLogic struct {
	svcCtx *svc.ServiceContext
}

func NewAuthLogic(svcCtx *svc.ServiceContext) *AuthLogic {
	return &AuthLogic{
		svcCtx: svcCtx,
	}
}

func (l *AuthLogic) Register(ctx context.Context, req *types.RegisterReq) (*types.AuthResp, error) {
	// Validate input
	if req.Username == "" || req.Password == "" {
		return nil, fmt.Errorf("username and password are required")
	}

	// Check if username exists
	existingUser, err := l.svcCtx.UserModel.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	userID, err := l.svcCtx.UserModel.Create(req.Username, string(hash))
	if err != nil {
		return nil, err
	}

	// Create default character stats
	if err := l.svcCtx.CharacterModel.Create(userID); err != nil {
		return nil, err
	}

	// Generate JWT
	token, err := l.generateToken(userID)
	if err != nil {
		return nil, err
	}

	// Get user info
	user, err := l.svcCtx.UserModel.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &types.AuthResp{
		Token: token,
		User: types.UserInfo{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Avatar:      user.Avatar,
			TgChatID:    user.TgChatID,
			TgUsername:  user.TgUsername,
		},
	}, nil
}

func (l *AuthLogic) Login(ctx context.Context, req *types.LoginReq) (*types.AuthResp, error) {
	// Find user by username
	user, err := l.svcCtx.UserModel.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	// Generate JWT
	token, err := l.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &types.AuthResp{
		Token: token,
		User: types.UserInfo{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Avatar:      user.Avatar,
			TgChatID:    user.TgChatID,
			TgUsername:  user.TgUsername,
		},
	}, nil
}

func (l *AuthLogic) GetMe(ctx context.Context, userID int64) (*types.UserInfo, error) {
	user, err := l.svcCtx.UserModel.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return &types.UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Avatar:      user.Avatar,
		TgChatID:    user.TgChatID,
		TgUsername:  user.TgUsername,
	}, nil
}

func (l *AuthLogic) generateToken(userID int64) (string, error) {
	claims := &middleware.CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(l.svcCtx.Config.Auth.Expire) * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(l.svcCtx.Config.Auth.Secret))
}
