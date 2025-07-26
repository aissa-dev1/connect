package messageservice

import (
	"connect/internal/db"
	messageconstants "connect/internal/message/constants"
	"context"
	"fmt"
	"log"
)

func CreateTableIfNotExists() {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			senderId INT NOT NULL,
			receiverId INT NOT NULL,
			text VARCHAR(1000) NOT NULL,
			readStatus VARCHAR(20) DEFAULT '%s'
		)
	`, messageconstants.ReadStatusSent)
	_, err := db.Pool().Exec(context.Background(), query)

	if err != nil {
		log.Fatalf("Failed to create table messages %v\n", err)
	}
}
