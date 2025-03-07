package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
