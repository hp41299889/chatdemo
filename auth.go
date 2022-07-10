package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Username string
	Password string
	Nickname string
}

func connetDB() *gorm.DB {
	viper.SetConfigName("mysql")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("viper error: ", err)
	}
	USERNAME := viper.GetString("db.username")
	PASSWORD := viper.GetString("db.password")
	NETWORK := viper.GetString("db.network")
	SERVER := viper.GetString("db.server")
	PORT := viper.GetInt32("db.port")
	DATABASE := viper.GetString("db.database")
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("gorm error: ", err)
	}
	return db
}

func LoginAuth(c *gin.Context) {
	var (
		username string
		password string
	)
	if in, isExist := c.GetPostForm("username"); isExist && in != "" {
		username = in
	} else {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error": errors.New("username can't be empty."),
		})
		return
	}
	if in, isExist := c.GetPostForm("password"); isExist && in != "" {
		password = in
	} else {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error": errors.New("password can't be empty"),
		})
		return
	}
	if nickname, err := Auth(username, password); err != nil {
		c.HTML(http.StatusUnauthorized, "index.html", gin.H{
			"error": err,
		})
		return
	} else {
		c.Redirect(http.StatusMovedPermanently, "/chat?nickname="+nickname)
	}
}

func CheckPassword(user *User, passwordIn string) error {
	if passwordIn != user.Password {
		return errors.New("password is not correct")
	}
	return nil
}

func FindUser(db *gorm.DB, username string) (*User, error) {
	user := new(User)
	user.Username = username
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, err
}

func Auth(username, password string) (string, error) {
	db := connetDB()
	if user, err := FindUser(db, username); err != nil {
		return "", err
	} else {
		if err := CheckPassword(user, password); err != nil {
			return "", err
		}
		return user.Username, nil
	}
}
