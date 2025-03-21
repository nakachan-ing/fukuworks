package http

import (
	"net/http"
	"strconv"
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
		return
	}

	parsedTime, err := time.Parse("2006-01-02", projectRequest.Deadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date format is invalid"})
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

	err = h.projectRepo.Create(&newProject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
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
	}

	c.IndentedJSON(http.StatusCreated, projectResponse)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	userName := c.Param("user")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is invalid"})
	}

	project, err := h.projectRepo.Find(userName, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Project not found"})
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
	}
	c.IndentedJSON(http.StatusOK, projectResponse)
}

// for user
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	userName := c.Param("user")
	projectName := c.Param("projects")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is invalid"})
	}

	var projectRequest dto.ProjectUpdateRequest
	if err := c.BindJSON(&projectRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something is invalid"})
		return
	}

	parsedTime, err := time.Parse("2006-01-02", projectRequest.Deadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date format is invalid"})
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

	updatedProject, err := h.projectRepo.Update(userName, projectName, uint(id), &targetProject)
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

	c.IndentedJSON(http.StatusCreated, projectResponse)

}

func (h *ProjectHandler) SoftDeleteProject(c *gin.Context) {
	userName := c.Param("user")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is invalid"})
	}

	err = h.projectRepo.SoftDelete(userName, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.Status(http.StatusNoContent)

}

func (h *ProjectHandler) GetAllProjectsByUser(c *gin.Context) {
	userName := c.Param("user")
	projects, err := h.projectRepo.FindAll(userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
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
		})
	}
	c.IndentedJSON(http.StatusOK, projectResponse)
}

func (h *ProjectHandler) HardDeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is invalid"})
		return
	}

	err = h.projectRepo.HardDelete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.Status(http.StatusNoContent)
}
