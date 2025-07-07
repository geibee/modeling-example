package models

import (
	"database/sql"
	"time"
)

type ItemBasic struct {
	ItemID       string `json:"item_id" db:"品目id"`
	ItemName     string `json:"item_name" db:"品目名"`
	CategoryType string `json:"category_type" db:"品種区分"`
	ItemCode     string `json:"item_code" db:"品目コード"`
}

type ItemTypeA struct {
	ItemID   string          `json:"item_id" db:"品目id"`
	Capacity sql.NullFloat64 `json:"capacity" db:"容量"`
	Material sql.NullString  `json:"material" db:"材質"`
}

type ItemTypeB struct {
	ItemID        string          `json:"item_id" db:"品目id"`
	InnerDiameter sql.NullFloat64 `json:"inner_diameter" db:"内径"`
	OuterDiameter sql.NullFloat64 `json:"outer_diameter" db:"外径"`
}

type ItemWithDetails struct {
	ItemBasic
	TypeA *ItemTypeA `json:"type_a,omitempty"`
	TypeB *ItemTypeB `json:"type_b,omitempty"`
}

type ItemCreateRequest struct {
	ItemID       string   `json:"item_id"`
	ItemName     string   `json:"item_name"`
	CategoryType string   `json:"category_type"`
	ItemCode     string   `json:"item_code"`
	Capacity     *float64 `json:"capacity,omitempty"`
	Material     *string  `json:"material,omitempty"`
	InnerDiameter *float64 `json:"inner_diameter,omitempty"`
	OuterDiameter *float64 `json:"outer_diameter,omitempty"`
}

type ItemUpdateRequest struct {
	ItemName     *string  `json:"item_name,omitempty"`
	ItemCode     *string  `json:"item_code,omitempty"`
	Capacity     *float64 `json:"capacity,omitempty"`
	Material     *string  `json:"material,omitempty"`
	InnerDiameter *float64 `json:"inner_diameter,omitempty"`
	OuterDiameter *float64 `json:"outer_diameter,omitempty"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ItemListResponse struct {
	Items     []ItemWithDetails `json:"items"`
	Total     int               `json:"total"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
	Timestamp time.Time         `json:"timestamp"`
}