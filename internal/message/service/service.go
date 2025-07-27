package messageservice

import (
	"connect/internal/db"
	messageconstants "connect/internal/message/constants"
	messagemodel "connect/internal/message/model"
	errormessage "connect/internal/pkg/error_message"
	"context"
	"errors"
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
		);
	`, messageconstants.ReadStatusSent)
	_, err := db.Pool().Exec(context.Background(), query)

	if err != nil {
		log.Fatalf("Failed to create table messages %v\n", err)
	}
}

func InsertMessage(message messagemodel.Message) error {
	_, err := db.Pool().Exec(context.Background(), "INSERT INTO messages (senderId, receiverId, text, readStatus) VALUES ($1, $2, $3, $4);", message.SenderId, message.ReceiverId, message.Text, message.ReadStatus)

	if err != nil {
		return errors.New(errormessage.InternalServerError) 
	}

	return nil
}

func InsertMessageAndGetId(message messagemodel.Message) (int, error) {
	var id int

	err := db.Pool().QueryRow(context.Background(), "INSERT INTO messages (senderId, receiverId, text, readStatus) VALUES ($1, $2, $3, $4) RETURNING id;", message.SenderId, message.ReceiverId, message.Text, message.ReadStatus).Scan(&id)

	if err != nil {
		return 0, errors.New(errormessage.InternalServerError) 
	}

	return id, nil
}

func GetMessagesBetweenUsers(senderId int, receiverId int) ([]messagemodel.Message, error) {
	messages := []messagemodel.Message{}
	
	rows, rowsErr := db.Pool().Query(context.Background(), "SELECT id, senderId, receiverId, text, readStatus FROM messages WHERE (senderId = $1 AND receiverId = $2) OR (senderId = $2 AND receiverId = $1);", senderId, receiverId)

	if rowsErr != nil {
		return messages, errors.New(errormessage.InternalServerError)
	}

	defer rows.Close()

	for rows.Next() {
		var message messagemodel.Message

		scanErr := rows.Scan(&message.Id, &message.SenderId, &message.ReceiverId, &message.Text, &message.ReadStatus)

		if scanErr != nil {
			return messages, errors.New(errormessage.InternalServerError)
		}

		messages = append(messages, message)
	}

	return messages, nil
}
