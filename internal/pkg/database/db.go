package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Cludch/csgo-tools/pkg/config"
	log "github.com/sirupsen/logrus"
)

var db *gorm.DB

func init() {
	dbConfig := config.GetConfiguration().Database
	conn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Europe/Berlin",
		dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Database, dbConfig.Port)
	database, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect to database")
	}

	db = database

	// Create / migrate tables
	db.AutoMigrate(&ShareCode{}) //nolint
	db.AutoMigrate(&CSGOUser{})  //nolint
	db.AutoMigrate(&Match{})     //nolint
}

// GetDatabase returns a database connection
func GetDatabase() *gorm.DB {
	return db
}
