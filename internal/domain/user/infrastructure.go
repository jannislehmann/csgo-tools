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

	if _, err := collection.Indexes().CreateMany(context.Background(), models, opts); err != nil {
		log.Error(err)
	}
}

func (r *RepositoryMongo) Create(u *User) error {
	collection := r.getCollection()
	_, err := collection.InsertOne(ctx, u)
	return handleError(err)
}

func (r *RepositoryMongo) Find(id entity.ID) (*User, error) {
	filterConfig := bson.M{"_id": id}
	u, err := r.filterOne(filterConfig)
	return u, handleError(err)
}

func (r *RepositoryMongo) FindBySteamId(id uint64) (*User, error) {
	filterConfig := bson.M{"steam.id": id}
	u, err := r.filterOne(filterConfig)
	return u, handleError(err)
}

func (r *RepositoryMongo) FindByFaceitId(id entity.ID) (*User, error) {
	filterConfig := bson.M{"faceit.id": id}
	u, err := r.filterOne(filterConfig)
	return u, handleError(err)
}

func (r *RepositoryMongo) FindUsersContainingAuthenticationCode() ([]*User, error) {
	filterConfig := bson.M{"steam.apiEnabled": true}
	u, err := r.filter(filterConfig)
	return u, handleError(err)
}

func (r *RepositoryMongo) List() ([]*User, error) {
	filterConfig := bson.M{}
	u, err := r.filter(filterConfig)
	return u, handleError(err)
}

func (r *RepositoryMongo) Delete(id entity.ID) error {
	filter := bson.M{"_id": id}

	res, err := r.getCollection().DeleteOne(ctx, filter)
	if err != nil {
		return handleError(err)
	}

	if res.DeletedCount == 0 {
		log.Debug("user: no user was deleted")
	}

	return nil
}

func (r *RepositoryMongo) UpdateMatchAuthCode(u *User) error {
	filter := bson.M{"_id": u.ID}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "steam.authCode", Value: u.Steam.AuthCode},
	}}}

	t := &User{}
	return handleError(r.getCollection().FindOneAndUpdate(ctx, filter, update).Decode(t))
}

func (r *RepositoryMongo) UpdateLatestShareCode(u *User) error {
	filter := bson.M{"_id": u.ID}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "steam.lastShareCode", Value: u.Steam.LastShareCode},
	}}}

	t := &User{}
	return handleError(r.getCollection().FindOneAndUpdate(ctx, filter, update).Decode(t))
}

func (r *RepositoryMongo) UpdateSteamAPIUsage(u *User) error {
	filter := bson.M{"_id": u.ID}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "steam.apiEnabled", Value: u.Steam.APIEnabled},
	}}}

	t := &User{}
	return handleError(r.getCollection().FindOneAndUpdate(ctx, filter, update).Decode(t))
}

func (r *RepositoryMongo) getCollection() *mongo.Collection {
	return r.db.GetCollection("users")
}

func (r *RepositoryMongo) filterOne(filter interface{}) (*User, error) {
	var u *User
	res := r.getCollection().FindOne(ctx, filter)
	if err := res.Decode(&u); err != nil {
		return nil, handleError(err)
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
		if err := handleError(cur.Decode(&m)); err != nil {
			return users, err
		}

		users = append(users, m)
	}

	if err := handleError(cur.Err()); err != nil {
		return users, err
	}

	cur.Close(ctx)

	if len(users) == 0 {
		return users, nil
	}

	return users, nil
}

func handleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return entity.ErrNotFound
	} else {
		const msg = "user.infrastructure: %s"
		log.Debugf(msg, err)
		return entity.ErrUnknownInfrastructureError
	}
}
