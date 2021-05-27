package user

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()

type RepositoryMongo struct {
	db *entity.Service
}

func NewRepositoryMongo(db *entity.Service) *RepositoryMongo {
	r := &RepositoryMongo{
		db: db,
	}

	r.createIndex()

	return r
}

func (r *RepositoryMongo) createIndex() {
	collection := r.getCollection()
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "steam.id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("steam_id"),
		},
		{
			Keys:    bson.D{{Key: "faceit.id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("faceit_id"),
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), models, opts)
	if err != nil {
		log.Error(err)
	}
}

func (r *RepositoryMongo) Create(u *User) error {
	collection := r.getCollection()
	_, err := collection.InsertOne(ctx, u)
	return err
}

func (r *RepositoryMongo) Find(id entity.ID) (*User, error) {
	filterConfig := bson.M{"_id": id}
	return r.filterOne(filterConfig)
}

func (r *RepositoryMongo) FindBySteamId(id uint64) (*User, error) {
	filterConfig := bson.M{"steam.id": id}
	return r.filterOne(filterConfig)
}

func (r *RepositoryMongo) FindByFaceitId(id entity.ID) (*User, error) {
	filterConfig := bson.M{"faceit.id": id}
	return r.filterOne(filterConfig)
}

func (r *RepositoryMongo) FindUsersContainingAuthenticationCode() ([]*User, error) {
	filterConfig := bson.M{"steam.apiEnabled": true}
	return r.filter(filterConfig)
}

func (r *RepositoryMongo) List() ([]*User, error) {
	filterConfig := bson.M{}
	return r.filter(filterConfig)
}

func (r *RepositoryMongo) Delete(id entity.ID) error {
	filter := bson.M{"_id": id}

	res, err := r.getCollection().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no user was deleted")
	}

	return nil
}

func (r *RepositoryMongo) UpdateMatchAuthCode(u *User) error {
	filter := bson.M{"_id": u.ID}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "steam.authCode", Value: u.Steam.AuthCode},
	}}}

	t := &User{}
	return r.getCollection().FindOneAndUpdate(ctx, filter, update).Decode(t)
}

func (r *RepositoryMongo) UpdateLatestShareCode(u *User) error {
	filter := bson.M{"_id": u.ID}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "steam.lastShareCode", Value: u.Steam.LastShareCode},
	}}}

	t := &User{}
	return r.getCollection().FindOneAndUpdate(ctx, filter, update).Decode(t)
}

func (r *RepositoryMongo) UpdateSteamAPIUsage(u *User) error {
	filter := bson.M{"_id": u.ID}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "steam.apiEnabled", Value: u.Steam.APIEnabled},
	}}}

	t := &User{}
	return r.getCollection().FindOneAndUpdate(ctx, filter, update).Decode(t)
}

func (r *RepositoryMongo) getCollection() *mongo.Collection {
	return r.db.GetCollection("users")
}

func (r *RepositoryMongo) filterOne(filter interface{}) (*User, error) {
	var u *User
	res := r.getCollection().FindOne(ctx, filter)
	if err := res.Decode(&u); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, entity.ErrNotFound
		} else {
			log.Debugf("match.infrastructure: %s", err)
			return nil, entity.ErrUnknownInfrastructureError
		}
	}

	return u, nil
}

func (r *RepositoryMongo) filter(filter interface{}) ([]*User, error) {
	var users []*User

	cur, err := r.getCollection().Find(ctx, filter)
	if err != nil {
		return users, err
	}

	for cur.Next(ctx) {
		var m *User
		err := cur.Decode(&m)
		if err != nil {
			return users, err
		}

		users = append(users, m)
	}

	if err := cur.Err(); err != nil {
		return users, err
	}

	cur.Close(ctx)

	if len(users) == 0 {
		return users, nil
	}

	return users, nil
}
