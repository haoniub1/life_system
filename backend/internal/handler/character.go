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

func GetCharacterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    401,
				Message: "unauthorized",
			})
			return
		}

		char := logic.NewCharacterLogic(svcCtx)
		resp, err := char.GetCharacter(r.Context(), userID)
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

func UpdateCharacterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    401,
				Message: "unauthorized",
			})
			return
		}

		var req types.UpdateCharacterReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    400,
				Message: "invalid request",
			})
			return
		}

		char := logic.NewCharacterLogic(svcCtx)
		resp, err := char.UpdateCharacter(r.Context(), userID, &req)
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
