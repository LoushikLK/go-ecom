package main

import (
	configs "ecommerce/configs"
	app_middlewares "ecommerce/middlewares"
	"ecommerce/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	var envConfig configs.EnvConfig = configs.AppEnv()
	configs.InitDatabase(&configs.DBConfig{
		Host: envConfig.DB_HOST,
		Port: envConfig.DB_PORT,
		Name: envConfig.DB_NAME,
		User: envConfig.DB_USER,
		Pass: envConfig.DB_PASS,
	}) //initialize database

	configs.InitMemeCache(envConfig.MEMECACHE_SERVER)

	app_middlewares.TopLevelMiddleware(app) //setup middlewares
	routes.InitRoutes(app)                  //setup routes
	app_middlewares.ErrorMiddleware(app)    //parse errors

	log.Fatal(app.Listen(":8000"))
}
