// Пакет api содержит обработчики HTTP-запросов для работы с балансом кошельков и транзакциями.
package api

import (
	"InfotecsTA/internal/models"
	"InfotecsTA/internal/repo"
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	wr *repo.WalletRepository
	tr *repo.TransactionRepository
}

func NewTransactionHandler(walletRepo *repo.WalletRepository, transactionRepo *repo.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{wr: walletRepo, tr: transactionRepo}
}

// GET /api/wallet/:address/balance
// Получение баланса кошелька по адресу
//
//	{
//	    "balance": 1000.0
//	}
func (h *TransactionHandler) GetBalance(c *fiber.Ctx) error {
	address := c.Params("address")
	balance, err := h.wr.GetBalance(address)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).SendString("Wallet not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get balance")
	}
	return c.JSON(fiber.Map{
		"balance": balance,
	})
}

// POST /api/send
// Отправка денег с одного кошелька на другой
//
//	 {
//	 	"from": "address1",
//	 	"to": "address2",
//		"amount": 100.0
//	 }
//
// Возвращает статус операции
func (h *TransactionHandler) SendMoney(c *fiber.Ctx) error {
	var TransactionRequest models.Transaction

	if err := c.BodyParser(&TransactionRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	if err := h.wr.SendMoney(TransactionRequest); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return c.Status(fiber.StatusNotFound).SendString("Wallet not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to send money: " + err.Error())
	}

	return c.Status(fiber.StatusOK).SendString("Transaction completed successfully")
}

// GET /api/transactions?count=N
// Получение последних N транзакций в формате JSON
func (h *TransactionHandler) GetLast(c *fiber.Ctx) error {
	count, err := strconv.Atoi(c.Query("count"))
	if err != nil || count <= 0 {
		count = 10
	}
	transactions, err := h.tr.GetLast(count)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get transactions")
	}
	return c.JSON(transactions)
}
