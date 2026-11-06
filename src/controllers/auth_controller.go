package controllers

import (
	"net/http"

	"flexoo-service-order/src/middlewares"
	"flexoo-service-order/src/models"
	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (ac *AuthController) ShowLogin(c *gin.Context) {
	if cookie, err := c.Cookie("flexoo_session"); err == nil && cookie != "" {
		if _, ok := middlewares.ParseSessionToken(cookie); ok {
			c.Redirect(http.StatusFound, "/")
			return
		}
	}
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (ac *AuthController) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	client, err := models.AuthenticateClient(username, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"Error":    "Invalid username or password. Please try again.",
			"Username": username,
		})
		return
	}

	middlewares.SetAuthCookie(c, client)
	c.Redirect(http.StatusSeeOther, "/")
}

func (ac *AuthController) Logout(c *gin.Context) {
	middlewares.ClearAuthCookie(c)
	c.Redirect(http.StatusSeeOther, "/login")
}

