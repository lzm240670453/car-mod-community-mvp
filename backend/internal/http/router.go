package http

import (
	"net/http"

	"retrofit/backend/internal/config"
	"retrofit/backend/internal/http/handler"
	"retrofit/backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouterDeps struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client
	Logger *zap.Logger
}

func NewRouter(deps RouterDeps) *gin.Engine {
	if deps.Config.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger(deps.Logger))

	health := handler.NewHealthHandler(deps.DB, deps.Redis)
	router.GET("/healthz", health.Check)
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
	})

	api := router.Group("/api/v1")
	registerAPIRoutes(api, handler.Deps{DB: deps.DB})

	return router
}

func registerAPIRoutes(api *gin.RouterGroup, deps handler.Deps) {
	auth := handler.NewAuthHandler(deps)
	users := handler.NewUserHandler(deps)
	vehicles := handler.NewVehicleHandler(deps)
	posts := handler.NewPostHandler(deps)
	parts := handler.NewPartHandler(deps)
	intents := handler.NewIntentHandler(deps)
	messages := handler.NewMessageHandler(deps)
	reports := handler.NewReportHandler(deps)
	uploads := handler.NewUploadHandler()
	admin := handler.NewAdminHandler(deps)

	api.POST("/auth/wechat-login", auth.WeChatLogin)
	api.GET("/auth/me", auth.Me)
	api.POST("/auth/logout", auth.Logout)
	api.POST("/auth/bind-phone", auth.BindPhone)

	api.GET("/users/me", users.Me)
	api.PUT("/users/me", users.UpdateMe)
	api.GET("/users/me/garages", users.ListGarages)
	api.POST("/users/me/garages", users.CreateGarage)
	api.PUT("/users/me/garages/:garageId", users.UpdateGarage)
	api.DELETE("/users/me/garages/:garageId", users.DeleteGarage)

	api.GET("/vehicles/brands", vehicles.ListBrands)
	api.GET("/vehicles/brands/:brandId/series", vehicles.ListSeries)
	api.GET("/vehicles/series/:seriesId/models", vehicles.ListModels)
	api.GET("/vehicles/search", vehicles.Search)

	api.GET("/posts", posts.List)
	api.GET("/posts/:postId", posts.Get)
	api.POST("/posts", posts.Create)
	api.PUT("/posts/:postId", posts.Update)
	api.DELETE("/posts/:postId", posts.Delete)
	api.POST("/posts/:postId/like", posts.Like)
	api.DELETE("/posts/:postId/like", posts.Unlike)
	api.POST("/posts/:postId/favorite", posts.Favorite)
	api.DELETE("/posts/:postId/favorite", posts.Unfavorite)
	api.GET("/posts/:postId/comments", posts.ListComments)
	api.POST("/posts/:postId/comments", posts.CreateComment)
	api.DELETE("/comments/:commentId", posts.DeleteComment)

	api.GET("/parts", parts.List)
	api.GET("/parts/:partId", parts.Get)
	api.POST("/parts", parts.Create)
	api.PUT("/parts/:partId", parts.Update)
	api.DELETE("/parts/:partId", parts.Delete)
	api.POST("/parts/:partId/favorite", parts.Favorite)
	api.DELETE("/parts/:partId/favorite", parts.Unfavorite)
	api.GET("/part-categories", parts.ListCategories)

	api.POST("/parts/:partId/intents", intents.Create)
	api.GET("/intents", intents.List)
	api.GET("/intents/:intentId", intents.Get)
	api.POST("/intents/:intentId/close", intents.Close)

	api.GET("/messages", messages.List)
	api.GET("/messages/unread-count", messages.UnreadCount)
	api.POST("/messages/read/:messageId", messages.MarkRead)
	api.POST("/messages/read-all", messages.MarkAllRead)

	api.POST("/reports", reports.Create)
	api.GET("/reports/my", reports.My)

	api.POST("/uploads/signature", uploads.Signature)
	api.POST("/uploads/complete", uploads.Complete)

	adminAPI := api.Group("/admin")
	adminAPI.POST("/auth/login", admin.Login)
	adminAPI.GET("/users", admin.ListUsers)
	adminAPI.PUT("/users/:userId/status", admin.UpdateUserStatus)
	adminAPI.GET("/posts/review", admin.ListPostsForReview)
	adminAPI.GET("/posts/pending-review", admin.ListPendingReviewPosts)
	adminAPI.POST("/posts/:postId/approve", admin.ApprovePost)
	adminAPI.POST("/posts/:postId/hide", admin.HidePost)
	adminAPI.POST("/posts/:postId/restore", admin.RestorePost)
	adminAPI.GET("/parts/review", admin.ListPartsForReview)
	adminAPI.GET("/parts/pending-review", admin.ListPendingReviewParts)
	adminAPI.POST("/parts/:partId/approve", admin.ApprovePart)
	adminAPI.POST("/parts/:partId/hide", admin.HidePart)
	adminAPI.POST("/parts/:partId/restore", admin.RestorePart)
	adminAPI.GET("/reports", admin.ListReports)
	adminAPI.POST("/reports/:reportId/process", admin.ProcessReport)
	adminAPI.GET("/vehicles", admin.ListVehicles)
	adminAPI.POST("/vehicles/brands", admin.CreateBrand)
	adminAPI.POST("/vehicles/series", admin.CreateSeries)
	adminAPI.POST("/vehicles/models", admin.CreateModel)
	adminAPI.GET("/categories", admin.ListCategories)
	adminAPI.POST("/categories", admin.CreateCategory)
	adminAPI.PUT("/categories/:categoryId", admin.UpdateCategory)

}
