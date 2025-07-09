package models

import (
	"time"
)

type RegionDeliveryMethod struct {
	Region               string    `json:"region"`
	DeliveryMethod       string    `json:"delivery_method"`
	StandardDeliveryDays int       `json:"standard_delivery_days"`
	StandardDeliveryFee  float64   `json:"standard_delivery_fee"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type Prefecture struct {
	Prefecture string    `json:"prefecture"`
	Region     string    `json:"region"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Customer struct {
	CustomerID     int       `json:"customer_id"`
	CustomerName   string    `json:"customer_name"`
	Prefecture     string    `json:"prefecture"`
	DeliveryMethod string    `json:"delivery_method"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CustomerWithDetails struct {
	CustomerID           int       `json:"customer_id"`
	CustomerName         string    `json:"customer_name"`
	Prefecture           string    `json:"prefecture"`
	DeliveryMethod       string    `json:"delivery_method"`
	Region               string    `json:"region"`
	StandardDeliveryDays int       `json:"standard_delivery_days"`
	StandardDeliveryFee  float64   `json:"standard_delivery_fee"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}