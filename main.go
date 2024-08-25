package main

import (
	"github.com/gin-gonic/gin"
	"greatmanbackend/common"
	"greatmanbackend/controllers"
	"os"
)

func main() {
	common.InitDB()

	r := gin.Default()

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return
	}

	v0 := r.Group("/v0")
	v0Auth := r.Group("/v0")

	auth := gin.BasicAuth(gin.Accounts{
		os.Getenv("USERNAME"): os.Getenv("PASSWORD"),
	})

	v0Auth.Use(auth)
	{
		v0.GET("/men", controllers.GetGreatMen)
		v0.GET("/man/:id", controllers.GetGreatMan)
		v0.GET("/man/:id/talks", controllers.GetTalks)
		v0.GET("/talk/:id/content", controllers.GetTalkDetail)

		v0Auth.POST("/man", controllers.CreatNewMan)
		v0Auth.PUT("/man/:id", controllers.UpdateMan)
		v0Auth.DELETE("/man/:id", controllers.DeleteMan)

		v0Auth.POST("/man/:id/talk")

		v0Auth.PUT("/talk/:id")
		v0Auth.DELETE("/talk/:id")
	}

	_ = r.Run("localhost:" + os.Getenv("PORT"))
}
