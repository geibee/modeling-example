package models

import (
	"database/sql"
	"time"
)

type Department struct {
	DepartmentID   string    `json:"department_id"`
	DepartmentName string    `json:"department_name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type DepartmentRepository struct {
	db *sql.DB
}

func NewDepartmentRepository(db *sql.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) GetAll() ([]Department, error) {
	query := `SELECT department_id, department_name, created_at, updated_at 
			  FROM departments ORDER BY department_id`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []Department
	for rows.Next() {
		var d Department
		err := rows.Scan(&d.DepartmentID, &d.DepartmentName, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, err
		}
		departments = append(departments, d)
	}

	return departments, nil
}

func (r *DepartmentRepository) GetByID(id string) (*Department, error) {
	query := `SELECT department_id, department_name, created_at, updated_at 
			  FROM departments WHERE department_id = $1`
	
	var d Department
	err := r.db.QueryRow(query, id).Scan(&d.DepartmentID, &d.DepartmentName, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &d, nil
}

func (r *DepartmentRepository) Create(d *Department) error {
	query := `INSERT INTO departments (department_id, department_name) 
			  VALUES ($1, $2) 
			  RETURNING created_at, updated_at`
	
	err := r.db.QueryRow(query, d.DepartmentID, d.DepartmentName).Scan(&d.CreatedAt, &d.UpdatedAt)
	return err
}

func (r *DepartmentRepository) Update(d *Department) error {
	query := `UPDATE departments 
			  SET department_name = $2, updated_at = CURRENT_TIMESTAMP 
			  WHERE department_id = $1 
			  RETURNING updated_at`
	
	err := r.db.QueryRow(query, d.DepartmentID, d.DepartmentName).Scan(&d.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (r *DepartmentRepository) Delete(id string) error {
	query := `DELETE FROM departments WHERE department_id = $1`
	
	result, err := r.db.Exec(query, id)
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