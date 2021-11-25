package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main (){
	auth := func(c *gin.Context){
		v,err := c.Cookie("username")
		if err != nil{
			c.JSON(403,gin.H{
				"错误":"访问失败，您还没有登录！",
			})
			c.Abort()
		}else {
			c.Set("cookie",v)
		}
	}
	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		if username == "1336636739" && password == "2021214174"{
			c.SetCookie("username",username,3060,"/","",false,true)
			//c.SetCookie("password",password,3060,"/","",false,true)
			c.JSON(200,gin.H{
				"恭喜您":"登录成功",
			})
			//c.Redirect(http.StatusMovedPermanently,"http://localhost/hello")
		}else{
			c.JSON(200,gin.H{
				"警告":"您输入的账号或密码有误",
			})
		}
		fmt.Println("账号:",username,"密码:",password)
	})

	r.GET("/hello",auth,func(c *gin.Context) {
		cookie,_ := c.Get("cookie")
		str := cookie.(string)
		c.JSON(200,"Hello  " + str)
	})

	_ = r.Run()
}
