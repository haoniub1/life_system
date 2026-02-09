package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"life-system-backend/internal/middleware"
	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
)

// SetBarkKeyHandler sets the user's Bark push notification key
func SetBarkKeyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 401, Message: "unauthorized"})
			return
		}

		var req types.SetBarkKeyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: "invalid request"})
			return
		}

		// Update bark key
		if err := svcCtx.UserModel.UpdateBarkKey(userID, req.BarkKey); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 500, Message: "failed to update bark key"})
			return
		}

		httpx.OkJson(w, types.CommonResp{
			Code:    0,
			Message: "Bark key 设置成功",
		})
	}
}

// GetBarkStatusHandler returns the user's Bark configuration status
func GetBarkStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 401, Message: "unauthorized"})
			return
		}

		barkKey, err := svcCtx.UserModel.GetBarkKey(userID)
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 500, Message: "failed to get bark status"})
			return
		}

		// Mask the bark key for security (show first 8 chars + ***)
		maskedKey := ""
		if barkKey != "" {
			if len(barkKey) > 8 {
				maskedKey = barkKey[:8] + "***"
			} else {
				maskedKey = barkKey
			}
		}

		resp := types.BarkStatusResp{
			Enabled: svcCtx.Config.Bark.Enabled && barkKey != "",
			BarkKey: maskedKey,
		}

		httpx.OkJson(w, types.CommonResp{
			Code:    0,
			Message: "success",
			Data:    resp,
		})
	}
}

// TestBarkHandler sends a test notification to verify Bark setup
func TestBarkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 401, Message: "unauthorized"})
			return
		}

		if svcCtx.BarkClient == nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: "Bark 服务未启用"})
			return
		}

		barkKey, err := svcCtx.UserModel.GetBarkKey(userID)
		if err != nil || barkKey == "" {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: "请先设置 Bark Key"})
			return
		}

		var req types.TestBarkReq
		if err := httpx.Parse(r, &req); err != nil {
			// Use default test message
			req.Title = "🔔 Life System 测试"
			req.Body = "恭喜！Bark 推送配置成功！"
		}

		if req.Title == "" {
			req.Title = "🔔 Life System 测试"
		}
		if req.Body == "" {
			req.Body = "恭喜！Bark 推送配置成功！"
		}

		// Send test notification with alarm style
		if err := svcCtx.BarkClient.PushAlarm(barkKey, req.Title, req.Body); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 500, Message: "推送失败: " + err.Error()})
			return
		}

		httpx.OkJson(w, types.CommonResp{
			Code:    0,
			Message: "测试推送已发送，请检查手机通知",
		})
	}
}

// DeleteBarkKeyHandler removes the user's Bark key
func DeleteBarkKeyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 401, Message: "unauthorized"})
			return
		}

		if err := svcCtx.UserModel.UpdateBarkKey(userID, ""); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 500, Message: "failed to delete bark key"})
			return
		}

		httpx.OkJson(w, types.CommonResp{
			Code:    0,
			Message: "Bark key 已删除",
		})
	}
}
