package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

var (
	receiveID    string
	Usermb1pwd   string
	Usermb2pwd   string
	mb1, mb2     string
	ID           string
	UserName     string
	UserPassword string
	c            *gin.Context
	db           *sql.DB
	err          error
	whoSend      string
	sendMessage  string
	MessageMap = make(map[string]interface{})
)

const load = "root:root@tcp(127.0.0.1)/test"

func main() {
	Begin()
	r := gin.Default()
	//注册
	r.POST("/register", func(c *gin.Context) {

		inputName := c.PostForm("username")
		inputPassword := c.PostForm("password")
		have := Register(inputName, inputPassword)
		switch have {
		case 0:
			c.JSON(200, "恭喜您注册成功")
		case 1:
			c.Abort()
			c.JSON(403, "您的账号已被注册")
		}

	})

	//登录
	r.POST("/login", func(c *gin.Context) {
		inputName := c.PostForm("username")
		inputPassword := c.PostForm("password")

		have := Login(inputName, inputPassword)
		switch have {
		case 0:
			c.SetCookie("ID", ID, 3600, "/", "", false, true)
			c.SetCookie("user", UserName, 3600, "/", "", false, true)
			c.JSON(200, "恭喜您登录成功！")
			has := CheckMiBao(UserName)
			if has == 2 {
				c.JSON(200, "您还没有设置密保哟~请尽快设置呀！")
			}
		case 1:
			c.JSON(403, "账号不存在！")
		case 2:
			c.JSON(403, "密码错误！")
		}
	})

	//查看和设置密保
	r.POST("/settingmb", auth, func(c *gin.Context) {
		UserName, _ = c.Cookie("user")
		setmbName1 := c.PostForm("mbName1")
		setmbPwd1 := c.PostForm("mbPwd1")
		setmbName2 := c.PostForm("mbName2")
		setmbPwd2 := c.PostForm("mbPwd2")
		have := CheckMiBao(UserName)
		switch have {
		case 0:
			//密保存在
			c.JSON(403, gin.H{
				"warm:": "您已有密保！若想修改密保请访问/changemb",
				"密保1":   mb1,
				"密保2":   mb2,
			})
			c.Abort()
		case 1:
			c.JSON(403, "用户不存在！")
			c.Abort()
		case 2:
			//密保不存在，设置密保
			has := false
			has = AddMiBao(setmbName1, setmbPwd1, setmbName2, setmbPwd2)
			if has {
				//c.SetCookie("mb1", setmbName1, 3600, "/", "", false, true)
				//c.SetCookie("mb1", setmbName2, 3600, "/", "", false, true)
				c.JSON(200, "添加密保成功！")
			} else {
				c.JSON(403, "添加失败！")
				c.Abort()
			}
		}

	})

	//修改密保
	r.POST("/changemb", auth, func(c *gin.Context) {
		UserName, _ = c.Cookie("user")
		ID, _ = c.Cookie("ID")
		/*修改密保请添加key:
		changembname1,changembpassword1,changename2,changembpassword2
		*/
		changembName1 := c.PostForm("changembname1")
		changembPassword1 := c.PostForm("changembpassword1")
		changembName2 := c.PostForm("changembname2")
		changembPassword2 := c.PostForm("changembpassword2")
		have := CheckMiBao(UserName)
		switch have {
		case 0:
			has := Changemb(changembName1, changembPassword1, changembName2, changembPassword2)
			switch has {
			case 0:
				c.JSON(200, "修改成功！")
			case 1:
				c.JSON(403, "您的输入不能为空！")
			case 2:
				c.JSON(403, "您的输入格式有误！请重新输入")
			}
		case 2:
			c.JSON(403, "甘霖娘，都没有设置密保还来修改密保！快设置一个去！")
			c.Abort()
		}

	})

	//改密码
	r.POST("/changepassword", func(c *gin.Context) {
		UserName = c.PostForm("username")
		ChangedPassword := c.PostForm("changePassword")
		have := CheckMiBao(UserName)
		switch have {
		case 0:
			c.JSON(200, gin.H{
				"message": "请选择您验证的密保：",
				"密保1":     mb1,
				"密保2":     mb2,
			})
			mb1pwd := c.PostForm("mbPwd1")
			mb2pwd := c.PostForm("mbPwd2")
			has := ChangePassword(mb1pwd, mb2pwd, ChangedPassword)
			if has == 0 {
				c.JSON(200, "密码修改成功！")
			} else {
				if has == 2 {
					c.JSON(403, "您的输入有误！")
					c.Abort()
				} else {
					c.JSON(403, "密保答案错误！")
				}
			}
		case 1:
			c.JSON(403, gin.H{
				"warm:": "您的账号不存在！",
			})
			c.Abort()
		case 2:
			c.JSON(403, gin.H{
				"message:": "您还没有设置密保不能修改密码喔,请尽快设置密保喔，亲~",
			})
			c.Abort()
		}

	})

	//留言系统
	r.POST("/message", auth,func(c *gin.Context) {
		ID, _ = c.Cookie("ID")
		UserName, _ = c.Cookie("user")
		//需要设置 接收人key:toWhom  消息key:message
		toWhom := c.PostForm("toWhom")
		message := c.PostForm("message")

		//登录此网站之后自动查询是否有人发消息给自己
		has := CheckMessage(ID)
		switch has {
		case 0:
			c.JSON(200, MessageMap)
		case 1:
			c.JSON(403,"您还没有收到信息")
		}
		have := Message(ID, UserName, toWhom, message)
		switch have {
		case 0:
			c.JSON(200, "发送成功！")
		case 1:
			c.Abort()
			c.JSON(403, "发送失败：发送人不能为空！请重新发送!")
		case 2:
			c.JSON(403, "发送失败！")
		case 3:
			c.JSON(4003, "接收人不存在！")
		case 4:
			c.JSON(403, "发送失败！")
		}
	})

	//访问小站
	r.GET("/hello", auth, func(c *gin.Context) {
		path := c.FullPath()
		c.JSON(200, gin.H{
			"code":       path,
			"欢迎来到小霖的小站！": UserName,
		})
	})

	_ = r.Run()
	defer func() {
		_ = db.Close()
	}()
}

func auth(c *gin.Context) {
	UserName, err = c.Cookie("user")
	if err != nil {
		c.JSON(403, "您好游客，请登录")
		c.Abort()
	}
}

// Begin 连接数据库
func Begin() {
	db, err = sql.Open("mysql", load)
	fmt.Println(err)
	CheckErr(err)
}

// Register 注册
func Register(name string, password string) int {
	err = db.QueryRow("select name from user where name = ?", name).Scan(&UserName)
	if err == nil {
		return 1
	}
	hash, err1 := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	CheckErr(err1)
	_, err = db.Exec("insert into user (name,password) values(?,?)", name, string(hash))
	fmt.Println("err2", err)
	return 0
}

// AddMiBao 添加密保
func AddMiBao(mb1, mb2, mb1pwd, mb2pwd string) bool {
	stmt, err1 := db.Prepare("insert into mibao values (?,?,?,?,?)")
	CheckErr(err1)
	_, err = stmt.Exec(ID, mb1, mb1pwd, mb2, mb2pwd)
	if err != nil {
		c.JSON(403, "密保不能为空！")
		c.Abort()
	}
	return true
}

// Changemb 修改密保
func Changemb(mb1, mb1pwd, mb2, mb2pwd string) int {
	if mb1 == "" || mb2 == "" || mb1pwd == "" || mb2pwd == "" {
		return 2
	}
	stmt, err1 := db.Prepare("update mibao set mb1=?,mb1pwd=?,mb2=?,mb2pwd=? where id=?")
	CheckErr(err1)
	_, err = stmt.Exec(mb1, mb1pwd, mb2, mb2pwd, ID)
	if err != nil {
		return 1
	} else {
		return 0
	}
}

// CheckMiBao 检查密保是否存在
func CheckMiBao(name string) int {
	err = db.QueryRow("select id from user where name = ?", name).Scan(&ID)
	if err != nil {
		//用户不存在
		return 1
	}
	err = db.QueryRow("select mb1,mb2 from mibao where id = ?", ID).Scan(&mb1, &mb2)
	if err != nil {
		//密保不存在
		return 2
	} else {
		//密保已存在
		return 0
	}
}

// ChangePassword 修改密码
func ChangePassword(pwd1 string, pwd2 string, changePassword string) int {
	err = db.QueryRow("select mb1pwd,mb2pwd from mibao").Scan(&Usermb1pwd, &Usermb2pwd)
	CheckErr(err)
	if pwd1 == Usermb1pwd || pwd2 == Usermb2pwd {
		stmt, err1 := db.Prepare("update user set password=? where name=?")
		CheckErr(err1)
		_, err = stmt.Exec(changePassword, UserName)
		if changePassword == "" {
			return 2
		}
		if err != nil {
			c.JSON(403, "您的密码输入格式错误！")
		} else {
			return 0
		}
	}
	return 1
}

// Login 登录
func Login(name string, password string) int {
	err = db.QueryRow("select id,name,password from user where name = ?", name).Scan(&ID, &UserName, &UserPassword)
	if err != nil {
		//账号不存在
		return 1
	}
	err = bcrypt.CompareHashAndPassword([]byte(UserPassword), []byte(password))
	if name == UserName && err == nil {
		//密码正确
		return 0
	} else {
		//密码错误
		return 2
	}
}

// Message 发送信息函数
func Message(id, name, toWhom, message string) int {
	if id == "" || name == "" || toWhom == "" {
		//关键信息不能为空
		return 1
	}
	//查询接收人的信息
	err = db.QueryRow("select id from user where name = ?", toWhom).Scan(&receiveID)
	if err != nil {
		//接收人不存在
		return 3
	}
	//发送消息
	stmt, _ := db.Prepare("insert into send values(?,?,?,?)")
	_, err = stmt.Exec(id, name, toWhom, message)
	if err != nil {
		//插入失败
		return 2
	}

	//接收消息同时更新信息库
	stmt, _ = db.Prepare("insert into receive values(?,?,?)")
	_, err = stmt.Exec(receiveID, name, message)
	if err != nil {
		//插入失败
		return 4
	}
	return 0
}

func CheckMessage(id string) int {
	rows,err1 := db.Query("select whoSend,message from receive where id =?",id)
	if err1 != nil{
		return 1
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next(){
		i := 0
		i++
		err = rows.Scan(&whoSend,&sendMessage)
		if err != nil {
			c.JSON(403,err)
			c.Abort()
		}
		MessageMap[whoSend+strconv.Itoa(i)]=sendMessage
	}
	err = rows.Err()
	CheckErr(err)
	return 0
}

// CheckErr 检查err的函数
func CheckErr(err error) {
	if err != nil {
		c.JSON(403, err)
		c.Abort()
	}
}
