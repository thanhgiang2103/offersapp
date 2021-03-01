package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	route := gin.Default()

	//usersGroup := route.Group("user")

	route.Run(":3000")

}
