package match

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
			Keys:    bson.D{{Key: "shareCode.encoded", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetName("shareCode_encoded"),
		},
		{
			Keys:    bson.D{{Key: "faceitMatchId", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetName("faceitMatchId"),
		},
		{
			Keys:    bson.D{{Key: "filename", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true).SetName("filename"),
		},
	}

	if _, err := collection.Indexes().CreateMany(context.Background(), models, opts); err != nil {
		log.Error(err)
	}
}

func (r *RepositoryMongo) Create(m *Match) error {
	collection := r.getCollection()
	_, err := collection.InsertOne(ctx, m)
	return handleError(err)
}

func (r *RepositoryMongo) Find(id entity.ID) (*Match, error) {
	filterConfig := bson.M{"_id": id}
	m, err := r.filterOne(filterConfig)
	return m, handleError(err)
}

func (r *RepositoryMongo) FindByFilename(filename string) (*Match, error) {
	filterConfig := bson.M{"filename": filename}
	m, err := r.filterOne(filterConfig)
	return m, handleError(err)
}

func (r *RepositoryMongo) FindByFaceitId(id entity.ID) (*Match, error) {
	filterConfig := bson.M{"faceitMatchId": id}
	m, err := r.filterOne(filterConfig)
	return m, handleError(err)
}

func (r *RepositoryMongo) FindByValveId(id uint64) (*Match, error) {
	filterConfig := bson.M{"shareCode.matchId": id}
	m, err := r.filterOne(filterConfig)
	return m, handleError(err)
}

func (r *RepositoryMongo) FindByValveOutcomeId(id uint64) (*Match, error) {
	filterConfig := bson.D{{Key: "shareCode.outcomeId", Value: id}}
	m, err := r.filterOne(filterConfig)
	return m, handleError(err)
}

func (r *RepositoryMongo) List() ([]*Match, error) {
	filterConfig := bson.M{}
	m, err := r.filter(filterConfig)
	return m, handleError(err)
}

func (r *RepositoryMongo) ListDownloadedMatches() ([]*Match, error) {
	filterConfig := bson.M{"status": Downloaded}
	m, err := r.filter(filterConfig)
	return m, handleError(err)
}

func (r *RepositoryMongo) ListDownloadableMatches() ([]*Match, error) {
	filterConfig := bson.M{"status": Downloadable}
	m, err := r.filter(filterConfig)
	return m, handleError(err)
}

func (r *RepositoryMongo) ListParsedMatches() ([]*Match, error) {
	filterConfig := bson.M{"status": Parsed}
	m, err := r.filter(filterConfig)
	return m, handleError(err)
}

func (r *RepositoryMongo) ListValveMatchesMissingDownloadUrl() ([]*Match, error) {
	filterConfig := bson.M{
		"$and": []bson.M{
			{"source": MatchMaking},
			{"status": Created},
		},
	}
	m, err := r.filter(filterConfig)
	return m, handleError(err)
}

func (r *RepositoryMongo) Delete(id entity.ID) error {
	filter := bson.M{"_id": id}

	res, err := r.getCollection().DeleteOne(ctx, filter)
	if err != nil {
		return handleError(err)
	}

	if res.DeletedCount == 0 {
		log.Debug("match: no match was deleted")
	}

	return nil
}

func (r *RepositoryMongo) UpdateResult(m *Match) error {
	filter := bson.M{"_id": m.ID}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "result", Value: m.Result},
	}}}

	t := &Match{}
	return handleError(r.getCollection().FindOneAndUpdate(ctx, filter, update).Decode(t))
}

func (r *RepositoryMongo) UpdateStatus(m *Match) error {
	filter := bson.M{"_id": m.ID}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "status", Value: m.Status},
	}}}

	t := &Match{}
	return handleError(r.getCollection().FindOneAndUpdate(ctx, filter, update).Decode(t))
}

func (r *RepositoryMongo) UpdateStatusAndFilename(m *Match) error {
	filter := bson.M{"_id": m.ID}
	var update primitive.D

	if m.Filename != "" {
		update = bson.D{primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "status", Value: m.Status},
			primitive.E{Key: "filename", Value: m.Filename},
		}}}
	} else {
		update = bson.D{primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "status", Value: m.Status},
		}}}
	}

	t := &Match{}
	return handleError(r.getCollection().FindOneAndUpdate(ctx, filter, update).Decode(t))
}

func (r *RepositoryMongo) UpdateDownloadInformation(m *Match) error {
	filter := bson.M{"_id": m.ID}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "status", Value: m.Status},
		primitive.E{Key: "time", Value: m.Time},
		primitive.E{Key: "url", Value: m.DownloadURL},
	}}}

	t := &Match{}
	return handleError(r.getCollection().FindOneAndUpdate(ctx, filter, update).Decode(t))
}

func (r *RepositoryMongo) getCollection() *mongo.Collection {
	return r.db.GetCollection("matches")
}

func (r *RepositoryMongo) filterOne(filter interface{}) (*Match, error) {
	var m *Match
	res := r.getCollection().FindOne(ctx, filter)
	if err := res.Decode(&m); err != nil {
		return nil, handleError(err)
	}

	return m, nil
}

func (r *RepositoryMongo) filter(filter interface{}) ([]*Match, error) {
	var matches []*Match

	cur, err := r.getCollection().Find(ctx, filter)
	if err != nil {
		return matches, err
	}

	for cur.Next(ctx) {
		var m Match
		if err := handleError(cur.Decode(&m)); err != nil {
			return matches, err
		}

		matches = append(matches, &m)
	}

	if err := handleError(cur.Err()); err != nil {
		return matches, err
	}

	cur.Close(ctx)

	if len(matches) == 0 {
		return matches, nil
	}

	return matches, nil
}

func handleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, entity.ErrNotFound) {
		return entity.ErrNotFound
	} else {
		const msg = "match.infrastructure: %s"
		log.Debugf(msg, err)
		return entity.ErrUnknownInfrastructureError
	}
}
