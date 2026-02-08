package handler

import (
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"life-system-backend/internal/logic"
	"life-system-backend/internal/middleware"
	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			httpx.OkJson(w, types.CommonResp{
				Code:    401,
				Message: err.Error(),
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
