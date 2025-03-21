package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

// for user
func (h *ProjectHandler) PostProject(c *gin.Context) {
	var projectRequest dto.ProjectCreateRequest

	if err := c.BindJSON(&projectRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something is invalid"})
	}

	parsedTime, err := time.Parse("2006-01-02", projectRequest.Deadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date format is invalid"})
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

	err = h.projectRepo.Create(&newProject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
	}

	projectResponse := dto.ProjectResponse{
		ID:           newProject.ID,
		Title:        newProject.Title,
		Description:  newProject.Description,
		Platform:     newProject.Platform,
		Client:       newProject.Client,
		EstimatedFee: newProject.EstimatedFee,
		Status:       newProject.Status,
		Deadline:     newProject.Deadline.Format("2006-01-02"),
	}

	c.IndentedJSON(http.StatusCreated, projectResponse)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	// userName := c.Param("user")
	projectName := c.Param("project")
	// id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "ID is invalid"})
	// }

	project, err := h.projectRepo.Find(projectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Project not found"})
	}
	projectResponse := dto.ProjectResponse{
		ID:           project.ID,
		Title:        project.Title,
		Description:  project.Description,
		Platform:     project.Platform,
		Client:       project.Client,
		EstimatedFee: project.EstimatedFee,
		Status:       project.Status,
		Deadline:     project.Deadline.Format("2006-01-02"),
	}
	c.IndentedJSON(http.StatusOK, projectResponse)
}
