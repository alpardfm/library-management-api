package main

import (
	"net/http"

	"github.com/alpardfm/library-management-api/configs"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerSystemRoutes(router *gin.Engine, db *gorm.DB, cfg *configs.Config) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"app":     cfg.AppName,
			"version": cfg.AppVersion,
			"env":     cfg.AppEnv,
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"error":  "database connection unavailable",
			})
			return
		}

		if err := sqlDB.PingContext(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"error":  "database ping failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})
}
