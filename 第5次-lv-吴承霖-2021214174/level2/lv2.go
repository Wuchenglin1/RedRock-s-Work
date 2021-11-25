package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var usermap = make(map[string]UserInfo)

var (
	datefile *os.File
)

func main(){
	//先读取一次数据，再将数据文件关闭，之后数据文件随用随关
	Init()
	ReloadData()
	_ = datefile.Close()

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
			}
			if len(username)<6||len(password)<6{
				c.JSON(403,"您的输入有误，请输入长度大于6位的账号和密码")
			}else{
				Init()
				WriteData(user.Username,user)
				_ = datefile.Close()
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
		//登录时比对密码信息
		err := bcrypt.CompareHashAndPassword([]byte(usermap[username].Password),[]byte(password))
		if username == loadUser.Username && err == nil{
			c.SetCookie("username",username,3060,"/","",false,true)
			c.JSON(200,"恭喜您，登录成功")
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
		//读取密码时，直接将读取经过bcrypt.GenerateFromPassword处理的密码值
		usermap[u.Username] = u
	}
}

func WriteData(name string,user UserInfo){
	//先将密码加密处理,再把加密过后的密码传递给user.Password
	hash,err :=bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)
	if err != nil{
		fmt.Println(err)
		return
	}
	user.Password = string(hash)

	//再将密码写进数据文件中
	usermap[name] = user
	jsonByte,_ := json.Marshal(usermap[name])
	jsonString := string(jsonByte)
	_,_ = datefile.WriteString(jsonString+"\r\n")
}

