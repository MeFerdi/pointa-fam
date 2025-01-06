package controllers

import (
    "net/http"
    "pointafam/backend/models"
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func AddToCart(c *gin.Context) {
    var cartItem models.CartItem
    if err := c.ShouldBindJSON(&cartItem); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    db := c.MustGet("db").(*gorm.DB)
    if err := cartItem.AddCartItem(db); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add to cart"})
        return
    }

    c.JSON(http.StatusCreated, cartItem)
}

func RemoveFromCart(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    db := c.MustGet("db").(*gorm.DB)
    if err := models.DeleteCartItem(db, uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove from cart"})
        return
    }

    c.JSON(http.StatusNoContent, nil)
}

func ViewCart(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    retailerIDStr := c.Param("retailer_id")
    retailerID, err := strconv.ParseUint(retailerIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid retailer ID"})
        return
    }

    cart, err := models.GetCart(db, uint(retailerID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
        return
    }

    c.JSON(http.StatusOK, cart)
}

func DeleteFromCart(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    itemIDStr := c.Param("item_id")
    itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
        return
    }

    if err := models.DeleteCartItem(db, uint(itemID)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item from cart"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
}
func GetCartItems(c *gin.Context) {
    retailerIDStr := c.Param("retailer_id")
    retailerID, err := strconv.ParseUint(retailerIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid retailer ID"})
        return
    }

    db := c.MustGet("db").(*gorm.DB)
    cart, err := models.GetCart(db, uint(retailerID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch cart items"})
        return
    }

    c.JSON(http.StatusOK, cart.Items)
}