package routes

import (
	"fmt"
	"net/http"
	"offersapp/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func ItemIndex(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	items, err := models.GetAllItems(&conn)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func ItemCreate(c *gin.Context) {

	item := models.Item{}
	err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	userID := c.GetString("user_id")

	err = item.Create(&conn, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func ItemForSaleByCurrentUser(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	userID := c.GetString("user_id")

	items, err := models.GetItemsBeingSoldByUser(userID, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func ItemUpdate(c *gin.Context) {
	itemSent := models.Item{}

	err := c.ShouldBindJSON(&itemSent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	userID := c.GetString("user_id")

	itemBeingUpdated, err := models.FindItemById(itemSent.ID, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if itemBeingUpdated.SellerID.String() != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You aren't authorized to update this item"})
		return
	}
	itemSent.SellerID = itemBeingUpdated.SellerID

	err = itemSent.Update(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": itemSent})
}
