package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"flexoo-service-order/src/models"
	"github.com/gin-gonic/gin"
)

const sessionCookieName = "flexoo_session"

func secretKey() []byte {
	s := os.Getenv("SESSION_SECRET")
	if s == "" {
		s = "flexoo-service-order-default-secret-key-2026"
	}
	return []byte(s)
}

func CreateSessionToken(client *models.Client) string {
	payload := fmt.Sprintf("%d|%s|%s|%s", client.ClientID, client.Username, client.FullName, client.Email)
	encoded := base64.RawURLEncoding.EncodeToString([]byte(payload))
	mac := hmac.New(sha256.New, secretKey())
	mac.Write([]byte(encoded))
	return encoded + "." + hex.EncodeToString(mac.Sum(nil))
}

func ParseSessionToken(token string) (*models.Client, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, false
	}
	mac := hmac.New(sha256.New, secretKey())
	mac.Write([]byte(parts[0]))
	if !hmac.Equal([]byte(parts[1]), []byte(hex.EncodeToString(mac.Sum(nil)))) {
		return nil, false
	}

	raw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, false
	}
	fields := strings.SplitN(string(raw), "|", 4)
	if len(fields) != 4 {
		return nil, false
	}
	id, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return nil, false
	}
	return &models.Client{
		ClientID: id,
		Username: fields[1],
		FullName: fields[2],
		Email:    fields[3],
	}, true
}

func SetAuthCookie(c *gin.Context, client *models.Client) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(sessionCookieName, CreateSessionToken(client), 86400*7, "/", "", false, true)
}

func ClearAuthCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(sessionCookieName, "", -1, "/", "", false, true)
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(sessionCookieName)
		if err != nil || cookie == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		client, ok := ParseSessionToken(cookie)
		if !ok {
			ClearAuthCookie(c)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Set("client", client)
		c.Next()
	}
}

func CurrentClient(c *gin.Context) *models.Client {
	val, exists := c.Get("client")
	if !exists {
		return nil
	}
	client, _ := val.(*models.Client)
	return client
}

