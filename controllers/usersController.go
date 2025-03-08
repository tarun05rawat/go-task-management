package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tarun05rawat/go-task-management/database"
	"github.com/tarun05rawat/go-task-management/model"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	//get the email/pass of req body
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	//Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil { //Bind return nothing. if it does, it has to be an error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Create user
	user := model.User{Username: body.Username, Email: body.Email, PasswordHash: string(hash)}
	result := database.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or email already exists"})
		return
	}

	//Response
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	}})
}

func Login(c *gin.Context) {
	// get the email and password off req body
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	// Look up requested user
	var user model.User
	database.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email/password"})
		return
	}

	// Compare sent in password to user's actual password hash
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email/password"})
		return
	}

	//generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(os.Getenv("SECRET"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create token"})
		return
	}

	fmt.Println(tokenString, err)

	//send it back
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
