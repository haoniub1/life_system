package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"life-system-backend/internal/logic"
	"life-system-backend/internal/middleware"
	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
)

func GenerateBindCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    401,
				Message: "unauthorized",
			})
			return
		}

		tg := logic.NewTelegramLogic(svcCtx)
		resp, err := tg.GenerateBindCode(r.Context(), userID)
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
			Data:    resp,
		})
	}
}

func GetTelegramStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    401,
				Message: "unauthorized",
			})
			return
		}

		tg := logic.NewTelegramLogic(svcCtx)
		resp, err := tg.GetStatus(r.Context(), userID)
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
			Data:    resp,
		})
	}
}

func UnbindTelegramHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    401,
				Message: "unauthorized",
			})
			return
		}

		tg := logic.NewTelegramLogic(svcCtx)
		err = tg.Unbind(r.Context(), userID)
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
		})
	}
}
