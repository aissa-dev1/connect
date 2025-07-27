package blockservice

import (
	blockconstants "connect/internal/block/constants"
	blockmodel "connect/internal/block/model"
	"connect/internal/db"
	errormessage "connect/internal/pkg/error_message"
	"context"
	"database/sql"
	"errors"
	"log"
)

func CreateTableIfNotExists() {
	_, err := db.Pool().Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS blocks (
			id SERIAL PRIMARY KEY,
			blockerId INT NOT NULL,
			blockedId INT NOT NULL
		);
	`)

	if err != nil {
		log.Fatalf("Failed to create table blocks %v\n", err)
	}
}

func InsertBlock(block blockmodel.Block) error {
	_, err := db.Pool().Exec(context.Background(), "INSERT INTO blocks (blockerId, blockedId) VALUES ($1, $2);", block.BlockerId, block.BlockedId)

	if err != nil {
		return errors.New(errormessage.InternalServerError)
	}

	return nil
}

func BlockExists(blockerId int, blockedId int) (bool, error) {
	var exists bool

	err := db.Pool().QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM blocks WHERE blockerId = $1 AND blockedId = $2);", blockerId, blockedId).Scan(&exists)

	if err != nil {
		return exists, errors.New(errormessage.InternalServerError)
	}

	return exists, nil
}

func BlockExistsMutual(blockerId int, blockedId int) (bool, error) {
	var exists bool

	err := db.Pool().QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM blocks WHERE (blockerId = $1 AND blockedId = $2) OR (blockerId = $2 AND blockedId = $1));", blockerId, blockedId).Scan(&exists)

	if err != nil {
		return exists, errors.New(errormessage.InternalServerError)
	}

	return exists, nil
}

func GetBlock(blockerId int, blockedId int) (*blockmodel.Block, error) {
	block := new(blockmodel.Block)

	blockErr := db.Pool().QueryRow(context.Background(), "SELECT id, blockerId, blockedId FROM blocks WHERE blockerId = $1 AND blockedId = $2;", blockerId, blockedId).Scan(&block.Id, &block.BlockerId, &block.BlockedId)

	if blockErr != nil {
		if errors.Is(blockErr, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.New(errormessage.InternalServerError)
	}

	return block, nil
}

func MustGetBlock(blockerId int, blockedId int) (blockmodel.Block, error) {
	blockExists, blockExistsErr := BlockExists(blockerId, blockedId)
	
	if blockExistsErr != nil {
		return blockmodel.Block{}, blockExistsErr
	}
	if !blockExists {
		return blockmodel.Block{}, blockconstants.ErrBlockNotFound
	}

	block, blockErr := GetBlock(blockerId, blockedId)

	if blockErr != nil {
		return blockmodel.Block{}, blockErr
	}
	if block == nil {
		return blockmodel.Block{}, blockconstants.ErrBlockNotFound
	}

	return *block, nil
}

func GetBlocksBetweenUsers(blockerId int, blockedId int) ([]blockmodel.Block, error) {
	blocks := []blockmodel.Block{}
	
	rows, rowsErr := db.Pool().Query(context.Background(), "SELECT id, blockerId, blockedId FROM blocks WHERE (blockerId = $1 AND blockedId = $2) OR (blockerId = $2 AND blockedId = $1);", blockerId, blockedId)

	if rowsErr != nil {
		return blocks, errors.New(errormessage.InternalServerError)
	}

	defer rows.Close()

	for rows.Next() {
		var block blockmodel.Block

		scanErr := rows.Scan(&block.Id, &block.BlockerId, &block.BlockedId)

		if scanErr != nil {
			return blocks, errors.New(errormessage.InternalServerError)
		}

		blocks = append(blocks, block)
	}

	return blocks, nil
}

func GetBlockerBlocks(blockerId int) ([]blockmodel.Block, error) {
	blocks := []blockmodel.Block{}

	rows, rowsErr := db.Pool().Query(context.Background(), "SELECT id, blockerId, blockedId FROM blocks WHERE blockerId = $1;", blockerId)

	if rowsErr != nil {
		return blocks, errors.New(errormessage.InternalServerError)
	}

	defer rows.Close()

	for rows.Next() {
		var block blockmodel.Block

		scanErr := rows.Scan(&block.Id, &block.BlockerId, &block.BlockedId)

		if scanErr != nil {
			return blocks, errors.New(errormessage.InternalServerError)
		}

		blocks = append(blocks, block)
	}

	return blocks, nil
}
