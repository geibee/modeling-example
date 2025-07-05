package handlers

import (
	"net/http"
	"time"

	"org-hierarchy/models"

	"github.com/labstack/echo/v4"
)

type OrganizationAttributeHandler struct {
	repo *models.OrganizationAttributeRepository
}

func NewOrganizationAttributeHandler(repo *models.OrganizationAttributeRepository) *OrganizationAttributeHandler {
	return &OrganizationAttributeHandler{repo: repo}
}

func (h *OrganizationAttributeHandler) GetAll(c echo.Context) error {
	attrs, err := h.repo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, attrs)
}

func (h *OrganizationAttributeHandler) GetByID(c echo.Context) error {
	departmentID := c.Param("department_id")
	effectiveDateStr := c.Param("effective_date")
	
	effectiveDate, err := time.Parse("2006-01-02", effectiveDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format. Use YYYY-MM-DD"})
	}
	
	attr, err := h.repo.GetByID(departmentID, effectiveDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	if attr == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Organization attribute not found"})
	}
	
	return c.JSON(http.StatusOK, attr)
}

func (h *OrganizationAttributeHandler) GetDepartmentHistory(c echo.Context) error {
	departmentID := c.Param("id")
	
	attrs, err := h.repo.GetDepartmentHistory(departmentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	return c.JSON(http.StatusOK, attrs)
}

func (h *OrganizationAttributeHandler) Create(c echo.Context) error {
	var input struct {
		DepartmentID       string  `json:"department_id"`
		EffectiveDate      string  `json:"effective_date"`
		ParentDepartmentID *string `json:"parent_department_id"`
	}
	
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	if input.DepartmentID == "" || input.EffectiveDate == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Department ID and effective date are required"})
	}
	
	effectiveDate, err := time.Parse("2006-01-02", input.EffectiveDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format. Use YYYY-MM-DD"})
	}
	
	attr := &models.OrganizationAttribute{
		DepartmentID:       input.DepartmentID,
		EffectiveDate:      effectiveDate,
		ParentDepartmentID: input.ParentDepartmentID,
	}
	
	if err := h.repo.Create(attr); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	// 作成後、失効年月を含む完全なデータを取得
	createdAttr, err := h.repo.GetByID(attr.DepartmentID, attr.EffectiveDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	return c.JSON(http.StatusCreated, createdAttr)
}

func (h *OrganizationAttributeHandler) Update(c echo.Context) error {
	departmentID := c.Param("department_id")
	effectiveDateStr := c.Param("effective_date")
	
	effectiveDate, err := time.Parse("2006-01-02", effectiveDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format. Use YYYY-MM-DD"})
	}
	
	var input struct {
		ParentDepartmentID *string `json:"parent_department_id"`
	}
	
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	attr := &models.OrganizationAttribute{
		DepartmentID:       departmentID,
		EffectiveDate:      effectiveDate,
		ParentDepartmentID: input.ParentDepartmentID,
	}
	
	if err := h.repo.Update(attr); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	// 更新後、失効年月を含む完全なデータを取得
	updatedAttr, err := h.repo.GetByID(attr.DepartmentID, attr.EffectiveDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	return c.JSON(http.StatusOK, updatedAttr)
}

func (h *OrganizationAttributeHandler) Delete(c echo.Context) error {
	departmentID := c.Param("department_id")
	effectiveDateStr := c.Param("effective_date")
	
	effectiveDate, err := time.Parse("2006-01-02", effectiveDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format. Use YYYY-MM-DD"})
	}
	
	if err := h.repo.Delete(departmentID, effectiveDate); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	return c.NoContent(http.StatusNoContent)
}