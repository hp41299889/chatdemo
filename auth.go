package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Username string
	Password string
	Nickname string
}

func LoginAuth(c *gin.Context) {
	var (
		username string
		password string
	)
	if in, isExist := c.GetPostForm("username"); isExist && in != "" {
		username = in
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": errors.New("username can't be empty."),
		})
		return
	}
	if in, isExist := c.GetPostForm("password"); isExist && in != "" {
		password = in
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": errors.New("password can't be empty"),
		})
		return
	}
	if nickname, err := Auth(username, password); err == nil {
		c.Redirect(http.StatusMovedPermanently, "/chat?nickname="+nickname)

	} else {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": err,
		})
		return
	}
}

func CheckPassword(db *gorm.DB, user *User, passwordIn string) (string, error) {
	if passwordIn != user.Password {
		return "", errors.New("password is not correct")
	} else {
		return user.Nickname, nil
	}
}

func FindUser(db *gorm.DB, usernameIn string) (*User, error) {
	user := new(User)
	user.Username = usernameIn
	err := db.Where("username = ?", usernameIn).First(&user).Error
	return user, err
}

func Auth(usernameIn, passwordIn string) (string, error) {
	dsn := "query:query123@tcp(34.81.16.122:3306)/demo"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("connect to cloudsql error: ", err)
	}
	if user, err := FindUser(db, usernameIn); err != nil {
		return "", errors.New("user not found")
	} else {
		return CheckPassword(db, user, passwordIn)
	}
}
