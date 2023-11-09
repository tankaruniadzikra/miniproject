package main

import (
	"log"
	"miniproject/config"
	"miniproject/docs"
	"miniproject/handler"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()

	// Inisialisasi dokumen Swagger
	docs.SwaggerInfo.Title = "Manufacture"
	docs.SwaggerInfo.Description = "API for rental Manufacture"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Database initialization
	config.InitDB()

	// Routes
	e.POST("/users/register", handler.RegisterUser)
	e.POST("/users/login", handler.LoginUser)

	e.GET("/equipments", handler.GetAllEquipments)

	e.POST("/equipments/rent", handler.RentEquipment, handler.RequireAuth)
	e.DELETE("/equipments/rent/:id", handler.DeleteRentalHistory, handler.RequireAuth)
	e.GET("/equipments/rent", handler.GetAllRentalHistories, handler.RequireAuth)

	e.POST("/users/topup_deposit", handler.TopupDeposit, handler.RequireAuth)
	e.POST("/payments", handler.MakePayment, handler.RequireAuth)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start the server
	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))

}
