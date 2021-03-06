package entity

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Cludch/csgo-tools/internal/config"
	log "github.com/sirupsen/logrus"
)

var db *gorm.DB

func init() {
	dbConfig := config.GetConfiguration().Database
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Europe/Berlin",
		dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Database, dbConfig.Port)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect to database")
	}

	db = conn

	// Create / migrate tables.
	db.AutoMigrate(&ShareCode{}) //nolint
	db.AutoMigrate(&CSGOUser{})  //nolint
	db.AutoMigrate(&Match{})     //nolint
}

// GetDatabase returns a database connection.
func GetDatabase() *gorm.DB {
	return db
}
