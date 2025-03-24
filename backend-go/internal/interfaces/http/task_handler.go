package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/dto"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/repositories"
)

type TaskHandler struct {
	taskRepo repositories.TaskRepository
}

func NewTaskHandler(taskRepo repositories.TaskRepository) *TaskHandler {
	return &TaskHandler{taskRepo: taskRepo}
}

func taskBindAndValidate(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errors := make([]map[string]string, len(ve))
			for i, fe := range ve {
				errors[i] = map[string]string{
					"field":   fe.Field(),
					"message": taskValidationErrorMessage(fe),
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return false
	}
	return true
}

func taskValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "gte":
		return "must be greater than or equal to " + fe.Param()
	default:
		return "is invalid"
	}
}

// for user
func (h *TaskHandler) PostTask(c *gin.Context) {
	userName := c.Param("user")
	pid, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is invalid"})
		return
	}
	var taskRequest dto.TaskCreateRequest

	if !taskBindAndValidate(c, &taskRequest) {
		return
	}

	parsedTime, err := time.Parse("2006-01-02", taskRequest.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date format is invalid (yyyy-MM-dd)"})
		return
	}

	newTask := models.Task{
		Title:       taskRequest.Title,
		Description: taskRequest.Description,
		Status:      taskRequest.Status,
		Priority:    taskRequest.Priority,
		DueDate:     parsedTime,
	}

	if err := h.taskRepo.Create(userName, uint(pid), &newTask); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	taskResponse := dto.TaskResponse{
		Number:      newTask.Number,
		Title:       newTask.Title,
		Description: newTask.Description,
		Status:      newTask.Status,
		Priority:    newTask.Priority,
		DueDate:     newTask.DueDate.Format("2006-01-02"),
	}

	c.IndentedJSON(http.StatusCreated, taskResponse)
}

func (h *TaskHandler) GetAllTasksByProject(c *gin.Context) {
	userName := c.Param("user")
	pid, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is invalid"})
		return
	}

	tasks, err := h.taskRepo.FindAll(userName, uint(pid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
		return
	}

	var taskResponse []dto.TaskResponse
	for _, task := range tasks {
		taskResponse = append(taskResponse, dto.TaskResponse{
			Number:      task.Number,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
			DueDate:     task.DueDate.Format("2006-01-02"),
			CreatedAt:   task.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
		})
	}
	c.IndentedJSON(http.StatusOK, taskResponse)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	userName := c.Param("user")
	pid, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is invalid"})
		return
	}
	tid, err := strconv.ParseUint(c.Param("tid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is invalid"})
		return
	}

	task, err := h.taskRepo.Find(userName, uint(pid), uint(tid))

	taskResponse := dto.TaskResponse{
		Number:      task.Number,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     task.DueDate.Format("2006-01-02"),
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}
	c.IndentedJSON(http.StatusOK, taskResponse)

}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userName := c.Param("user")
	pid, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is invalid"})
		return
	}
	tid, err := strconv.ParseUint(c.Param("tid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is invalid"})
		return
	}

	var taskRequest dto.TaskUpdateRequest
	if !taskBindAndValidate(c, &taskRequest) {
		return
	}

	parsedTime, err := time.Parse("2006-01-02", taskRequest.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date format is invalid (yyyy-MM-dd)"})
		return
	}

	targetTask := models.Task{
		Title:       taskRequest.Title,
		Description: taskRequest.Description,
		Status:      taskRequest.Status,
		Priority:    taskRequest.Priority,
		DueDate:     parsedTime,
	}

	updatedTask, err := h.taskRepo.Update(userName, uint(pid), uint(tid), &targetTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	taskResponse := dto.TaskResponse{
		Number:      updatedTask.Number,
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Status:      updatedTask.Status,
		Priority:    updatedTask.Priority,
		DueDate:     updatedTask.DueDate.Format("2006-01-02"),
		CreatedAt:   updatedTask.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   updatedTask.UpdatedAt.Format(time.RFC3339),
	}

	c.IndentedJSON(http.StatusOK, taskResponse)
}

func (h *TaskHandler) SoftDeleteTask(c *gin.Context) {
	userName := c.Param("user")
	pid, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is invalid"})
		return
	}
	tid, err := strconv.ParseUint(c.Param("tid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is invalid"})
		return
	}

	if err := h.taskRepo.SoftDelete(userName, uint(pid), uint(tid)); err != nil {
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) GetAllTasksForOwner(c *gin.Context) {
	tasks, err := h.taskRepo.FindAllTasksForOwner()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all tasks"})
		return
	}

	var taskResponse []dto.TaskResponseForOwner
	for _, task := range tasks {
		taskResponse = append(taskResponse, dto.TaskResponseForOwner{
			ID:          task.ID,
			ProjectID:   task.ProjectID,
			Number:      task.Number,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
			DueDate:     task.DueDate.Format("2006-01-02"),
			CreatedAt:   task.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
			DeletedAt:   task.DeletedAt.Time.Format(time.RFC3339),
		})
	}
	c.IndentedJSON(http.StatusOK, taskResponse)
}

func (h *TaskHandler) HardDeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is invalid"})
		return
	}

	if err := h.taskRepo.HardDelete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.Status(http.StatusNoContent)
}
