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
		Role     string `json:"role"` // ✅ Allow role input (optional, defaults to "user")
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Default role to "user" if not provided
	if body.Role == "" {
		body.Role = "user"
	}

	// Create user
	user := model.User{
		Username:     body.Username,
		Email:        body.Email,
		PasswordHash: string(hash),
		Role:         body.Role, // ✅ Assign role
	}
	result := database.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Look up user by email
	var user model.User
	if err := database.DB.First(&user, "email = ?", body.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email/password"})
		return
	}

	// Compare password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email/password"})
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
		"role":    user.Role,                                  // ✅ Include role in token
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(), // Expires in 30 days
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	fmt.Println("Generated Token:", tokenString)

	// Set the JWT token in an HTTP-only cookie
	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", false, true)

	// ✅ Return token in response body (for easy access in Postman)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
		"token":   tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// Validate function (checks if user is authenticated)
func Validate(c *gin.Context) {
	user, exists := c.Get("user") // Retrieve user from context (set in RequireAuth)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// ✅ Ensure user is returned properly
	userData, ok := user.(model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User data error"})
		return
	}

	// Return the authenticated user's details
	c.JSON(http.StatusOK, gin.H{
		"message": "User is authenticated",
		"user": gin.H{
			"id":       userData.ID,
			"username": userData.Username,
			"email":    userData.Email,
			"role":     userData.Role,
		},
	})
}

// Logout function (clears the authentication cookie)
func Logout(c *gin.Context) {
	// Set the cookie to expire immediately
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GetAllUsers - Admin-only endpoint to fetch all users
func GetAllUsers(c *gin.Context) {
	// ✅ Get `user_id` from context (to verify admin access)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// ✅ Fetch the requesting user's details
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// ✅ Check if the user is an admin
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// ✅ Fetch all users
	var users []model.User
	database.DB.Find(&users)

	// ✅ Return list of users
	c.JSON(http.StatusOK, users)
}
