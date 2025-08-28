// Пакет repo содержит реализацию репозиториев для работы с базой данных.

package repo

import (
	"InfotecsTZ/internal/models"
	"database/sql"
)

// TransactionRepository предоставляет методы для работы с таблицей транзакций.
type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Возвращает последние N (count) транзакций
func (r *TransactionRepository) GetLast(count int) ([]models.TransactionStamp, error) {
	rows, err := r.db.Query("SELECT id, from_address, to_address, amount, timestamp, status FROM transactions ORDER BY timestamp DESC LIMIT $1", count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.TransactionStamp
	for rows.Next() {
		var tx models.TransactionStamp
		if err := rows.Scan(&tx.ID, &tx.From, &tx.To, &tx.Amount, &tx.Timestamp, &tx.Status); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// Добавляет запись о транзакции
func (r *TransactionRepository) AddStamp(tx *sql.Tx, from, to string, amount float64, status string) error {
	_, err := tx.Exec("INSERT INTO transactions (from_address, to_address, amount, status) VALUES ($1, $2, $3, $4)", from, to, amount, status)
	return err
}
