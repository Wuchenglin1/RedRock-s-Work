package main

import "github.com/gin-gonic/gin"

func main(){
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200,gin.H{
			"username":"张三",
			"age":"18",
			"sex":"male",
		})
	})

	_ = r.Run()
}
