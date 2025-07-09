package handlers

import (
	"delivery-app/internal/database"
	"fmt"
)

func ValidateCustomerDeliveryMethod(prefecture, deliveryMethod string) error {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM region_delivery_methods rdm
		JOIN prefectures p ON p.region = rdm.region
		WHERE p.prefecture = $1 AND rdm.delivery_method = $2
	`
	err := database.DB.QueryRow(query, prefecture, deliveryMethod).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to validate delivery method: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("指定された都道府県の地域では、この配達方法は利用できません")
	}
	return nil
}

func CheckCustomerDependencies(region, deliveryMethod string) error {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM customers c
		JOIN prefectures p ON p.prefecture = c.prefecture
		WHERE p.region = $1 AND c.delivery_method = $2
	`
	err := database.DB.QueryRow(query, region, deliveryMethod).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check dependencies: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("この地域と配達方法の組み合わせを使用している顧客が存在するため、削除できません")
	}
	return nil
}

func ValidatePrefectureUpdate(oldPrefecture, newRegion string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.Query(`
		SELECT DISTINCT c.delivery_method 
		FROM customers c 
		WHERE c.prefecture = $1
	`, oldPrefecture)
	if err != nil {
		return fmt.Errorf("failed to query customers: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var deliveryMethod string
		if err := rows.Scan(&deliveryMethod); err != nil {
			return fmt.Errorf("failed to scan delivery method: %w", err)
		}

		var count int
		err := tx.QueryRow(`
			SELECT COUNT(*) 
			FROM region_delivery_methods 
			WHERE region = $1 AND delivery_method = $2
		`, newRegion, deliveryMethod).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check region delivery method: %w", err)
		}
		if count == 0 {
			return fmt.Errorf("地域を変更すると、顧客の配達方法「%s」が無効になります", deliveryMethod)
		}
	}

	return nil
}