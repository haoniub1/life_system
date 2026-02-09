package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
	"life-system-backend/internal/middleware"
	"life-system-backend/internal/svc"
)

func RegisterRoutes(server *rest.Server, svcCtx *svc.ServiceContext) {
	// Public routes (no auth)
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  "POST",
				Path:    "/api/auth/register",
				Handler: RegisterHandler(svcCtx),
			},
			{
				Method:  "POST",
				Path:    "/api/auth/login",
				Handler: LoginHandler(svcCtx),
			},
		},
	)

	// Static file serving for uploads
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  "GET",
				Path:    "/uploads/:file",
				Handler: http.StripPrefix("/uploads/", http.FileServer(http.Dir("data/uploads"))).ServeHTTP,
			},
		},
	)

	// Protected routes (with auth)
	authMiddleware := middleware.AuthMiddleware(svcCtx.Config.Auth.Secret)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  "POST",
				Path:    "/api/auth/logout",
				Handler: authMiddleware(LogoutHandler(svcCtx)),
			},
			{
				Method:  "GET",
				Path:    "/api/auth/me",
				Handler: authMiddleware(GetMeHandler(svcCtx)),
			},
			// User profile
			{
				Method:  "PUT",
				Path:    "/api/user/profile",
				Handler: authMiddleware(UpdateProfileHandler(svcCtx)),
			},
			{
				Method:  "PUT",
				Path:    "/api/user/password",
				Handler: authMiddleware(ChangePasswordHandler(svcCtx)),
			},
			// Upload
			{
				Method:  "POST",
				Path:    "/api/upload",
				Handler: authMiddleware(UploadHandler(svcCtx)),
			},
			// Character
			{
				Method:  "GET",
				Path:    "/api/character",
				Handler: authMiddleware(GetCharacterHandler(svcCtx)),
			},
			{
				Method:  "PUT",
				Path:    "/api/character",
				Handler: authMiddleware(UpdateCharacterHandler(svcCtx)),
			},
			// Tasks
			{
				Method:  "GET",
				Path:    "/api/tasks",
				Handler: authMiddleware(ListTasksHandler(svcCtx)),
			},
			{
				Method:  "POST",
				Path:    "/api/tasks",
				Handler: authMiddleware(CreateTaskHandler(svcCtx)),
			},
			{
				Method:  "PUT",
				Path:    "/api/tasks/:id",
				Handler: authMiddleware(UpdateTaskHandler(svcCtx)),
			},
			{
				Method:  "POST",
				Path:    "/api/tasks/complete/:id",
				Handler: authMiddleware(CompleteTaskHandler(svcCtx)),
			},
			{
				Method:  "DELETE",
				Path:    "/api/tasks/:id",
				Handler: authMiddleware(DeleteTaskHandler(svcCtx)),
			},
			// Telegram
			{
				Method:  "POST",
				Path:    "/api/telegram/bindcode",
				Handler: authMiddleware(GenerateBindCodeHandler(svcCtx)),
			},
			{
				Method:  "GET",
				Path:    "/api/telegram/status",
				Handler: authMiddleware(GetTelegramStatusHandler(svcCtx)),
			},
			{
				Method:  "DELETE",
				Path:    "/api/telegram/unbind",
				Handler: authMiddleware(UnbindTelegramHandler(svcCtx)),
			},
			// Bark Push Notification
			{
				Method:  "PUT",
				Path:    "/api/bark/key",
				Handler: authMiddleware(SetBarkKeyHandler(svcCtx)),
			},
			{
				Method:  "GET",
				Path:    "/api/bark/status",
				Handler: authMiddleware(GetBarkStatusHandler(svcCtx)),
			},
			{
				Method:  "POST",
				Path:    "/api/bark/test",
				Handler: authMiddleware(TestBarkHandler(svcCtx)),
			},
			{
				Method:  "DELETE",
				Path:    "/api/bark/key",
				Handler: authMiddleware(DeleteBarkKeyHandler(svcCtx)),
			},
			// Timeline
			{
				Method:  "GET",
				Path:    "/api/timeline",
				Handler: authMiddleware(GetTimelineHandler(svcCtx)),
			},
			// Shop
			{
				Method:  "GET",
				Path:    "/api/shop/items",
				Handler: authMiddleware(GetShopItemsHandler(svcCtx)),
			},
			{
				Method:  "POST",
				Path:    "/api/shop/items",
				Handler: authMiddleware(CreateShopItemHandler(svcCtx)),
			},
			{
				Method:  "PUT",
				Path:    "/api/shop/items/:id",
				Handler: authMiddleware(UpdateShopItemHandler(svcCtx)),
			},
			{
				Method:  "DELETE",
				Path:    "/api/shop/items/:id",
				Handler: authMiddleware(DeleteShopItemHandler(svcCtx)),
			},
			{
				Method:  "POST",
				Path:    "/api/shop/purchase",
				Handler: authMiddleware(PurchaseItemHandler(svcCtx)),
			},
			{
				Method:  "GET",
				Path:    "/api/shop/inventory",
				Handler: authMiddleware(GetInventoryHandler(svcCtx)),
			},
			{
				Method:  "POST",
				Path:    "/api/shop/use",
				Handler: authMiddleware(UseItemHandler(svcCtx)),
			},
			{
				Method:  "GET",
				Path:    "/api/shop/history",
				Handler: authMiddleware(GetPurchaseHistoryHandler(svcCtx)),
			},
		},
	)
}
