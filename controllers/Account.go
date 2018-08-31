package controllers

import (
	"net/http"
	"strings"
	"sso-v2/models"
	"github.com/gin-gonic/gin"
	"sso-v2/utils"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"github.com/gin-contrib/sessions"
)


const (
	// 插入账户
	INSERTACCOUNT = `INSERT INTO ACCOUNT (username, password, nickname, sex, email) VALUES (?,?,?,?,?)`

	// 根据用户名查找账户
	QUERYACCOUNTBYUSERNAME = `SELECT username, password FROM ACCOUNT WHERE username = ?`
)

// 登录
func Login(context *gin.Context) {
	var jsonData AuthForm
	if err := context.ShouldBind(&jsonData) ;err != nil {
		utils.CheckError(context, err)
	}else{

		username := strings.TrimSpace(jsonData.Username)
		password := strings.TrimSpace(jsonData.Password)
		// 查找用户
		account, err := models.GetAccount(QUERYACCOUNTBYUSERNAME, username)
		if err != nil {
			utils.CheckError(context, errors.New("wrong username"))
			return
		}
		if account.HasPassword(password) == false {
			utils.CheckError(context, errors.New("wrong password"))
			return
		}
		// TODO 颁发TOKEN
		session := sessions.Default(context)
		session.Set("isLogin", true)
		session.Set("username", username)
		session.Save()
		context.JSON(http.StatusOK, gin.H{
			"message": "login success",
		})
	}

}

// 注册
//type AuthForm struct {
//	Username string `form:"username" binding:"required"`
//	Password string `form:"password" binding:"required"`
//}

type AuthForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Nickname string `form:"nickname"`
	Sex      uint8  `form:"sex"`
	Email    string `form:"email"`
}

func Signup(context *gin.Context) {

	var jsonData AuthForm
	if err := context.ShouldBind(&jsonData); err != nil {
		utils.CheckError(context, err)

	} else {
		jsonData.Username = strings.TrimSpace(jsonData.Username)
		jsonData.Password = strings.TrimSpace(jsonData.Password)

		// 密码加密
		if pw, err := bcrypt.GenerateFromPassword([]byte(jsonData.Password), bcrypt.DefaultCost) ;err != nil {
			utils.CheckError(context, err)
		}else {
			jsonData.Password = string(pw[:])
		}
		if jsonData.Sex != 1 && jsonData.Sex != 2 {
			jsonData.Sex = 0
		}
		stmt, err := models.GetMySQLDB().Prepare(INSERTACCOUNT)
		if err != nil {
			utils.CheckError(context, err)
		}

		_, err = stmt.Exec(jsonData.Username, jsonData.Password, jsonData.Nickname, jsonData.Sex, jsonData.Email)
		if err != nil {
			utils.CheckError(context, err)
		}

		context.JSON(http.StatusOK, gin.H{
			"message": "sign up success",
		})
	}
}

func GetSignup(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signup.html", nil)
	return
}

func GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/login.html", nil)
	return
}

// 认证

// 密码哈希
//func HashPassword(plain string){
//	// TODO 对密码的各种校验
//	// 对密码哈希处理
//	if pw, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost) ;err != nil {
//		utils.CheckError(context, err)
//	}else {
//		jsonData.Password = string(pw[:])
//	}
//}

type Auth struct {
	//State `form:"state"`
	RedirectUri string `form:"redirect_uri" binding:"required"`
	//ResponseType `form:"response_type"`
}

func Authorize(context *gin.Context){

	session := sessions.Default(context)
	isLogin :=session.Get("isLogin")

	if isLogin == nil || isLogin == false {
		context.Redirect(http.StatusFound, "/login")
		return
	}

	var authPara Auth
	err := context.ShouldBind(&authPara)
	if err != nil {
		utils.CheckError(context, err)
	}

}