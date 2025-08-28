// Пакет models содержит определения структур данных, используемых в приложении.
package models

// Структура Transaction представляет транзакцию между двумя кошельками.
type Transaction struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}
