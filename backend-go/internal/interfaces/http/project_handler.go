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

type ProjectHandler struct {
	projectRepo repositories.ProjectRepository
}

func NewProjectHandler(projectRepo repositories.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{projectRepo: projectRepo}
}

func projectBindAndValidate(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errors := make([]map[string]string, len(ve))
			for i, fe := range ve {
				errors[i] = map[string]string{
					"field":   fe.Field(),
					"message": projectValidationErrorMessage(fe),
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

func projectValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "gte":
		return "must be greater than or equal to " + fe.Param()
	default:
		return "is invalid"
	}
}

// ==================================================================================================================
// for user
func (h *ProjectHandler) PostProject(c *gin.Context) {
	userName := c.Param("user")
	var projectRequest dto.ProjectCreateRequest

	if !projectBindAndValidate(c, &projectRequest) {
		return
	}

	parsedTime, err := time.Parse("2006-01-02", projectRequest.Deadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date format is invalid (yyyy-MM-dd)"})
		return
	}

	newProject := models.Project{
		Title:        projectRequest.Title,
		Description:  projectRequest.Description,
		Platform:     projectRequest.Platform,
		Client:       projectRequest.Client,
		EstimatedFee: projectRequest.EstimatedFee,
		Status:       projectRequest.Status,
		Deadline:     parsedTime,
	}

	if err := h.projectRepo.Create(userName, &newProject); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	projectResponse := dto.ProjectResponse{
		Number:       newProject.Number,
		Title:        newProject.Title,
		Description:  newProject.Description,
		Platform:     newProject.Platform,
		Client:       newProject.Client,
		EstimatedFee: newProject.EstimatedFee,
		Status:       newProject.Status,
		Deadline:     newProject.Deadline.Format("2006-01-02"),
		CreatedAt:    newProject.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    newProject.UpdatedAt.Format(time.RFC3339),
	}

	c.IndentedJSON(http.StatusCreated, projectResponse)
}

func (h *ProjectHandler) GetAllProjectsByUser(c *gin.Context) {
	userName := c.Param("user")
	projects, err := h.projectRepo.FindAll(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
		return
	}

	var projectResponse []dto.ProjectResponse
	for _, project := range projects {
		projectResponse = append(projectResponse, dto.ProjectResponse{
			Number:       project.Number,
			Title:        project.Title,
			Description:  project.Description,
			Platform:     project.Platform,
			Client:       project.Client,
			EstimatedFee: project.EstimatedFee,
			Status:       project.Status,
			Deadline:     project.Deadline.Format("2006-01-02"),
			CreatedAt:    project.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    project.UpdatedAt.Format(time.RFC3339),
		})
	}
	c.IndentedJSON(http.StatusOK, projectResponse)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	userName := c.Param("user")
	id, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is invalid"})
		return
	}

	project, err := h.projectRepo.Find(userName, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	projectResponse := dto.ProjectResponse{
		Number:       project.Number,
		Title:        project.Title,
		Description:  project.Description,
		Platform:     project.Platform,
		Client:       project.Client,
		EstimatedFee: project.EstimatedFee,
		Status:       project.Status,
		Deadline:     project.Deadline.Format("2006-01-02"),
		CreatedAt:    project.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    project.UpdatedAt.Format(time.RFC3339),
	}
	c.IndentedJSON(http.StatusOK, projectResponse)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	userName := c.Param("user")
	id, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project ID is invalid"})
		return
	}

	var projectRequest dto.ProjectUpdateRequest
	if !projectBindAndValidate(c, &projectRequest) {
		return
	}

	parsedTime, err := time.Parse("2006-01-02", projectRequest.Deadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date format is invalid (yyyy-MM-dd)"})
		return
	}

	targetProject := models.Project{
		Title:        projectRequest.Title,
		Description:  projectRequest.Description,
		Platform:     projectRequest.Platform,
		Client:       projectRequest.Client,
		EstimatedFee: projectRequest.EstimatedFee,
		Status:       projectRequest.Status,
		Deadline:     parsedTime,
	}

	updatedProject, err := h.projectRepo.Update(userName, uint(id), &targetProject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	projectResponse := dto.ProjectResponse{
		Number:       updatedProject.Number,
		Title:        updatedProject.Title,
		Description:  updatedProject.Description,
		Platform:     updatedProject.Platform,
		Client:       updatedProject.Client,
		EstimatedFee: updatedProject.EstimatedFee,
		Status:       updatedProject.Status,
		Deadline:     updatedProject.Deadline.Format("2006-01-02"),
	}

	c.IndentedJSON(http.StatusOK, projectResponse)
}

func (h *ProjectHandler) SoftDeleteProject(c *gin.Context) {
	userName := c.Param("user")
	id, err := strconv.ParseUint(c.Param("pid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is invalid"})
		return
	}

	if err := h.projectRepo.SoftDelete(userName, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.Status(http.StatusNoContent)

}

// ==================================================================================================================

// ==================================================================================================================
// for owner
func (h *ProjectHandler) HardDeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is invalid"})
		return
	}

	if err := h.projectRepo.HardDelete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProjectHandler) GetAllProjectsForOwner(c *gin.Context) {
	projects, err := h.projectRepo.FindAllForOwner()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all users"})
		return
	}

	var projectResponse []dto.ProjectResponseForOwner
	for _, project := range projects {
		projectResponse = append(projectResponse, dto.ProjectResponseForOwner{
			ID:           project.ID,
			UserID:       project.UserID,
			Number:       project.Number,
			Title:        project.Title,
			Description:  project.Description,
			Platform:     project.Platform,
			Client:       project.Client,
			EstimatedFee: project.EstimatedFee,
			Status:       project.Status,
			Deadline:     project.Deadline.Format("2006-01-02"),
			CreatedAt:    project.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    project.UpdatedAt.Format(time.RFC3339),
			DeletedAt:    project.DeletedAt.Time.Format(time.RFC3339),
		})
	}
	c.IndentedJSON(http.StatusOK, projectResponse)
}

// ==================================================================================================================
