package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
	"life-system-backend/internal/logic"
	"life-system-backend/internal/middleware"
	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
)

func ListTasksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    401,
				Message: "unauthorized",
			})
			return
		}

		taskType := r.URL.Query().Get("type")
		status := r.URL.Query().Get("status")

		task := logic.NewTaskLogic(svcCtx)
		resp, err := task.ListTasks(r.Context(), userID, taskType, status)
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

func CreateTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    401,
				Message: "unauthorized",
			})
			return
		}

		var req types.CreateTaskReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    400,
				Message: "invalid request",
			})
			return
		}

		task := logic.NewTaskLogic(svcCtx)
		resp, err := task.CreateTask(r.Context(), userID, &req)
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

func UpdateTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    401,
				Message: "unauthorized",
			})
			return
		}

		taskID, err := strconv.ParseInt(pathvar.Vars(r)["id"], 10, 64)
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    400,
				Message: "invalid task id",
			})
			return
		}

		var req types.UpdateTaskReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    400,
				Message: "invalid request",
			})
			return
		}

		task := logic.NewTaskLogic(svcCtx)
		resp, err := task.UpdateTask(r.Context(), userID, taskID, &req)
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

func CompleteTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    401,
				Message: "unauthorized",
			})
			return
		}

		taskID, err := strconv.ParseInt(pathvar.Vars(r)["id"], 10, 64)
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    400,
				Message: "invalid task id",
			})
			return
		}

		task := logic.NewTaskLogic(svcCtx)
		result, err := task.CompleteTask(r.Context(), userID, taskID, "web")
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
			Data:    result,
		})
	}
}

func DeleteTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    401,
				Message: "unauthorized",
			})
			return
		}

		taskID, err := strconv.ParseInt(pathvar.Vars(r)["id"], 10, 64)
		if err != nil {
			httpx.OkJson(w, types.CommonResp{
				Code:    400,
				Message: "invalid task id",
			})
			return
		}

		task := logic.NewTaskLogic(svcCtx)
		err = task.DeleteTask(r.Context(), userID, taskID, "web")
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
