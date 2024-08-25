package main

import (
	"github.com/gin-gonic/gin"
	"greatmanbackend/common"
	"greatmanbackend/controllers"
	"net/http"
	"os"
)

func main() {
	common.InitDB()

	r := gin.Default()

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return
	}
	r.Use(Cors)

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

		v0Auth.POST("/man/:id/talk", controllers.NewTalk)

		v0Auth.PUT("/talk/:id", controllers.UpdateTalk)
		v0Auth.DELETE("/talk/:id", controllers.DeleteTalk)
	}

	_ = r.Run("localhost:" + os.Getenv("PORT"))
}

var origin = os.Getenv("CORS_ALLOW_ORIGIN")

func Cors(c *gin.Context) {
	method := c.Request.Method
	c.Header("Access-Control-Allow-Origin", origin)
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	c.Next()
}
