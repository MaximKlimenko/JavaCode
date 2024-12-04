package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/MaximKlimenko/JavaCode/models"
	"github.com/MaximKlimenko/JavaCode/storage"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetWalletByID(context *fiber.Ctx) error {
	walletModel := models.Wallet{}
	id := context.Params("WALLET_UUID")
	fmt.Println(id)
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id can't be empty"})
		return nil
	}
	err := r.DB.First(&walletModel, "UUID = ?", id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the wallet"})
		fmt.Println(err)
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "balance was received successfully",
			"balance": walletModel.Balance})

	return nil
}

func (r *Repository) ChageBalance(ctx *fiber.Ctx) error {
	var req models.WalletReq

	// Парсим и валидируем входящие данные
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var wallet models.Wallet

	// Используем транзакцию для атомарности
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Получаем или создаём кошелёк с блокировкой строки
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			FirstOrCreate(&wallet, models.Wallet{UUID: req.WalletID}).Error; err != nil {
			return err
		}

		// Обновляем баланс кошелька
		if err := wallet.UpdateBalance(tx, req.OperationType, req.Amount); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if err == gorm.ErrInvalidData {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Insufficient funds",
			})
		}
		fmt.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Возвращаем успешный ответ
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"walletId": wallet.UUID,
		"balance":  wallet.Balance,
	})
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Post("/wallet", r.ChageBalance)
	api.Get("/get/:WALLET_UUID", r.GetWalletByID)
}

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}
	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("could not load the database")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
