package configs

import (
	"fmt"
	"log"
	"os"

	"ecommerce/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host string
	Port string
	Name string
	User string
	Pass string
}

var DB *gorm.DB

func InitDatabase(env *DBConfig) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", env.Host, env.User, env.Pass, env.Name, env.Port)

	dbConnection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database! \n", err.Error())
		os.Exit(1)
	}

	dbConnection.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	log.Println("Running Migrations")
	err = dbConnection.AutoMigrate(&models.Account{}, &models.UserLogin{}, &models.UserOtp{}, &models.Address{})
	if err != nil {
		log.Fatal("Migration Failed:  \n", err.Error())
		os.Exit(1)
	}

	log.Println("Database Connected!")

	DB = dbConnection
}
