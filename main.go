package main

import (
	"RestGo/config"
	"RestGo/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	db := config.DBConnect()
	UserController := &controllers.UserDB{DB: db}
	SocialMediaController := &controllers.SocialMediaDB{DB: db}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello Golang",
		})
	})

	router.POST("/users/register", UserController.Register)
	router.POST("/users/login", UserController.Login)

	socialMediaRouter := router.Group("/socialmedias")
	{
		socialMediaRouter.POST("/", SocialMediaController.CreateSocialMedia)
	}

	router.Run(":3000")
}
