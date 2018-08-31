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
	"time"
	"strconv"
	"crypto/md5"
	"io"
	"fmt"
	"encoding/hex"
)



const (
	// 插入账户
	INSERTACCOUNT = `INSERT INTO ACCOUNT (username, password, nickname, sex, email, created_at, updated_at) VALUES (?,?,?,?,?,?,?)`

	// 根据用户名查找账户
	QUERYACCOUNTBYUSERNAME = `SELECT id, username, password FROM ACCOUNT WHERE username = ?`

	// 生成Ticket入库
	INSERTTICKET = `INSERT INTO Ticket (ticket, ttl, user_id, created_at, updated_at) VALUES (?,?,?,?,?)`

	// 查询Ticket //TODO 得好好优化
	QUERYTICKETBYTICKET = `SELECT ticket, ttl, created_at, updated_at, is_verify FROM Ticket WHERE is_verify = 0 AND ticket=?  order by created_at desc`

	UPDATETICKETINVAILD = `update Ticket set is_verify = 1`
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

		// 1、 设置session
		session := sessions.Default(context)
		session.Set("isLogin", true)
		session.Set("username", account.Username)
		session.Set("id", account.Id)
		session.Save()

		// 2、ticket处理
		ticketHandle(context, account)
	}

}

// 存入ticket
func ticketHandle(context *gin.Context, account *models.Account){

	// 2、生成ticket，并存入数据库
	ticket := GeneraterTicket(account.Username)
	var ticketObj = models.Ticket{Ticket:ticket,TTL: 60*5, User_id: account.Id}
	if err := ticketObj.CreateTicket(INSERTTICKET); err != nil {
		utils.CheckError(context, err)
		return
	}
	// 3、重定向到指定Callback
	var redirect = context.Query("redirect_uri")
	redirect = strings.Join([]string{redirect, fmt.Sprintf("ticket=%s", ticket)}, "?")
	context.Redirect(http.StatusFound, redirect)

	return
}


// ticket 有效期5分钟， 校验后马上失效
func GeneraterTicket(username string) string{

	rd := strings.Join([]string{strconv.FormatInt(time.Now().Unix(), 10), username},"")
	hash := md5.New()
	io.WriteString(hash, rd)

	ticket := hex.EncodeToString(hash.Sum(nil))
	fmt.Println(rd, ticket,"rd, ticket............")
	return  ticket

}



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
		timeStamp := time.Now().Unix()
		_, err = stmt.Exec(jsonData.Username, jsonData.Password, jsonData.Nickname, jsonData.Sex, jsonData.Email, timeStamp, timeStamp)
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


type Auth struct {
	//State `form:"state"`
	RedirectUri string `form:"redirect_uri" binding:"required"`
	//ResponseType `form:"response_type"`
}

func Authorize(context *gin.Context){
	session := sessions.Default(context)
	isLogin :=session.Get("isLogin")

	if isLogin == nil || isLogin == false {

		var redirect  = "/login"
		redirectUri := context.Query("redirect_uri")
		fmt.Println("redirect_uriredirect_uriredirect_uri   ", redirectUri)
		if len(redirectUri) == 0 || strings.HasPrefix(redirectUri, "http") == false {

			utils.CheckError(context, errors.New("redirect_uri parameter invalid ..."))
			return
		}else {
			redirect = strings.Join([]string{redirect, fmt.Sprintf("redirect_uri=%s", redirectUri)},"?")
		}

		context.HTML(http.StatusFound, "auth/login.html", gin.H{
			"redirect_uri": redirectUri,
		})
		return
	}

	var authPara Auth
	err := context.ShouldBind(&authPara)
	if err != nil {
		utils.CheckError(context, err)
	}

	username := session.Get("username")
	id := session.Get("id")
	ticketHandle(context, &models.Account{Username: username.(string), Id: id.(int)})
	return
}

// 校验ticket
func VerifyTicket(context *gin.Context){

	var ticket = context.PostForm("ticket")
	if len(ticket) == 0 {
		utils.CheckError(context, errors.New("ticket is empty"))
		return
	}
	ticketObj := models.Ticket{Ticket: ticket}
	if ticketObj, err := ticketObj.FindOneByTicket(QUERYTICKETBYTICKET);err != nil{

		utils.CheckError(context, err)
		return
	}else {
		curTimeStamp := time.Now().Unix()
		createAtTimeStamp := ticketObj.CreatedAt

		if int(curTimeStamp - createAtTimeStamp) > ticketObj.TTL {
			utils.CheckError(context, errors.New("ticket  expire"))
		}else{
			// 失效ticket
			if err := ticketObj.InvalidTicket(UPDATETICKETINVAILD); err != nil {
				utils.CheckError(context, err)
			}else{
				context.JSON(http.StatusOK, gin.H{
					"verirystatus": true,
				})
				return
			}
		}
	}
}

