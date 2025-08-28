package seed

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
)

func RandomAddress() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func SeedWallets(db *sql.DB, count int) error {
	var existing int
	if err := db.QueryRow("SELECT COUNT(*) FROM wallets").Scan(&existing); err != nil {
		return err
	}
	if existing == 10 {
		log.Println("Wallets already seeded")
		return nil
	}

	for i := 0; i < count; i++ {
		addr, err := RandomAddress()
		if err != nil {
			return err
		}
		_, err = db.Exec("INSERT INTO wallets (address, balance) VALUES ($1, $2)", addr, 100.00)
		if err != nil {

			return fmt.Errorf("failed to insert wallet: %w", err)
		}
	}

	return nil
}
