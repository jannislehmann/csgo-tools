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
	return err
}

func (r *RepositoryMongo) Find(id uint64) (*Player, error) {
	filterConfig := bson.M{"_id": id}
	return r.filterOne(filterConfig)
}

func (r *RepositoryMongo) FindByFaceitId(id entity.ID) (*Player, error) {
	filterConfig := bson.M{"faceitId": id}
	return r.filterOne(filterConfig)
}

func (r *RepositoryMongo) List() ([]*Player, error) {
	filterConfig := bson.M{}
	return r.filter(filterConfig)
}

func (r *RepositoryMongo) AddResult(p *Player, result *PlayerResult) error {
	filter := bson.M{"_id": p.ID}

	update := bson.D{primitive.E{Key: "$push", Value: bson.D{
		primitive.E{Key: "results", Value: result},
	}}}

	t := &Player{}
	return r.getCollection().FindOneAndUpdate(ctx, filter, update).Decode(t)
}

func (r *RepositoryMongo) DeleteResult(p *Player, matchId entity.ID) error {
	filter := bson.M{"_id": p.ID}

	pull := bson.D{primitive.E{Key: "$pull", Value: bson.D{
		primitive.E{Key: "results", Value: bson.D{{Key: "MatchID", Value: matchId}}},
	}}}

	t := &Player{}
	return r.getCollection().FindOneAndUpdate(ctx, filter, pull).Decode(t)
}

func (r *RepositoryMongo) getCollection() *mongo.Collection {
	return r.db.GetCollection("players")
}

func (r *RepositoryMongo) filterOne(filter interface{}) (*Player, error) {
	var p *Player
	res := r.getCollection().FindOne(ctx, filter)
	if err := res.Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return p, entity.ErrNotFound
		} else {
			log.Debugf("player.infrastructure: %s", err)
			return p, entity.ErrUnknownInfrastructureError
		}
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
		err := cur.Decode(&m)
		if err != nil {
			return players, err
		}

		players = append(players, &m)
	}

	if err := cur.Err(); err != nil {
		return players, err
	}

	cur.Close(ctx)

	if len(players) == 0 {
		return players, nil
	}

	return players, nil
}
