package friendshipservice

import (
	"connect/internal/db"
	friendshipconstants "connect/internal/friendship/constants"
	friendshipmodel "connect/internal/friendship/model"
	errormessage "connect/internal/pkg/error_message"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func CreateTableIfNotExists() {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS friendships (
			id SERIAL PRIMARY KEY,
			requesterId INT NOT NULL,
			receiverId INT NOT NULL,
			status VARCHAR(20) DEFAULT '%s'
		);`, friendshipconstants.StatusPending)

	_, err := db.Pool().Exec(context.Background(), query)

	if err != nil {
		log.Fatalf("Failed to create table friendships %v\n", err)
	}
}

func InsertFriendship(friendship friendshipmodel.Friendship) error {
	var insertError error

	if friendship.Status != nil {
		_, err := db.Pool().Exec(context.Background(), `
			INSERT INTO friendships ( requesterId, receiverId, status ) VALUES ( $1, $2, $3 )
		`, friendship.RequesterId, friendship.ReceiverId, friendship.Status)

		insertError = err
	} else {
		_, err := db.Pool().Exec(context.Background(), `
			INSERT INTO friendships ( requesterId, receiverId ) VALUES ( $1, $2 )
		`, friendship.RequesterId, friendship.ReceiverId)

		insertError = err
	}

	if insertError != nil {
		return errors.New(errormessage.InternalServerError)
	}

	return nil
}

func FriendshipExists(requesterId int, receiverId int) (bool, error) {
	var exists bool

	existsRowErr := db.Pool().QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM friendships WHERE (requesterId = $1 AND receiverId = $2) OR (requesterId = $2 AND receiverId = $1));", requesterId, receiverId).Scan(&exists)

	if existsRowErr != nil {
		return exists, errors.New(errormessage.InternalServerError)
	}

	return exists, nil
}

func GetFriendship(requesterId int, receiverId int) (*friendshipmodel.Friendship, error) {
	friendship := new(friendshipmodel.Friendship)

	friendshipErr := db.Pool().QueryRow(context.Background(), "SELECT id, requesterId, receiverId, status FROM friendships WHERE (requesterId = $1 AND receiverId = $2) OR (requesterId = $2 AND receiverId = $1);", requesterId, receiverId).Scan(&friendship.Id, &friendship.RequesterId, &friendship.ReceiverId, &friendship.Status)

	if friendshipErr != nil {
		if errors.Is(friendshipErr, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.New(errormessage.InternalServerError)
	}

	return friendship, nil
}

func MustGetFriendship(requesterId int, receiverId int) (friendshipmodel.Friendship, error) {
	friendshipExists, friendshipExistsErr := FriendshipExists(requesterId, receiverId)

	if friendshipExistsErr != nil {
		return friendshipmodel.Friendship{}, friendshipExistsErr
	}
	if !friendshipExists {
		return friendshipmodel.Friendship{}, friendshipconstants.ErrFriendshipNotFound
	}

	friendship, friendshipErr := GetFriendship(requesterId, receiverId)

	if friendshipErr != nil {
		return friendshipmodel.Friendship{}, friendshipErr
	}
	if friendship == nil {
		return friendshipmodel.Friendship{}, friendshipconstants.ErrFriendshipNotFound
	}

	return *friendship, nil
}

func GetReceiverFriendships(receiverId int) ([]friendshipmodel.Friendship, error) {
	friendships := []friendshipmodel.Friendship{}

	rows, rowsErr := db.Pool().Query(context.Background(), "SELECT id, requesterId, receiverId, status FROM friendships WHERE receiverId = $1;", receiverId)

	if rowsErr != nil {
		return friendships, errors.New(errormessage.InternalServerError)
	}

	defer rows.Close()

	for rows.Next() {
		var friendship friendshipmodel.Friendship

		scanErr := rows.Scan(&friendship.Id, &friendship.RequesterId, &friendship.ReceiverId, &friendship.Status)

		if scanErr != nil {
			return friendships, errors.New(errormessage.InternalServerError)
		}

		friendships = append(friendships, friendship)
	}

	return friendships, nil
}

func GetReceiverFriendshipsByStatus(receiverId int, status friendshipmodel.FriendshipStatus) ([]friendshipmodel.Friendship, error) {
	friendships := []friendshipmodel.Friendship{}

	rows, rowsErr := db.Pool().Query(context.Background(), "SELECT id, requesterId, receiverId, status FROM friendships WHERE receiverId = $1 AND status = $2;", receiverId, status)

	if rowsErr != nil {
		return friendships, errors.New(errormessage.InternalServerError)
	}

	defer rows.Close()

	for rows.Next() {
		var friendship friendshipmodel.Friendship

		scanErr := rows.Scan(&friendship.Id, &friendship.RequesterId, &friendship.ReceiverId, &friendship.Status)

		if scanErr != nil {
			return friendships, errors.New(errormessage.InternalServerError)
		}

		friendships = append(friendships, friendship)
	}

	return friendships, nil
}

func GetRequesterOrReceiverFriendshipsByStatus(userId int, status friendshipmodel.FriendshipStatus) ([]friendshipmodel.Friendship, error) {
	friendships := []friendshipmodel.Friendship{}
	rows, rowsErr := db.Pool().Query(context.Background(), "SELECT id, requesterId, receiverId, status FROM friendships WHERE (requesterId = $1 OR receiverId = $1) AND (status = $2);", userId, status)

	if rowsErr != nil {
		return friendships, errors.New(errormessage.InternalServerError)
	}

	defer rows.Close()

	for rows.Next() {
		var friendship friendshipmodel.Friendship

		scanErr := rows.Scan(&friendship.Id, &friendship.RequesterId, &friendship.ReceiverId, &friendship.Status)

		if scanErr != nil {
			return friendships, errors.New(errormessage.InternalServerError)
		}

		friendships = append(friendships, friendship)
	}

	return friendships, nil
}

func UpdateFriendshipStatus(requesterId int, receiverId int, status friendshipmodel.FriendshipStatus) error {
	_, updateErr := db.Pool().Exec(context.Background(), "UPDATE friendships SET status = $1 WHERE requesterId = $2 AND receiverId = $3;", status, requesterId, receiverId)

	if updateErr != nil {
		return errors.New(errormessage.InternalServerError)
	}

	return nil
}

func MustUpdateFriendshipStatus(requesterId int, receiverId int, status friendshipmodel.FriendshipStatus) error {
	friendshipExists, friendshipExistsErr := FriendshipExists(requesterId, receiverId)

	if friendshipExistsErr != nil {
		return friendshipExistsErr
	}
	if !friendshipExists {
		return friendshipconstants.ErrFriendshipNotFound
	}

	updateErr := UpdateFriendshipStatus(requesterId, receiverId, status)

	if updateErr != nil {
		return updateErr
	}

	return nil
}

func DeleteFriendship(requesterId int, receiverId int) error {
	_, deleteErr := db.Pool().Exec(context.Background(), "DELETE FROM friendships WHERE (requesterId = $1 AND receiverId = $2) OR (requesterId = $2 AND receiverId = $1);", requesterId, receiverId)

	if deleteErr != nil {
		return errors.New(errormessage.InternalServerError)
	}

	return nil
}
