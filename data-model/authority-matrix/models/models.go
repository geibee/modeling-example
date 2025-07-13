package models

import (
	"time"
)

// UserScreen はUSER_SCREENテーブルを表現します
type UserScreen struct {
	UserID   string `json:"user_id" db:"user_id"`
	ScreenID string `json:"screen_id" db:"screen_id"`
	Role     string `json:"role" db:"role"`
}

// APIScreen はAPI_SCREENテーブルを表現します
type APIScreen struct {
	APIID    string `json:"api_id" db:"api_id"`
	ScreenID string `json:"screen_id" db:"screen_id"`
}

// Department はユーザーの部署情報を表現します
type Department struct {
	UserID       string `json:"user_id" db:"user_id"`
	DepartmentID string `json:"department_id" db:"department_id"`
}

// Budget は予算データを表現します
type Budget struct {
	ID           string    `json:"id" db:"id"`
	DepartmentID string    `json:"department_id" db:"department_id"`
	OrgID        string    `json:"org_id" db:"org_id"`
	Amount       float64   `json:"amount" db:"amount"`
	Description  string    `json:"description" db:"description"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// BudgetRequest は予算操作用のリクエストペイロードを表現します
type BudgetRequest struct {
	BudgetList []Budget `json:"budget_list"`
}

// Organization は組織情報を表現します
type Organization struct {
	OrgID   string `json:"org_id" db:"ORG_ID"`
	OrgName string `json:"org_name" db:"ORG_NAME"`
}