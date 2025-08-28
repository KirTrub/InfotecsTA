// Пакет models содержит определения структур данных, используемых в приложении.
package models

// TransactionStamp представляет запись транзакции с деталями.
type TransactionStamp struct {
	ID        int64   `json:"id"`
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
	Timestamp string  `json:"timestamp"`
	Status    string  `json:"status"`
}
