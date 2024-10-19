package api

import (
	"SkinRest/internal/database"
	"SkinRest/internal/middleware"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewAppCtx(db *sql.DB, logger *zap.Logger) *database.AppContext {
	return &database.AppContext{
		DB:     db,
		Logger: logger,
	}
}

// Set app-context middleware
func ContextMiddleware(appCtx *database.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("appCtx", appCtx)
		c.Next()
	}
}

func NewRouter(logger *zap.Logger, DB *sql.DB) *gin.Engine {
	appCtx := NewAppCtx(DB, logger) // initialize AppContext

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(ContextMiddleware(appCtx)) // use AppContext for all handlers

	v1 := r.Group("/api/v1")
	v1.GET("/", HealthCheck)

	auth := v1.Group("/user")

	auth.POST("/register", RegisterHandler)
	auth.POST("/login", LoginHandler)
	auth.GET("/me", middleware.ApiKeyAuth(), middleware.ValidateAuthToken(), AboutMe)

	skins := v1.Group("/skins", middleware.ApiKeyAuth())

	skins.POST("/add", AddNewSkin)
	skins.GET("/", GetSkinsCollection)
	skins.GET("/:id", GetSkin)
	skins.DELETE("/:id", DeleteSkin)

	r.SetTrustedProxies(nil)

	return r
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
