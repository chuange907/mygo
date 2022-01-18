package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func helloHandler(c *gin.Context)  {
	c.JSON(http.StatusOK,gin.H{
		"message":"Hello qimi!",
	})

}

func setupRouter() *gin.Engine  {
	r := gin.Default()
	r.GET("/hello",helloHandler)
	return r

}