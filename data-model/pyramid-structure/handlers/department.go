package handlers

import (
	"net/http"

	"org-hierarchy/models"

	"github.com/labstack/echo/v4"
)

type DepartmentHandler struct {
	repo *models.DepartmentRepository
}

func NewDepartmentHandler(repo *models.DepartmentRepository) *DepartmentHandler {
	return &DepartmentHandler{repo: repo}
}

func (h *DepartmentHandler) GetAll(c echo.Context) error {
	departments, err := h.repo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, departments)
}

func (h *DepartmentHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	
	department, err := h.repo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	if department == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Department not found"})
	}
	
	return c.JSON(http.StatusOK, department)
}

func (h *DepartmentHandler) Create(c echo.Context) error {
	var department models.Department
	if err := c.Bind(&department); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	if department.DepartmentID == "" || department.DepartmentName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Department ID and name are required"})
	}
	
	if err := h.repo.Create(&department); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	return c.JSON(http.StatusCreated, department)
}

func (h *DepartmentHandler) Update(c echo.Context) error {
	id := c.Param("id")
	
	var department models.Department
	if err := c.Bind(&department); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	department.DepartmentID = id
	
	if department.DepartmentName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Department name is required"})
	}
	
	if err := h.repo.Update(&department); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	return c.JSON(http.StatusOK, department)
}

func (h *DepartmentHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	
	if err := h.repo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	return c.NoContent(http.StatusNoContent)
}