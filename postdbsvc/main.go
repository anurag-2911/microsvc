package main

import (
	"postdbsvc/db"
	"postdbsvc/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	defer db.CloseDB()

	router:=gin.Default()

	router.GET("/posts",handlers.GetPosts)
	router.POST("/posts",handlers.CreatePost)

	router.Run(":9090")

}