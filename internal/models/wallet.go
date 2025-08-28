// Пакет models содержит определения структур данных, используемых в приложении.

package models

// Wallet представляет кошелек с уникальным адресом и балансом.
type Wallet struct {
	ID      int64   `json:"id"`
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}
