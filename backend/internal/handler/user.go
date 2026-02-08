package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"life-system-backend/internal/logic"
	"life-system-backend/internal/middleware"
	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
)

func UpdateProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 401, Message: "unauthorized"})
			return
		}

		var req types.UpdateProfileReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: "invalid request"})
			return
		}

		l := logic.NewUserLogic(svcCtx)
		resp, err := l.UpdateProfile(r.Context(), userID, &req)
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: err.Error()})
			return
		}

		httpx.OkJson(w, types.CommonResp{
			Code:    0,
			Message: "success",
			Data:    resp,
		})
	}
}

func ChangePasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 401, Message: "unauthorized"})
			return
		}

		var req types.ChangePasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: "invalid request"})
			return
		}

		l := logic.NewUserLogic(svcCtx)
		if err := l.ChangePassword(r.Context(), userID, &req); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: err.Error()})
			return
		}

		httpx.OkJson(w, types.CommonResp{
			Code:    0,
			Message: "密码修改成功",
		})
	}
}
