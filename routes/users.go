package routes

import (
	"net/http"
	"offersapp/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func UserLogin(c *gin.Context) {
	user := models.User{}
	err := c.ShouldBindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	err = user.IsAuthenticated(&conn)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := user.GetAuthToken()
	if err == nil {

		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "There is an error authenticating "})
}

func UsersRegister(c *gin.Context) {

	user := models.User{}
	err := c.ShouldBindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	err = user.Register(&conn)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := user.GetAuthToken()
	if err == nil {

		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": user.ID})
}
