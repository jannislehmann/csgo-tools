package entity

import (
	"context"
	"fmt"
	"log"

	"github.com/Cludch/csgo-tools/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	configurationService config.UseCase
	ctx                  context.Context
	client               *mongo.Client
}

func NewService(c config.UseCase) *Service {
	s := &Service{
		configurationService: c,
	}
	s.connect()
	return s
}

func (s *Service) connect() {
	dbConfig := s.configurationService.GetConfig().Database
	const connString = "mongodb://%v:%v@%v:%v/%v"
	dsn := fmt.Sprintf(connString,
		dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	clientOptions := options.Client().ApplyURI(dsn)
	mongoClient, err := mongo.Connect(s.ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	s.client = mongoClient

	if err = s.client.Ping(s.ctx, nil); err != nil {
		log.Fatal(err)
	}
}

func (s *Service) GetCollection(collection string) *mongo.Collection {
	return s.client.Database(s.configurationService.GetConfig().Database.Database).Collection(collection)
}
