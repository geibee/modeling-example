package main

import (
	"delivery-app/internal/database"
	"delivery-app/internal/handlers"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Static("/", "web")
	e.Static("/static", "web/static")

	api := e.Group("/api")
	{
		api.GET("/customers", handlers.GetCustomers)
		api.POST("/customers", handlers.CreateCustomer)
		api.PUT("/customers/:id", handlers.UpdateCustomer)
		
		api.GET("/region-delivery-methods", handlers.GetRegionDeliveryMethods)
		api.DELETE("/region-delivery-methods", handlers.DeleteRegionDeliveryMethod)
		
		api.GET("/prefectures", handlers.GetPrefectures)
		api.PUT("/prefectures/:prefecture", handlers.UpdatePrefecture)
	}

	e.GET("/", func(c echo.Context) error {
		return c.File("web/index.html")
	})

	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}