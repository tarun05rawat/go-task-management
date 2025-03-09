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

// Signup function for registering new users
func Signup(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
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

	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// Login function (sets JWT token in an HTTP-only cookie)
func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	// Look up user by email
	var user model.User
	if err := database.DB.First(&user, "email = ?", body.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email/password"})
		return
	}

	// Compare password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email/password"})
		return
	}

	// Get JWT secret key (fallback if not set)
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "your-fallback-secret-key"
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(), // Expires in 30 days
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	fmt.Println("Generated Token:", tokenString)

	// Set the JWT token in an HTTP-only cookie
	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", false, true)

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

func Validate(c *gin.Context) {
	user, exists := c.Get("user") // Retrieve user from context (set in RequireAuth)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Return the authenticated user's details
	c.JSON(http.StatusOK, gin.H{
		"message": "User is authenticated",
		"user":    user,
	})
}

// Logout function (clears the authentication cookie)
func Logout(c *gin.Context) {
	// Set the cookie to expire immediately
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
