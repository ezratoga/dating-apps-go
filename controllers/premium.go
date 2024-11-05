package controllers

import (
    "dating-apps-go/database"
    "dating-apps-go/models"
    "github.com/labstack/echo/v4"
    "net/http"
	"github.com/golang-jwt/jwt"
)

// PurchasePremium handles premium upgrades
func PurchasePremium(header echo.Context) error {
    // var profile models.User
    authenticationHeader := header.Get("user").(*jwt.Token)
	userData := authenticationHeader.Claims.(jwt.MapClaims)
	// Retrieve userID from userData
    userID := uint(userData["userID"].(float64)) // Cast to uint

    var user models.User
	var invoice models.Invoice
	var transaction models.InvoicePayload

	if err := header.Bind(&transaction); err != nil {
        return header.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
    }

	if err := database.DB.Model(&user).Where("id = ?", userID).First(&user).Error; err != nil {
        return header.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
    }

    if err := database.DB.Model(&invoice).Where("invoice_id = ? and user_id = ?", transaction.InvoiceID, userID, "complete").Error; err != nil {
        return header.JSON(http.StatusNotFound, map[string]string{"error": "Invoice not found"})
    }

	if invoice.Status == "complete" {
		return header.JSON(http.StatusBadRequest, map[string]string{"error": "your account has premium feature"})
	}

    user.IsPremium = true
	user.Verified = true
    database.DB.Save(&user)

    return header.JSON(http.StatusOK, "User upgraded to premium")
}
