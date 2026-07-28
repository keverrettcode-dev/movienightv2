package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/keverrettcode-dev/movienightv2/backend/server/database"
	"github.com/keverrettcode-dev/movienightv2/backend/server/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection("users")

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

		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":"Validation Failed","details": err.Error()})
			return
		}

		HashPassword, err := HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Unable to hash password"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		count, err := userCollection.CountDocuments(ctx, bson.M{"email":user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to add existing user"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error":"User already exists"})
			return
		}

		user.UserId = bson.NewObjectID().Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.Password = HashPassword

		result, err := userCollection.InsertOne(ctx, user)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, result)


	}
}

func UserLogin() gin.HandlerFunc{
	return func(c *gin.Context) {
		var userLogin models.UserLogin
		if err := c.ShouldBindJSON(&userLogin); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid input data"})
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var foundUser models.User

		err := userCollection.FindOne(ctx, bson.M{"email": userLogin.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLogin.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
	}
}

type SignedDetails struct {
	Email string
	FirstName string
	LastName string
	Role string
	UserId string
	jwt.RegisteredClaims
}
var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email, firstName, lastName, role, userId string)(string, string, error) {
	claims := &SignedDetails{
		Email: email,
		FirstName: firstName,
		LastName: lastName,
		Role: role,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "MagicStream",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate((time.Now().Add(24 * time.Hour))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshedClaims := &SignedDetails{
		Email: email,
		FirstName: firstName,
		LastName: lastName,
		Role: role,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "MagicStream",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate((time.Now().Add(24 * time.Hour))),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshedClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}