package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"examples/bloggy/pkg/models"
)

type MongoStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func CreateMongoStore(ctx context.Context, database string, collection string) (Storage, error) {
	// Create mongo client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		return nil, err
	}

	// Create mongo collection
	coll := client.Database(database).Collection(collection)

	// Create index
	mod := mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = coll.Indexes().CreateOne(ctx, mod)

	if err != nil {
		return nil, err
	}

	return &MongoStore{client, coll}, nil
}

func (m *MongoStore) Insert(ctx context.Context, post models.Post) error {
	_, err := m.coll.InsertOne(ctx, post)

	if mongo.IsDuplicateKeyError(err) {
		return ErrAlreadyExists
	}

	if err != nil {
		return err
	}

	return nil
}

func (m *MongoStore) Find(ctx context.Context, title string) (models.Post, error) {
	var foundPost models.Post

	result := m.coll.FindOne(ctx, bson.M{"title": title})
	err := result.Err()

	if err == mongo.ErrNoDocuments {
		return foundPost, ErrDoesNotExist
	}

	if err != nil {
		return foundPost, err
	}

	result.Decode(&foundPost)

	return foundPost, nil
}

func (m *MongoStore) Remove(ctx context.Context, title string) error {
	result, err := m.coll.DeleteOne(ctx, bson.M{"title": title})

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrDoesNotExist
	}

	return nil
}

func (m *MongoStore) Modify(ctx context.Context, title string, post models.Post) error {
	result, err := m.coll.UpdateOne(ctx, bson.M{"title": title}, bson.M{"$set": post})

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrDoesNotExist
	}

	return nil
}

func (m *MongoStore) All(ctx context.Context) ([]models.Post, error) {
	var allPosts []models.Post

	cursor, err := m.coll.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &allPosts)

	if err != nil {
		return nil, err
	}

	return allPosts, nil
}

func (m *MongoStore) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *MongoStore) Clean(ctx context.Context) error {
	_, err := m.coll.DeleteMany(ctx, bson.M{})

	return err
}
