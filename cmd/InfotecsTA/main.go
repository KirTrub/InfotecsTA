package main

import (
	"InfotecsTA/internal/api"
	"InfotecsTA/internal/db"
	"InfotecsTA/internal/repo"
	"InfotecsTA/internal/seed"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	connStr := os.Getenv("DB_URL")
	database, err := db.Connect(connStr)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer database.Close()

	initSQL, err := os.ReadFile("internal/sql/init.sql")
	if err != nil {
		log.Fatalf("Failed to read init SQL file: %v", err)
	}

	if _, err := database.Exec(string(initSQL)); err != nil {
		log.Fatalf("Failed to execute init SQL: %v", err)
	}

	log.Println("Database initialized successfully")

	if err := seed.SeedWallets(database, 10); err != nil {
		log.Fatalf("Seeding wallets failed: %v", err)
	}

	app := fiber.New()

	tr := repo.NewTransactionRepository(database)
	wr := repo.NewWalletRepository(database, tr)
	th := api.NewTransactionHandler(wr, tr)

	app.Post("/api/send", th.SendMoney)
	app.Get("/api/wallet/:address/balance", th.GetBalance)
	app.Get("/api/transactions", th.GetLast)

	app.Listen(":8081")

}
