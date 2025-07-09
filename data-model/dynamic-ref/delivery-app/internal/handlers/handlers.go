package handlers

import (
	"database/sql"
	"delivery-app/internal/database"
	"delivery-app/internal/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetCustomers(c echo.Context) error {
	query := `
		SELECT 
			c.customer_id,
			c.customer_name,
			c.prefecture,
			c.delivery_method,
			p.region,
			rdm.standard_delivery_days,
			rdm.standard_delivery_fee,
			c.created_at,
			c.updated_at
		FROM customers c
		JOIN prefectures p ON p.prefecture = c.prefecture
		JOIN region_delivery_methods rdm ON rdm.region = p.region AND rdm.delivery_method = c.delivery_method
		ORDER BY c.customer_id
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	var customers []models.CustomerWithDetails
	for rows.Next() {
		var customer models.CustomerWithDetails
		err := rows.Scan(
			&customer.CustomerID,
			&customer.CustomerName,
			&customer.Prefecture,
			&customer.DeliveryMethod,
			&customer.Region,
			&customer.StandardDeliveryDays,
			&customer.StandardDeliveryFee,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		customers = append(customers, customer)
	}

	return c.JSON(http.StatusOK, customers)
}

func CreateCustomer(c echo.Context) error {
	var customer models.Customer
	if err := c.Bind(&customer); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := ValidateCustomerDeliveryMethod(customer.Prefecture, customer.DeliveryMethod); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	query := `
		INSERT INTO customers (customer_name, prefecture, delivery_method)
		VALUES ($1, $2, $3)
		RETURNING customer_id, created_at, updated_at
	`
	err := database.DB.QueryRow(query, customer.CustomerName, customer.Prefecture, customer.DeliveryMethod).
		Scan(&customer.CustomerID, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, customer)
}

func UpdateCustomer(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid customer ID"})
	}

	var customer models.Customer
	if err := c.Bind(&customer); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := ValidateCustomerDeliveryMethod(customer.Prefecture, customer.DeliveryMethod); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	query := `
		UPDATE customers 
		SET customer_name = $1, prefecture = $2, delivery_method = $3, updated_at = CURRENT_TIMESTAMP
		WHERE customer_id = $4
		RETURNING updated_at
	`
	err = database.DB.QueryRow(query, customer.CustomerName, customer.Prefecture, customer.DeliveryMethod, id).
		Scan(&customer.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Customer not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	customer.CustomerID = id
	return c.JSON(http.StatusOK, customer)
}

func DeleteRegionDeliveryMethod(c echo.Context) error {
	region := c.QueryParam("region")
	deliveryMethod := c.QueryParam("delivery_method")

	if region == "" || deliveryMethod == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Region and delivery_method are required"})
	}

	if err := CheckCustomerDependencies(region, deliveryMethod); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	query := `DELETE FROM region_delivery_methods WHERE region = $1 AND delivery_method = $2`
	result, err := database.DB.Exec(query, region, deliveryMethod)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Region delivery method not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Deleted successfully"})
}

func UpdatePrefecture(c echo.Context) error {
	prefecture := c.Param("prefecture")
	
	var req struct {
		Region string `json:"region"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := ValidatePrefectureUpdate(prefecture, req.Region); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	query := `UPDATE prefectures SET region = $1, updated_at = CURRENT_TIMESTAMP WHERE prefecture = $2`
	result, err := database.DB.Exec(query, req.Region, prefecture)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Prefecture not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Updated successfully"})
}

func GetRegionDeliveryMethods(c echo.Context) error {
	rows, err := database.DB.Query(`
		SELECT region, delivery_method, standard_delivery_days, standard_delivery_fee, created_at, updated_at
		FROM region_delivery_methods
		ORDER BY region, delivery_method
	`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	var methods []models.RegionDeliveryMethod
	for rows.Next() {
		var method models.RegionDeliveryMethod
		err := rows.Scan(
			&method.Region,
			&method.DeliveryMethod,
			&method.StandardDeliveryDays,
			&method.StandardDeliveryFee,
			&method.CreatedAt,
			&method.UpdatedAt,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		methods = append(methods, method)
	}

	return c.JSON(http.StatusOK, methods)
}

func GetPrefectures(c echo.Context) error {
	rows, err := database.DB.Query(`
		SELECT prefecture, region, created_at, updated_at
		FROM prefectures
		ORDER BY prefecture
	`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	var prefectures []models.Prefecture
	for rows.Next() {
		var prefecture models.Prefecture
		err := rows.Scan(
			&prefecture.Prefecture,
			&prefecture.Region,
			&prefecture.CreatedAt,
			&prefecture.UpdatedAt,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		prefectures = append(prefectures, prefecture)
	}

	return c.JSON(http.StatusOK, prefectures)
}