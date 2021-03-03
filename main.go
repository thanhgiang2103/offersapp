package main

import (
	"context"
	"fmt"
	"net/http"
	"offersapp/models"
	"offersapp/routes"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v4"
)

func main() {
	conn, err := connectDB()
	if err != nil {
		return
	}

	route := gin.Default()
	route.Use(dbMiddleware(*conn))
	usersGroup := route.Group("users")
	{
		usersGroup.POST("register", routes.UsersRegister)
		usersGroup.POST("login", routes.UserLogin)
	}
	itemsGroup := route.Group("items")
	{
		itemsGroup.GET("index", routes.ItemIndex)
		itemsGroup.POST("create", authMiddleWare(), routes.ItemCreate)
		itemsGroup.GET("sold_by_user", authMiddleWare(), routes.ItemForSaleByCurrentUser)
		itemsGroup.PUT("update", authMiddleWare(), routes.ItemUpdate)
	}
	route.Run(":3000")

}

func connectDB() (conn *pgx.Conn, err error) {
	conn, err = pgx.Connect(context.Background(), "postgresql://postgres:15151515a@localhost:5432/offersapp")
	if err != nil {
		fmt.Println("Error connecting to DB")
		fmt.Println(err.Error())
	}
	_ = conn.Ping(context.Background())
	return
}

func dbMiddleware(conn pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	}
}

func authMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		split := strings.Split(bearer, "Bearer ")

		if len(split) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}
		token := split[1]

		isValid, userID := models.IsTokenValid(token)

		if isValid == false {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		} else {
			c.Set("user_id", userID)
			c.Next()
		}
	}
}
