package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

type UserInfo struct {
	Username string
	Password string
}

var usermap = make(map[string]UserInfo)

var (
	datefile *os.File
)

func main(){
	Init()
	ReloadData()

	r := gin.Default()
//注册账号
	r.POST("/register", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		user :=UserInfo{
			Username: username,
			Password: password,
		}
			_,ok := usermap[username]
			if ok{
				c.JSON(403,"账号已存在")
				c.Abort()
			}else {
				WriteData(user.Username,user)
				c.SetCookie("userInfo",username,3060,"/","",false,true)
				c.JSON(200,gin.H{
					"系统通知": "恭喜您，注册成功",
				})
			}
	})

//登录账号
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		_,ok := usermap[username]
		if ok == false{
			c.JSON(403,gin.H{
				"警告":"您的账号不存在",
			})
			return
		}
		loadUser := usermap[username]
		if username == loadUser.Username && password == loadUser.Password{
			c.JSON(200,"恭喜您，登录成功")
			c.SetCookie("username",username,3060,"/","",false,true)
		}else{
			c.JSON(403,gin.H{
				"警告":"您的输入有误，请重新输入",
			})
			c.Abort()
		}
	})

// /hello网页
	r.GET("/hello",auth,func(c *gin.Context) {
		c.JSON(200,"欢迎访问小霖的小站！")
	})

	_ = r.Run()
}

//中间件判断是否登录
func auth (c *gin.Context){
	path := c.FullPath()
	v,err := c.Cookie("userInfo")
	if err != nil{
		c.JSON(200,gin.H{
			"code":path,
			"message":"游客您好",
		})
	}else{
		c.Set("userInfo",v)
		c.JSON(200,gin.H{
			"code":path,
			"message":v+"您好",
		})
	}
}

//打开文件
func Init(){
	file,err := os.OpenFile("./userInfo.data",os.O_CREATE|os.O_RDWR|os.O_APPEND,0666)
	datefile = file
	if err != nil{
		fmt.Println(err)
		return
	}
	_ = datefile.Close()
}

//读取数据
func ReloadData(){
	reader := bufio.NewReader(datefile)

	var u UserInfo

	for {
		str, err2 := reader.ReadString('\n')
		if err2 == io.EOF{
			break
		}
		if err2 != nil{
			fmt.Println(err2)
			return
		}
			err3 := json.Unmarshal([]byte(str), &u)
			if err3 != nil {
				fmt.Println(err2)
				return
		}
		usermap[u.Username] = u
	}
}

func WriteData(name string,user UserInfo){
	usermap[name] = user
	jsonByte,_ := json.Marshal(usermap[name])
	jsonString := string(jsonByte)
	datefile.WriteString(jsonString)
}

