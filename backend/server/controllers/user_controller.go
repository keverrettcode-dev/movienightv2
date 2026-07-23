package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	models "github.com/keverrettcode-dev/movienightv2/backend/server/models"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string)(string, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(HashPassword), nil
}

func RegisterUser() gin.HandlerFunc{
	return func(c *gin.Context) {
		var user models.User

		if err:=c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input."})
			return
		}
		validator := validator.New()

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":"Validation Failed","details": err.Error()})
			return
		}
	}
}