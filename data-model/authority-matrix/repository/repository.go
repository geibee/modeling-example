package repository

import (
	"database/sql"
	"fmt"

	"github.com/modeling-example/data-model/authority-matrix/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}



// UpdateBudget は予算レコードを更新します
func (r *Repository) UpdateBudget(budget models.Budget) error {
	query := `UPDATE budget SET amount = $1, description = $2, updated_at = CURRENT_TIMESTAMP 
	          WHERE id = $3`
	_, err := r.db.Exec(query, budget.Amount, budget.Description, budget.ID)
	if err != nil {
		return fmt.Errorf("予算の更新に失敗: %w", err)
	}
	return nil
}

// CreateBudget は新しい予算レコードを作成します
func (r *Repository) CreateBudget(budget models.Budget) error {
	query := `INSERT INTO budget (id, department_id, org_id, amount, description, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	_, err := r.db.Exec(query, budget.ID, budget.DepartmentID, budget.OrgID, budget.Amount, budget.Description)
	if err != nil {
		return fmt.Errorf("予算の作成に失敗: %w", err)
	}
	return nil
}

// GetUserRoles はユーザーのロール一覧を取得します
func (r *Repository) GetUserRoles(userID string) ([]string, error) {
	query := `SELECT role_id FROM user_role WHERE user_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザーロールの取得に失敗: %w", err)
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var roleID string
		if err := rows.Scan(&roleID); err != nil {
			return nil, fmt.Errorf("ロールのスキャンに失敗: %w", err)
		}
		roles = append(roles, roleID)
	}
	return roles, nil
}

// HasRoleScreenPermission はロールが特定の画面へのアクセス権を持っているかチェックします
func (r *Repository) HasRoleScreenPermission(roleID, screenID string) bool {
	query := `SELECT COUNT(*) FROM role_screen_permission WHERE role_id = $1 AND screen_id = $2`
	var count int
	err := r.db.QueryRow(query, roleID, screenID).Scan(&count)
	return err == nil && count > 0
}

// GetRoleScreenPermissionLevel はロールが特定の画面に対して持つ権限レベルを取得します
func (r *Repository) GetRoleScreenPermissionLevel(roleID, screenID string) (string, error) {
	query := `SELECT permission_level FROM role_screen_permission WHERE role_id = $1 AND screen_id = $2`
	var level string
	err := r.db.QueryRow(query, roleID, screenID).Scan(&level)
	if err != nil {
		return "", err
	}
	return level, nil
}


// GetRoleScreenPermissions はロールが持つすべての画面権限を取得します
func (r *Repository) GetRoleScreenPermissions(roleID string) ([]string, error) {
	query := `SELECT screen_id FROM role_screen_permission WHERE role_id = $1`
	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, fmt.Errorf("ロール画面権限の取得に失敗: %w", err)
	}
	defer rows.Close()

	var screens []string
	for rows.Next() {
		var screenID string
		if err := rows.Scan(&screenID); err != nil {
			return nil, fmt.Errorf("画面IDのスキャンに失敗: %w", err)
		}
		screens = append(screens, screenID)
	}
	return screens, nil
}

// GetScreensForAPI はAPIに紐づく画面一覧を取得します
func (r *Repository) GetScreensForAPI(apiID string) ([]string, error) {
	query := `SELECT screen_id FROM api_screen_mapping WHERE api_id = $1`
	rows, err := r.db.Query(query, apiID)
	if err != nil {
		return nil, fmt.Errorf("API画面マッピングの取得に失敗: %w", err)
	}
	defer rows.Close()

	var screens []string
	for rows.Next() {
		var screenID string
		if err := rows.Scan(&screenID); err != nil {
			return nil, fmt.Errorf("画面IDのスキャンに失敗: %w", err)
		}
		screens = append(screens, screenID)
	}
	return screens, nil
}

// GetUserOrganizations はユーザーが所属する組織情報を取得します
func (r *Repository) GetUserOrganizations(userID string) ([]models.Organization, error) {
	query := `
		SELECT o.org_id, o.org_name 
		FROM user_organization uo
		JOIN organization o ON uo.org_id = o.org_id
		WHERE uo.user_id = $1
		ORDER BY o.ORG_ID
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザー組織の取得に失敗: %w", err)
	}
	defer rows.Close()

	var organizations []models.Organization
	for rows.Next() {
		var org models.Organization
		if err := rows.Scan(&org.OrgID, &org.OrgName); err != nil {
			return nil, fmt.Errorf("組織情報のスキャンに失敗: %w", err)
		}
		organizations = append(organizations, org)
	}
	return organizations, nil
}