package repository

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	proto "test_service/protobuf/generated"
)

type Repository struct {
	// DbConn is the database connection
	DbConn *gorm.DB
}

func NewRepository(dbConfig *proto.DatastoreConfig) (*Repository, error) {
	dbConn, err := initializeDbConn(dbConfig)
	if err != nil {
		return nil, err
	}

	return &Repository{DbConn: dbConn}, nil
}

// initializeDBConn will create a connection to the db based on db config
func initializeDbConn(dbConfig *proto.DatastoreConfig) (*gorm.DB, error) {
	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		dbConfig.FqdnOrIP, dbConfig.Port, dbConfig.Username, dbConfig.DbName, dbConfig.Password)

	dbConn, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// migrate the schema
	//dbConn.AutoMigrate(&models.MyTableDefition{})

	return dbConn, nil
}
