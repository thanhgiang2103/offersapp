package main

import (
	"context"
	"fmt"
	"offersapp/routes"

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
	usersGroup := route.Group("user")
	{
		usersGroup.POST("register", routes.UsersRegister)
		usersGroup.POST("login", routes.UserLogin)
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
