package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"code-system/models"

	"github.com/labstack/echo/v4"
)

var DB *sql.DB

func SetDB(database *sql.DB) {
	DB = database
}

func GetItems(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	
	categoryType := c.QueryParam("category_type")
	
	offset := (page - 1) * pageSize
	
	query := `
		SELECT i.品目ID, i.品目名, i.品種区分, i.品目コード,
			   a.容量, a.材質,
			   b.内径, b.外径
		FROM 品目基本属性 i
		LEFT JOIN A品種品目属性 a ON i.品目ID = a.品目ID AND i.品種区分 = 'A'
		LEFT JOIN B品種品目属性 b ON i.品目ID = b.品目ID AND i.品種区分 = 'B'
		WHERE 1=1`
	
	countQuery := "SELECT COUNT(*) FROM 品目基本属性 WHERE 1=1"
	
	args := []interface{}{}
	countArgs := []interface{}{}
	argPos := 1
	
	if categoryType != "" {
		query += " AND i.品種区分 = $" + strconv.Itoa(argPos)
		countQuery += " AND 品種区分 = $" + strconv.Itoa(argPos)
		args = append(args, categoryType)
		countArgs = append(countArgs, categoryType)
		argPos++
	}
	
	query += " ORDER BY i.品目ID LIMIT $" + strconv.Itoa(argPos) + " OFFSET $" + strconv.Itoa(argPos+1)
	args = append(args, pageSize, offset)
	
	var total int
	err := DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to count items",
		})
	}
	
	rows, err := DB.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to fetch items",
		})
	}
	defer rows.Close()
	
	items := []models.ItemWithDetails{}
	for rows.Next() {
		var item models.ItemWithDetails
		var typeA models.ItemTypeA
		var typeB models.ItemTypeB
		
		err := rows.Scan(
			&item.ItemID, &item.ItemName, &item.CategoryType, &item.ItemCode,
			&typeA.Capacity, &typeA.Material,
			&typeB.InnerDiameter, &typeB.OuterDiameter,
		)
		if err != nil {
			continue
		}
		
		typeA.ItemID = item.ItemID
		typeB.ItemID = item.ItemID
		
		if item.CategoryType == "A" && (typeA.Capacity.Valid || typeA.Material.Valid) {
			item.TypeA = &typeA
		} else if item.CategoryType == "B" && (typeB.InnerDiameter.Valid || typeB.OuterDiameter.Valid) {
			item.TypeB = &typeB
		}
		
		items = append(items, item)
	}
	
	return c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data: models.ItemListResponse{
			Items:    items,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

func GetItem(c echo.Context) error {
	itemID := c.Param("id")
	
	query := `
		SELECT i.品目ID, i.品目名, i.品種区分, i.品目コード,
			   a.容量, a.材質,
			   b.内径, b.外径
		FROM 品目基本属性 i
		LEFT JOIN A品種品目属性 a ON i.品目ID = a.品目ID
		LEFT JOIN B品種品目属性 b ON i.品目ID = b.品目ID
		WHERE i.品目ID = $1`
	
	var item models.ItemWithDetails
	var typeA models.ItemTypeA
	var typeB models.ItemTypeB
	
	err := DB.QueryRow(query, itemID).Scan(
		&item.ItemID, &item.ItemName, &item.CategoryType, &item.ItemCode,
		&typeA.Capacity, &typeA.Material,
		&typeB.InnerDiameter, &typeB.OuterDiameter,
	)
	
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Item not found",
		})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to fetch item",
		})
	}
	
	typeA.ItemID = item.ItemID
	typeB.ItemID = item.ItemID
	
	if item.CategoryType == "A" && (typeA.Capacity.Valid || typeA.Material.Valid) {
		item.TypeA = &typeA
	} else if item.CategoryType == "B" && (typeB.InnerDiameter.Valid || typeB.OuterDiameter.Valid) {
		item.TypeB = &typeB
	}
	
	return c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    item,
	})
}

func CreateItem(c echo.Context) error {
	var req models.ItemCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body",
		})
	}
	
	if req.CategoryType != "A" && req.CategoryType != "B" {
		return c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Category type must be 'A' or 'B'",
		})
	}
	
	tx, err := DB.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to start transaction",
		})
	}
	defer tx.Rollback()
	
	_, err = tx.Exec(
		"INSERT INTO 品目基本属性 (品目ID, 品目名, 品種区分, 品目コード) VALUES ($1, $2, $3, $4)",
		req.ItemID, req.ItemName, req.CategoryType, req.ItemCode,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Failed to create item: " + err.Error(),
		})
	}
	
	if req.CategoryType == "A" {
		_, err = tx.Exec(
			"INSERT INTO A品種品目属性 (品目ID, 容量, 材質) VALUES ($1, $2, $3)",
			req.ItemID, req.Capacity, req.Material,
		)
	} else if req.CategoryType == "B" {
		_, err = tx.Exec(
			"INSERT INTO B品種品目属性 (品目ID, 内径, 外径) VALUES ($1, $2, $3)",
			req.ItemID, req.InnerDiameter, req.OuterDiameter,
		)
	}
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to create item attributes",
		})
	}
	
	if err = tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to commit transaction",
		})
	}
	
	return GetItem(c)
}

func UpdateItem(c echo.Context) error {
	itemID := c.Param("id")
	
	var req models.ItemUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request body",
		})
	}
	
	var categoryType string
	err := DB.QueryRow("SELECT 品種区分 FROM 品目基本属性 WHERE 品目ID = $1", itemID).Scan(&categoryType)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Item not found",
		})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to fetch item",
		})
	}
	
	tx, err := DB.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to start transaction",
		})
	}
	defer tx.Rollback()
	
	if req.ItemName != nil || req.ItemCode != nil {
		updateQuery := "UPDATE 品目基本属性 SET "
		args := []interface{}{}
		argPos := 1
		
		if req.ItemName != nil {
			updateQuery += "品目名 = $" + strconv.Itoa(argPos) + ", "
			args = append(args, *req.ItemName)
			argPos++
		}
		
		if req.ItemCode != nil {
			updateQuery += "品目コード = $" + strconv.Itoa(argPos) + ", "
			args = append(args, *req.ItemCode)
			argPos++
		}
		
		updateQuery = updateQuery[:len(updateQuery)-2]
		updateQuery += " WHERE 品目ID = $" + strconv.Itoa(argPos)
		args = append(args, itemID)
		
		_, err = tx.Exec(updateQuery, args...)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Error:   "Failed to update item",
			})
		}
	}
	
	if categoryType == "A" && (req.Capacity != nil || req.Material != nil) {
		updateQuery := "UPDATE A品種品目属性 SET "
		args := []interface{}{}
		argPos := 1
		
		if req.Capacity != nil {
			updateQuery += "容量 = $" + strconv.Itoa(argPos) + ", "
			args = append(args, *req.Capacity)
			argPos++
		}
		
		if req.Material != nil {
			updateQuery += "材質 = $" + strconv.Itoa(argPos) + ", "
			args = append(args, *req.Material)
			argPos++
		}
		
		updateQuery = updateQuery[:len(updateQuery)-2]
		updateQuery += " WHERE 品目ID = $" + strconv.Itoa(argPos)
		args = append(args, itemID)
		
		_, err = tx.Exec(updateQuery, args...)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Error:   "Failed to update item attributes",
			})
		}
	} else if categoryType == "B" && (req.InnerDiameter != nil || req.OuterDiameter != nil) {
		updateQuery := "UPDATE B品種品目属性 SET "
		args := []interface{}{}
		argPos := 1
		
		if req.InnerDiameter != nil {
			updateQuery += "内径 = $" + strconv.Itoa(argPos) + ", "
			args = append(args, *req.InnerDiameter)
			argPos++
		}
		
		if req.OuterDiameter != nil {
			updateQuery += "外径 = $" + strconv.Itoa(argPos) + ", "
			args = append(args, *req.OuterDiameter)
			argPos++
		}
		
		updateQuery = updateQuery[:len(updateQuery)-2]
		updateQuery += " WHERE 品目ID = $" + strconv.Itoa(argPos)
		args = append(args, itemID)
		
		_, err = tx.Exec(updateQuery, args...)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Error:   "Failed to update item attributes",
			})
		}
	}
	
	if err = tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to commit transaction",
		})
	}
	
	return GetItem(c)
}

func DeleteItem(c echo.Context) error {
	itemID := c.Param("id")
	
	result, err := DB.Exec("DELETE FROM 品目基本属性 WHERE 品目ID = $1", itemID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to delete item",
		})
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Error:   "Item not found",
		})
	}
	
	return c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Item deleted successfully",
	})
}