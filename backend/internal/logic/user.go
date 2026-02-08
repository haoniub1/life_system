package logic

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
)

type UserLogic struct {
	svcCtx *svc.ServiceContext
}

func NewUserLogic(svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		svcCtx: svcCtx,
	}
}

func (l *UserLogic) UpdateProfile(ctx context.Context, userID int64, req *types.UpdateProfileReq) (*types.UserInfo, error) {
	user, err := l.svcCtx.UserModel.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("用户不存在")
	}

	displayName := req.DisplayName
	if displayName == "" {
		displayName = user.DisplayName
	}

	avatar := req.Avatar
	if avatar == "" {
		avatar = user.Avatar
	}

	if err := l.svcCtx.UserModel.UpdateProfile(userID, displayName, avatar); err != nil {
		return nil, err
	}

	return &types.UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: displayName,
		Avatar:      avatar,
		TgChatID:    user.TgChatID,
		TgUsername:  user.TgUsername,
	}, nil
}

func (l *UserLogic) ChangePassword(ctx context.Context, userID int64, req *types.ChangePasswordReq) error {
	if req.OldPassword == "" || req.NewPassword == "" {
		return fmt.Errorf("旧密码和新密码不能为空")
	}

	if len(req.NewPassword) < 6 {
		return fmt.Errorf("新密码长度不能少于6位")
	}

	user, err := l.svcCtx.UserModel.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("用户不存在")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return fmt.Errorf("旧密码不正确")
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return l.svcCtx.UserModel.UpdatePassword(userID, string(hash))
}
