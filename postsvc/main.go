package main

import (
	"fmt"
	"postsvc/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println()
	router := gin.Default()

	router.GET("/posts", handlers.GetPosts)
	router.GET("posts/:id", handlers.GetPost)
	router.PUT("posts/:id", handlers.UpdatePost)
	router.POST("posts", handlers.CreatePost)
	router.DELETE("posts/:id", handlers.DeletePost)

	router.Run(":8080")
}
