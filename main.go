package main

import (
	"fmt"
	"log"
	"os"
	"project/cmd/routes"
	"project/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	logger.Init()
	defer logger.Sync()

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error cargando el archivo .env")
	}

	router := echo.New()

	r := routes.Routes(router)

	r.Use(logger.LoggerMiddleware)

	PORT := os.Getenv("PORT")
	port := fmt.Sprintf(":%s", PORT)

	r.Logger.Fatal(r.Start(port))
}
