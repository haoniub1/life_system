package handler

import (
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
	"life-system-backend/internal/logic"
	"life-system-backend/internal/middleware"
	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
)

func GetShopItemsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 401, Message: "unauthorized"})
			return
		}

		l := logic.NewShopLogic(svcCtx)
		resp, err := l.GetShopItems(r.Context(), userID)
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

func CreateShopItemHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 401, Message: "unauthorized"})
			return
		}

		var req types.CreateShopItemReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: "invalid request"})
			return
		}

		l := logic.NewShopLogic(svcCtx)
		resp, err := l.CreateShopItem(r.Context(), userID, &req)
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

func UpdateShopItemHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 401, Message: "unauthorized"})
			return
		}

		vars := pathvar.Vars(r)
		itemID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: "invalid item id"})
			return
		}

		var req types.UpdateShopItemReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: "invalid request"})
			return
		}

		l := logic.NewShopLogic(svcCtx)
		resp, err := l.UpdateShopItem(r.Context(), userID, itemID, &req)
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

func DeleteShopItemHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 401, Message: "unauthorized"})
			return
		}

		vars := pathvar.Vars(r)
		itemID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: "invalid item id"})
			return
		}

		l := logic.NewShopLogic(svcCtx)
		if err := l.DeleteShopItem(r.Context(), userID, itemID); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: err.Error()})
			return
		}

		httpx.OkJson(w, types.CommonResp{
			Code:    0,
			Message: "success",
		})
	}
}

func PurchaseItemHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 401, Message: "unauthorized"})
			return
		}

		var req types.PurchaseItemReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: "invalid request"})
			return
		}

		l := logic.NewShopLogic(svcCtx)
		resp, err := l.PurchaseItem(r.Context(), userID, &req)
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: err.Error()})
			return
		}

		httpx.OkJson(w, types.CommonResp{
			Code:    0,
			Message: resp.Message,
			Data:    resp,
		})
	}
}

func GetInventoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 401, Message: "unauthorized"})
			return
		}

		l := logic.NewShopLogic(svcCtx)
		resp, err := l.GetInventory(r.Context(), userID)
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

func UseItemHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 401, Message: "unauthorized"})
			return
		}

		var req types.UseItemReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: "invalid request"})
			return
		}

		l := logic.NewShopLogic(svcCtx)
		resp, err := l.UseItem(r.Context(), userID, &req)
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: err.Error()})
			return
		}

		httpx.OkJson(w, types.CommonResp{
			Code:    0,
			Message: resp.Message,
			Data:    resp,
		})
	}
}

func GetPurchaseHistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 401, Message: "unauthorized"})
			return
		}

		l := logic.NewShopLogic(svcCtx)
		resp, err := l.GetPurchaseHistory(r.Context(), userID)
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
