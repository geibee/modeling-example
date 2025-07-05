package models

import (
	"database/sql"
	"time"
)

type OrganizationAttribute struct {
	DepartmentID       string     `json:"department_id"`
	EffectiveDate      time.Time  `json:"effective_date"`
	ParentDepartmentID *string    `json:"parent_department_id"`
	ExpirationDate     *time.Time `json:"expiration_date"` // 導出属性
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type OrganizationAttributeRepository struct {
	db *sql.DB
}

func NewOrganizationAttributeRepository(db *sql.DB) *OrganizationAttributeRepository {
	return &OrganizationAttributeRepository{db: db}
}

func (r *OrganizationAttributeRepository) GetAll() ([]OrganizationAttribute, error) {
	query := `
		WITH org_attrs AS (
			SELECT 
				oa1.department_id,
				oa1.effective_date,
				oa1.parent_department_id,
				oa1.created_at,
				oa1.updated_at,
				(
					SELECT MIN(oa2.effective_date) - INTERVAL '1 day'
					FROM organization_attributes oa2
					WHERE oa2.department_id = oa1.department_id
					AND oa2.effective_date > oa1.effective_date
				) AS expiration_date
			FROM organization_attributes oa1
		)
		SELECT 
			department_id,
			effective_date,
			parent_department_id,
			expiration_date,
			created_at,
			updated_at
		FROM org_attrs
		ORDER BY department_id, effective_date DESC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attrs []OrganizationAttribute
	for rows.Next() {
		var a OrganizationAttribute
		err := rows.Scan(&a.DepartmentID, &a.EffectiveDate, &a.ParentDepartmentID, 
			&a.ExpirationDate, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		attrs = append(attrs, a)
	}

	return attrs, nil
}

func (r *OrganizationAttributeRepository) GetByID(departmentID string, effectiveDate time.Time) (*OrganizationAttribute, error) {
	query := `
		SELECT 
			oa1.department_id,
			oa1.effective_date,
			oa1.parent_department_id,
			(
				SELECT MIN(oa2.effective_date) - INTERVAL '1 day'
				FROM organization_attributes oa2
				WHERE oa2.department_id = oa1.department_id
				AND oa2.effective_date > oa1.effective_date
			) AS expiration_date,
			oa1.created_at,
			oa1.updated_at
		FROM organization_attributes oa1
		WHERE oa1.department_id = $1 AND oa1.effective_date = $2`
	
	var a OrganizationAttribute
	err := r.db.QueryRow(query, departmentID, effectiveDate).Scan(
		&a.DepartmentID, &a.EffectiveDate, &a.ParentDepartmentID, 
		&a.ExpirationDate, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &a, nil
}

func (r *OrganizationAttributeRepository) GetDepartmentHistory(departmentID string) ([]OrganizationAttribute, error) {
	query := `
		WITH org_attrs AS (
			SELECT 
				oa1.department_id,
				oa1.effective_date,
				oa1.parent_department_id,
				oa1.created_at,
				oa1.updated_at,
				(
					SELECT MIN(oa2.effective_date) - INTERVAL '1 day'
					FROM organization_attributes oa2
					WHERE oa2.department_id = oa1.department_id
					AND oa2.effective_date > oa1.effective_date
				) AS expiration_date
			FROM organization_attributes oa1
			WHERE oa1.department_id = $1
		)
		SELECT 
			department_id,
			effective_date,
			parent_department_id,
			expiration_date,
			created_at,
			updated_at
		FROM org_attrs
		ORDER BY effective_date DESC`
	
	rows, err := r.db.Query(query, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attrs []OrganizationAttribute
	for rows.Next() {
		var a OrganizationAttribute
		err := rows.Scan(&a.DepartmentID, &a.EffectiveDate, &a.ParentDepartmentID, 
			&a.ExpirationDate, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		attrs = append(attrs, a)
	}

	return attrs, nil
}

func (r *OrganizationAttributeRepository) Create(a *OrganizationAttribute) error {
	query := `INSERT INTO organization_attributes (department_id, effective_date, parent_department_id) 
			  VALUES ($1, $2, $3) 
			  RETURNING created_at, updated_at`
	
	err := r.db.QueryRow(query, a.DepartmentID, a.EffectiveDate, a.ParentDepartmentID).
		Scan(&a.CreatedAt, &a.UpdatedAt)
	return err
}

func (r *OrganizationAttributeRepository) Update(a *OrganizationAttribute) error {
	query := `UPDATE organization_attributes 
			  SET parent_department_id = $3, updated_at = CURRENT_TIMESTAMP 
			  WHERE department_id = $1 AND effective_date = $2 
			  RETURNING updated_at`
	
	err := r.db.QueryRow(query, a.DepartmentID, a.EffectiveDate, a.ParentDepartmentID).
		Scan(&a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (r *OrganizationAttributeRepository) Delete(departmentID string, effectiveDate time.Time) error {
	query := `DELETE FROM organization_attributes WHERE department_id = $1 AND effective_date = $2`
	
	result, err := r.db.Exec(query, departmentID, effectiveDate)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}