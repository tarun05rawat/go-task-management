package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tarun05rawat/go-task-management/database"
	"github.com/tarun05rawat/go-task-management/model"
)

func CreateTask(c *gin.Context) {
	var task model.Task

	// ✅ Get `user_id` from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID"})
		return
	}

	// ✅ Bind request JSON to task struct
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Assign `UserID` to the task
	task.UserID = userIDUint

	// ✅ Fetch last task ID for this user
	var lastTask model.Task
	result := database.DB.Where("user_id = ?", task.UserID).Order("task_id DESC").First(&lastTask)

	// ✅ If no previous task exists, start from 1
	if result.Error != nil {
		task.TaskID = 1
	} else {
		task.TaskID = lastTask.TaskID + 1
	}

	// ✅ Save Task
	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// ✅ Get All Tasks (For Logged-In User)
func GetTasks(c *gin.Context) {
	var tasks []model.Task

	// ✅ Get `user_id`
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// ✅ Fetch only tasks that belong to the user
	database.DB.Where("user_id = ?", userID).Find(&tasks)

	c.JSON(http.StatusOK, tasks)
}

// ✅ Get Task by ID (Ensuring User Can Only Access Their Own Tasks)
func GetTaskByID(c *gin.Context) {
	var task model.Task

	// ✅ Get `user_id`
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// ✅ Get `task_id` from URL param
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// ✅ Fetch task (Ensuring it's owned by the user)
	result := database.DB.Where("user_id = ? AND task_id = ?", userID, taskID).First(&task)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or does not belong to you"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// ✅ Update Task (Ensuring User Can Only Update Their Own Tasks)
func UpdateTask(c *gin.Context) {
	var task model.Task

	// ✅ Get `user_id`
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// ✅ Get `task_id` from request parameters
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// ✅ Check if the task exists and belongs to the logged-in user
	result := database.DB.Where("user_id = ? AND task_id = ?", userID, taskID).First(&task)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or does not belong to you"})
		return
	}

	// ✅ Bind new task data from request body
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Save the updated task
	database.DB.Save(&task)

	// ✅ Return success response
	c.JSON(http.StatusOK, task)
}

// ✅ Delete Task (Ensuring User Can Only Delete Their Own Tasks)
func DeleteTask(c *gin.Context) {
	var task model.Task

	// ✅ Get `user_id`
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// ✅ Get `task_id` from request parameters
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// ✅ Check if the task exists and belongs to the logged-in user
	result := database.DB.Where("user_id = ? AND task_id = ?", userID, taskID).First(&task)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or does not belong to you"})
		return
	}

	// ✅ Delete Task
	database.DB.Delete(&task)

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// ✅ Get All Users (Admin Only)
func GetAllUsers(c *gin.Context) {
	// ✅ Get `user_id` from context
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
