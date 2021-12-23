package player

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.TODO()

type RepositoryMongo struct {
	db *entity.Service
}

func NewRepositoryMongo(db *entity.Service) *RepositoryMongo {
	r := &RepositoryMongo{
		db: db,
	}

	return r
}

func (r *RepositoryMongo) Create(m *Player) error {
	collection := r.getCollection()
	_, err := collection.InsertOne(ctx, m)
	return handleError(err)
}

func (r *RepositoryMongo) Find(id uint64) (*Player, error) {
	filterConfig := bson.M{"_id": id}
	p, err := r.filterOne(filterConfig)
	return p, handleError(err)
}

func (r *RepositoryMongo) FindByFaceitId(id entity.ID) (*Player, error) {
	filterConfig := bson.M{"faceitId": id}
	return r.filterOne(filterConfig)
}

func (r *RepositoryMongo) List() ([]*Player, error) {
	filterConfig := bson.M{}
	p, err := r.filter(filterConfig)
	return p, handleError(err)
}

func (r *RepositoryMongo) AddResult(p *Player, result *PlayerResult) error {
	filter := bson.M{"_id": p.ID}

	update := bson.D{primitive.E{Key: "$push", Value: bson.D{
		primitive.E{Key: "results", Value: result},
	}}}

	t := &Player{}
	return handleError(r.getCollection().FindOneAndUpdate(ctx, filter, update).Decode(t))
}

func (r *RepositoryMongo) DeleteResult(p *Player, matchId entity.ID) error {
	filter := bson.M{"_id": p.ID}

	pull := bson.D{primitive.E{Key: "$pull", Value: bson.D{
		primitive.E{Key: "results", Value: bson.D{{Key: "matchId", Value: matchId}}},
	}}}

	t := &Player{}
	return handleError(r.getCollection().FindOneAndUpdate(ctx, filter, pull).Decode(t))
}

func (r *RepositoryMongo) getCollection() *mongo.Collection {
	return r.db.GetCollection("players")
}

func (r *RepositoryMongo) filterOne(filter interface{}) (*Player, error) {
	var p *Player
	res := r.getCollection().FindOne(ctx, filter)
	if err := res.Decode(&p); err != nil {
		return nil, handleError(err)
	}

	return p, nil
}

func (r *RepositoryMongo) filter(filter interface{}) ([]*Player, error) {
	var players []*Player

	cur, err := r.getCollection().Find(ctx, filter)
	if err != nil {
		return players, err
	}

	for cur.Next(ctx) {
		var m Player
		if err := handleError(cur.Decode(&m)); err != nil {
			return players, err
		}

		players = append(players, &m)
	}

	if err := handleError(cur.Err()); err != nil {
		return players, err
	}

	cur.Close(ctx)

	if len(players) == 0 {
		return players, nil
	}

	return players, nil
}

func handleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, entity.ErrNotFound) {
		return entity.ErrNotFound
	} else {
		log.Debugf("player.infrastructure: %s", err)
		return entity.ErrUnknownInfrastructureError
	}
}
