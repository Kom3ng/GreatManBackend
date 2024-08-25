package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"greatmanbackend/common"
	"greatmanbackend/controllers"
	"os"
)

var db *gorm.DB

func main() {
	common.InitDB()

	r := gin.Default()

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return
	}

	v0 := r.Group("/v0")
	auth := gin.BasicAuth(gin.Accounts{
		os.Getenv("USERNAME"): os.Getenv("PASSWORD"),
	})
	{
		v0.GET("/man/:id", controllers.GetGreatMan)
		v0.GET("/man/:id/talks", controllers.GetTalks)
		v0.GET("/talk/:id/content", controllers.GetTalkDetail)

		v0.POST("/man/:id", controllers.CreatNewMan).Use(auth)
		v0.PUT("/man/:id").Use(auth)
		v0.DELETE("/man/:id").Use(auth)

		v0.POST("/man/:id/talk").Use(auth)

		v0.PUT("/talk/:id").Use(auth)
		v0.DELETE("/talk/:id").Use(auth)
	}

	_ = r.Run("localhost:" + os.Getenv("PORT"))
}
