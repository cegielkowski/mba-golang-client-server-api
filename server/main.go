package main

import (
	"github.com/cegielkowski/mba-golang-client-server-api/internal/entity"
	"github.com/cegielkowski/mba-golang-client-server-api/internal/infra/database"
	"github.com/cegielkowski/mba-golang-client-server-api/internal/infra/webserver/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("dollar.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(&entity.Dollar{})
	dollarDB := database.NewDollar(db)
	dollarHandler := handlers.NewDollarHandler(dollarDB)

	app := fiber.New()

	app.Get("/cotacao", dollarHandler.GetDollar)

	_ = app.Listen(":8080")
}
