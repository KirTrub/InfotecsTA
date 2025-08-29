// Пакет repo содержит реализацию репозиториев для работы с базой данных.
package repo

import (
	"InfotecsTA/internal/models"
	"database/sql"
	"errors"
)

// WalletRepository предоставляет методы для работы с таблицей кошельков.
type WalletRepository struct {
	db *sql.DB
	tr *TransactionRepository
}

func NewWalletRepository(db *sql.DB, transactionRepo *TransactionRepository) *WalletRepository {
	return &WalletRepository{db: db, tr: transactionRepo}
}

// Возвращает баланс кошелька по его адресу
// Принимает адрес кошелька и возвращает его баланс или ошибку, если кошелек не найден.
func (r *WalletRepository) GetBalance(address string) (float64, error) {
	var balance float64
	err := r.db.QueryRow("SELECT balance FROM wallets WHERE address = $1", address).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

// Переводит деньги с одного кошелька на другой
// Принимает структуру Transaction с информацией о переводе.
// Возвращает ошибку, если перевод не удался (например, недостаточно средств или кошелек не найден).
func (r *WalletRepository) SendMoney(tr models.Transaction) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var fromBalance float64
	err = tx.QueryRow("SELECT balance FROM wallets WHERE address = $1 FOR UPDATE", tr.From).Scan(&fromBalance)
	if err != nil {
		return err
	}

	if fromBalance < tr.Amount {
		err = r.tr.AddStamp(tx, tr.From, tr.To, tr.Amount, "FAILED")
		if err != nil {
			return err
		}
		tx.Commit()
		return errors.New("insufficient funds")
	}

	if _, err := tx.Exec("UPDATE wallets SET balance = balance - $1 WHERE address = $2", tr.Amount, tr.From); err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE wallets SET balance = balance + $1 WHERE address = $2", tr.Amount, tr.To); err != nil {
		return err
	}
	if err := r.tr.AddStamp(tx, tr.From, tr.To, tr.Amount, "OK"); err != nil {
		return err
	}

	return tx.Commit()
}
