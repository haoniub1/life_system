package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"life-system-backend/internal/logic"
	"life-system-backend/internal/middleware"
	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
	"life-system-backend/pkg/ratelimit"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := ratelimit.GetClientIP(r)

		// Check registration rate limit
		if !svcCtx.RateLimiter.CheckRegister(clientIP) {
			httpx.OkJson(w, types.CommonResp{
				Code:    429,
				Message: "今日注册次数已达上限，请明天再试",
			})
			return
		}

		var req types.RegisterReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    400,
				Message: "invalid request",
			})
			return
		}

		auth := logic.NewAuthLogic(svcCtx)
		resp, err := auth.Register(r.Context(), &req)
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    400,
				Message: err.Error(),
			})
			return
		}

		// Count successful registration
		svcCtx.RateLimiter.RecordRegister(clientIP)

		// Set JWT cookie (HttpOnly=false for development with Vite proxy)
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    resp.Token,
			Path:     "/",
			HttpOnly: false, // Changed to false for localhost development
			MaxAge:   int(svcCtx.Config.Auth.Expire),
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
		})

		httpx.OkJson(w, types.CommonResp{
			Code:    0,
			Message: "success",
			Data:    resp,
		})
	}
}

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := ratelimit.GetClientIP(r)

		// Check login rate limit
		if !svcCtx.RateLimiter.CheckLogin(clientIP) {
			httpx.OkJson(w, types.CommonResp{
				Code:    429,
				Message: "今日登录失败次数过多，账号已锁定，请明天再试",
			})
			return
		}

		var req types.LoginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    400,
				Message: "invalid request",
			})
			return
		}

		auth := logic.NewAuthLogic(svcCtx)
		resp, err := auth.Login(r.Context(), &req)
		if err != nil {
			svcCtx.RateLimiter.RecordLoginFailure(clientIP)
			remaining := svcCtx.RateLimiter.LoginFailuresRemaining(clientIP)
			httpx.OkJson(w, types.CommonResp{
				Code:    401,
				Message: fmt.Sprintf("用户名或密码错误（今日剩余尝试次数：%d）", remaining),
			})
			return
		}

		// Set JWT cookie (HttpOnly=false for development with Vite proxy)
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    resp.Token,
			Path:     "/",
			HttpOnly: false, // Changed to false for localhost development
			MaxAge:   int(svcCtx.Config.Auth.Expire),
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
		})

		httpx.OkJson(w, types.CommonResp{
			Code:    0,
			Message: "success",
			Data:    resp,
		})
	}
}

func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clear token cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
			SameSite: http.SameSiteLaxMode,
		})

		httpx.OkJson(w, types.CommonResp{
			Code:    0,
			Message: "success",
		})
	}
}

func GetMeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    401,
				Message: "unauthorized",
			})
			return
		}

		auth := logic.NewAuthLogic(svcCtx)
		user, err := auth.GetMe(r.Context(), userID)
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    400,
				Message: err.Error(),
			})
			return
		}

		httpx.OkJson(w, types.CommonResp{
			Code:    0,
			Message: "success",
			Data:    user,
		})
	}
}
