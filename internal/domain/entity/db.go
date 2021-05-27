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
	dsn := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v",
		dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	clientOptions := options.Client().ApplyURI(dsn)
	mongoClient, err := mongo.Connect(s.ctx, clientOptions)
	s.client = mongoClient
	if err != nil {
		log.Fatal(err)
	}

	err = s.client.Ping(s.ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) GetCollection(collection string) *mongo.Collection {
	return s.client.Database(s.configurationService.GetConfig().Database.Database).Collection(collection)
}
